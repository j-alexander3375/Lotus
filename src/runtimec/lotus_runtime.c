/* lotus_runtime.c — support runtime linked into compiled Lotus binaries.
 *
 * This file is embedded in the compiler (see llvm_runtime.go), written to a
 * temp file at build time, and passed to clang alongside the generated IR
 * whenever the program uses a runtime-backed stdlib function.
 *
 * ABI: every function takes/returns int64_t ("i64"); pointers travel as i64
 * (the codegen converts with ptrtoint/inttoptr). String-returning functions
 * return char* (i8* on the IR side). This matches the compiler's existing
 * "heap blocks addressed as i64" convention.
 */

#define _GNU_SOURCE
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdarg.h>
#include <unistd.h>
#include <fcntl.h>
#include <time.h>
#include <setjmp.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>

typedef int64_t i64;

/* ========================================================================
 * try/catch/finally/throw support.
 *
 * Implemented with plain libc setjmp/longjmp rather than real LLVM
 * exception handling (landingpad/invoke/a personality routine): Lotus has
 * no stable notion of an "exception object" type to hand a personality
 * function, and setjmp/longjmp gets the same practical behavior (a `throw`
 * anywhere on the call stack - including deep inside unrelated function
 * calls - transfers control straight back to the nearest enclosing `try`)
 * without needing every intervening function to cooperate or be recompiled
 * with unwind tables. The actual `setjmp()` call itself is emitted directly
 * in the generated IR at each `try` (see generateTryStatement in
 * llvm_codegen.go) - only the frame bookkeeping lives here in C.
 *
 * Frames are a plain per-thread stack; `throw` always resumes the MOST
 * RECENTLY entered still-open `try`, popping it first (so if the resulting
 * catch/finally itself throws again - including an explicit re-throw after
 * a finally-only try with no catch - it correctly continues unwinding into
 * the next-outer frame instead of looping back into the one that just
 * fired).
 * ======================================================================== */

#define LOTUS_MAX_TRY_DEPTH 256

static __thread jmp_buf lotus_try_frames[LOTUS_MAX_TRY_DEPTH];
static __thread int lotus_try_top = 0;
static __thread int64_t lotus_current_exception = 0;

/* Returns the buffer for a new try frame and pushes it; the caller passes
 * this pointer directly to the IR-level setjmp() call. */
void *lotus_try_push(void) {
	if (lotus_try_top >= LOTUS_MAX_TRY_DEPTH) {
		fprintf(stderr, "lotus: try nesting exceeds %d levels\n", LOTUS_MAX_TRY_DEPTH);
		exit(1);
	}
	return &lotus_try_frames[lotus_try_top++];
}

/* Pops the current try frame on normal (non-thrown) completion. */
void lotus_try_pop(void) {
	if (lotus_try_top > 0) {
		lotus_try_top--;
	}
}

/* Returns the payload of the exception currently being handled (valid from
 * the moment control resumes after a thrown try's setjmp() call). */
i64 lotus_get_exception(void) {
	return lotus_current_exception;
}

/* throw <expr>: val is the already-evaluated, i64-coerced payload (see
 * generateThrowStatement). */
void lotus_throw(i64 val) {
	lotus_current_exception = val;
	if (lotus_try_top == 0) {
		fprintf(stderr, "lotus: uncaught exception (value=%lld)\n", (long long)val);
		exit(1);
	}
	lotus_try_top--;
	longjmp(lotus_try_frames[lotus_try_top], 1);
}

/* ========================================================================
 * Hash map (int64 keys) — open addressing, linear probing, tombstones.
 * ======================================================================== */

typedef struct {
    i64 cap;    /* slot count, power of two */
    i64 len;    /* live entries */
    i64 tombs;  /* tombstones */
    i64 *keys;
    i64 *vals;
    uint8_t *state; /* 0 empty, 1 used, 2 tombstone */
} lm_imap;

static uint64_t lm_hash_i64(i64 k) {
    uint64_t x = (uint64_t)k;
    x ^= x >> 33; x *= 0xff51afd7ed558ccdULL;
    x ^= x >> 33; x *= 0xc4ceb9fe1a85ec53ULL;
    x ^= x >> 33;
    return x;
}

static lm_imap *lm_imap_alloc(i64 cap_hint) {
    i64 cap = 16;
    while (cap < cap_hint * 2) cap <<= 1;
    lm_imap *m = calloc(1, sizeof(lm_imap));
    m->cap = cap;
    m->keys = calloc(cap, sizeof(i64));
    m->vals = calloc(cap, sizeof(i64));
    m->state = calloc(cap, 1);
    return m;
}

static void lm_imap_grow(lm_imap *m);

/* returns slot of key, or slot to insert at (first tombstone/empty) */
static i64 lm_imap_slot(lm_imap *m, i64 key, int *found) {
    uint64_t mask = (uint64_t)m->cap - 1;
    uint64_t i = lm_hash_i64(key) & mask;
    i64 insert = -1;
    for (;;) {
        switch (m->state[i]) {
        case 0:
            *found = 0;
            return insert >= 0 ? insert : (i64)i;
        case 2:
            if (insert < 0) insert = (i64)i;
            break;
        default:
            if (m->keys[i] == key) { *found = 1; return (i64)i; }
        }
        i = (i + 1) & mask;
    }
}

static void lm_imap_put(lm_imap *m, i64 key, i64 val) {
    if ((m->len + m->tombs) * 10 >= m->cap * 7) lm_imap_grow(m);
    int found;
    i64 s = lm_imap_slot(m, key, &found);
    if (!found) {
        if (m->state[s] == 2) m->tombs--;
        m->state[s] = 1;
        m->keys[s] = key;
        m->len++;
    }
    m->vals[s] = val;
}

static void lm_imap_grow(lm_imap *m) {
    i64 oldcap = m->cap;
    i64 *ok = m->keys, *ov = m->vals;
    uint8_t *os = m->state;
    m->cap = oldcap * 2;
    m->len = 0;
    m->tombs = 0;
    m->keys = calloc(m->cap, sizeof(i64));
    m->vals = calloc(m->cap, sizeof(i64));
    m->state = calloc(m->cap, 1);
    for (i64 i = 0; i < oldcap; i++)
        if (os[i] == 1) lm_imap_put(m, ok[i], ov[i]);
    free(ok); free(ov); free(os);
}

i64 lotus_hashmap_int_new(i64 cap_hint) { return (i64)lm_imap_alloc(cap_hint); }

i64 lotus_hashmap_int_put(i64 mp, i64 key, i64 val) {
    lm_imap_put((lm_imap *)mp, key, val);
    return 0;
}

i64 lotus_hashmap_int_get(i64 mp, i64 key) {
    lm_imap *m = (lm_imap *)mp;
    int found;
    i64 s = lm_imap_slot(m, key, &found);
    return found ? m->vals[s] : 0;
}

i64 lotus_hashmap_int_contains(i64 mp, i64 key) {
    int found;
    lm_imap_slot((lm_imap *)mp, key, &found);
    return found;
}

i64 lotus_hashmap_int_remove(i64 mp, i64 key) {
    lm_imap *m = (lm_imap *)mp;
    int found;
    i64 s = lm_imap_slot(m, key, &found);
    if (!found) return 0;
    m->state[s] = 2;
    m->len--;
    m->tombs++;
    return 1;
}

i64 lotus_hashmap_int_len(i64 mp) { return ((lm_imap *)mp)->len; }

i64 lotus_hashmap_int_clear(i64 mp) {
    lm_imap *m = (lm_imap *)mp;
    memset(m->state, 0, m->cap);
    m->len = 0;
    m->tombs = 0;
    return 0;
}

i64 lotus_hashmap_int_free(i64 mp) {
    lm_imap *m = (lm_imap *)mp;
    free(m->keys); free(m->vals); free(m->state); free(m);
    return 0;
}

/* ========================================================================
 * Hash map (string keys) — same scheme, keys strdup'd.
 * ======================================================================== */

typedef struct {
    i64 cap, len, tombs;
    char **keys;
    i64 *vals;
    uint8_t *state;
} lm_smap;

static uint64_t lm_hash_str(const char *s) {
    uint64_t h = 5381;
    while (*s) h = ((h << 5) + h) ^ (uint8_t)*s++;
    return h;
}

static lm_smap *lm_smap_alloc(i64 cap_hint) {
    i64 cap = 16;
    while (cap < cap_hint * 2) cap <<= 1;
    lm_smap *m = calloc(1, sizeof(lm_smap));
    m->cap = cap;
    m->keys = calloc(cap, sizeof(char *));
    m->vals = calloc(cap, sizeof(i64));
    m->state = calloc(cap, 1);
    return m;
}

static void lm_smap_grow(lm_smap *m);

static i64 lm_smap_slot(lm_smap *m, const char *key, int *found) {
    uint64_t mask = (uint64_t)m->cap - 1;
    uint64_t i = lm_hash_str(key) & mask;
    i64 insert = -1;
    for (;;) {
        switch (m->state[i]) {
        case 0:
            *found = 0;
            return insert >= 0 ? insert : (i64)i;
        case 2:
            if (insert < 0) insert = (i64)i;
            break;
        default:
            if (strcmp(m->keys[i], key) == 0) { *found = 1; return (i64)i; }
        }
        i = (i + 1) & mask;
    }
}

static void lm_smap_put(lm_smap *m, const char *key, i64 val) {
    if ((m->len + m->tombs) * 10 >= m->cap * 7) lm_smap_grow(m);
    int found;
    i64 s = lm_smap_slot(m, key, &found);
    if (!found) {
        if (m->state[s] == 2) { m->tombs--; free(m->keys[s]); }
        m->state[s] = 1;
        m->keys[s] = strdup(key);
        m->len++;
    }
    m->vals[s] = val;
}

static void lm_smap_grow(lm_smap *m) {
    i64 oldcap = m->cap;
    char **ok = m->keys;
    i64 *ov = m->vals;
    uint8_t *os = m->state;
    m->cap = oldcap * 2;
    m->len = 0;
    m->tombs = 0;
    m->keys = calloc(m->cap, sizeof(char *));
    m->vals = calloc(m->cap, sizeof(i64));
    m->state = calloc(m->cap, 1);
    for (i64 i = 0; i < oldcap; i++) {
        if (os[i] == 1) {
            lm_smap_put(m, ok[i], ov[i]);
            free(ok[i]);
        }
    }
    free(ok); free(ov); free(os);
}

i64 lotus_hashmap_str_new(i64 cap_hint) { return (i64)lm_smap_alloc(cap_hint); }

i64 lotus_hashmap_str_put(i64 mp, i64 key, i64 val) {
    lm_smap_put((lm_smap *)mp, (const char *)key, val);
    return 0;
}

i64 lotus_hashmap_str_get(i64 mp, i64 key) {
    lm_smap *m = (lm_smap *)mp;
    int found;
    i64 s = lm_smap_slot(m, (const char *)key, &found);
    return found ? m->vals[s] : 0;
}

i64 lotus_hashmap_str_contains(i64 mp, i64 key) {
    int found;
    lm_smap_slot((lm_smap *)mp, (const char *)key, &found);
    return found;
}

i64 lotus_hashmap_str_remove(i64 mp, i64 key) {
    lm_smap *m = (lm_smap *)mp;
    int found;
    i64 s = lm_smap_slot(m, (const char *)key, &found);
    if (!found) return 0;
    m->state[s] = 2;
    free(m->keys[s]);
    m->keys[s] = NULL;
    m->len--;
    m->tombs++;
    return 1;
}

i64 lotus_hashmap_str_len(i64 mp) { return ((lm_smap *)mp)->len; }

i64 lotus_hashmap_str_clear(i64 mp) {
    lm_smap *m = (lm_smap *)mp;
    for (i64 i = 0; i < m->cap; i++)
        if (m->state[i] == 1) { free(m->keys[i]); m->keys[i] = NULL; }
    memset(m->state, 0, m->cap);
    m->len = 0;
    m->tombs = 0;
    return 0;
}

i64 lotus_hashmap_str_free(i64 mp) {
    lotus_hashmap_str_clear(mp);
    lm_smap *m = (lm_smap *)mp;
    free(m->keys); free(m->vals); free(m->state); free(m);
    return 0;
}

/* ========================================================================
 * Hash sets — thin wrappers over the maps (value ignored).
 * ======================================================================== */

i64 lotus_hashset_int_new(i64 cap_hint) { return lotus_hashmap_int_new(cap_hint); }
i64 lotus_hashset_int_add(i64 s, i64 v) { return lotus_hashmap_int_put(s, v, 1); }
i64 lotus_hashset_int_contains(i64 s, i64 v) { return lotus_hashmap_int_contains(s, v); }
i64 lotus_hashset_int_remove(i64 s, i64 v) { return lotus_hashmap_int_remove(s, v); }
i64 lotus_hashset_int_len(i64 s) { return lotus_hashmap_int_len(s); }
i64 lotus_hashset_int_clear(i64 s) { return lotus_hashmap_int_clear(s); }
i64 lotus_hashset_int_free(i64 s) { return lotus_hashmap_int_free(s); }

i64 lotus_hashset_str_new(i64 cap_hint) { return lotus_hashmap_str_new(cap_hint); }
i64 lotus_hashset_str_add(i64 s, i64 v) { return lotus_hashmap_str_put(s, v, 1); }
i64 lotus_hashset_str_contains(i64 s, i64 v) { return lotus_hashmap_str_contains(s, v); }
i64 lotus_hashset_str_remove(i64 s, i64 v) { return lotus_hashmap_str_remove(s, v); }
i64 lotus_hashset_str_len(i64 s) { return lotus_hashmap_str_len(s); }
i64 lotus_hashset_str_clear(i64 s) { return lotus_hashmap_str_clear(s); }
i64 lotus_hashset_str_free(i64 s) { return lotus_hashmap_str_free(s); }

/* ========================================================================
 * Sorted set / sorted map (int64) — plain BST. Fine for stdlib use;
 * worst-case O(n) on sorted insertion, documented in the registry.
 * ======================================================================== */

typedef struct lm_bst_node {
    i64 key, val;
    struct lm_bst_node *l, *r;
} lm_bst_node;

typedef struct {
    lm_bst_node *root;
    i64 len;
} lm_bst;

/* returns 1 if inserted, 0 if key already present (value updated) */
static int lm_bst_put(lm_bst *t, i64 key, i64 val) {
    lm_bst_node **p = &t->root;
    while (*p) {
        if (key == (*p)->key) { (*p)->val = val; return 0; }
        p = key < (*p)->key ? &(*p)->l : &(*p)->r;
    }
    lm_bst_node *n = calloc(1, sizeof(lm_bst_node));
    n->key = key;
    n->val = val;
    *p = n;
    t->len++;
    return 1;
}

static lm_bst_node *lm_bst_find(lm_bst *t, i64 key) {
    lm_bst_node *n = t->root;
    while (n && n->key != key) n = key < n->key ? n->l : n->r;
    return n;
}

static int lm_bst_remove(lm_bst *t, i64 key) {
    lm_bst_node **p = &t->root;
    while (*p && (*p)->key != key) p = key < (*p)->key ? &(*p)->l : &(*p)->r;
    if (!*p) return 0;
    lm_bst_node *n = *p;
    if (n->l && n->r) {
        /* replace with successor */
        lm_bst_node **s = &n->r;
        while ((*s)->l) s = &(*s)->l;
        n->key = (*s)->key;
        n->val = (*s)->val;
        lm_bst_node *dead = *s;
        *s = dead->r;
        free(dead);
    } else {
        *p = n->l ? n->l : n->r;
        free(n);
    }
    t->len--;
    return 1;
}

static void lm_bst_free_nodes(lm_bst_node *n) {
    if (!n) return;
    lm_bst_free_nodes(n->l);
    lm_bst_free_nodes(n->r);
    free(n);
}

i64 lotus_sortedset_int_new(void) { return (i64)calloc(1, sizeof(lm_bst)); }
i64 lotus_sortedset_int_add(i64 t, i64 v) { return lm_bst_put((lm_bst *)t, v, 0); }
i64 lotus_sortedset_int_contains(i64 t, i64 v) { return lm_bst_find((lm_bst *)t, v) != NULL; }
i64 lotus_sortedset_int_remove(i64 t, i64 v) { return lm_bst_remove((lm_bst *)t, v); }

i64 lotus_sortedset_int_min(i64 tp) {
    lm_bst_node *n = ((lm_bst *)tp)->root;
    if (!n) return 0;
    while (n->l) n = n->l;
    return n->key;
}

i64 lotus_sortedset_int_max(i64 tp) {
    lm_bst_node *n = ((lm_bst *)tp)->root;
    if (!n) return 0;
    while (n->r) n = n->r;
    return n->key;
}

i64 lotus_sortedset_int_len(i64 t) { return ((lm_bst *)t)->len; }

i64 lotus_sortedset_int_free(i64 tp) {
    lm_bst *t = (lm_bst *)tp;
    lm_bst_free_nodes(t->root);
    free(t);
    return 0;
}

i64 lotus_sortedmap_int_new(void) { return lotus_sortedset_int_new(); }
i64 lotus_sortedmap_int_put(i64 t, i64 k, i64 v) { return lm_bst_put((lm_bst *)t, k, v); }

i64 lotus_sortedmap_int_get(i64 t, i64 k) {
    lm_bst_node *n = lm_bst_find((lm_bst *)t, k);
    return n ? n->val : 0;
}

i64 lotus_sortedmap_int_contains(i64 t, i64 k) { return lm_bst_find((lm_bst *)t, k) != NULL; }
i64 lotus_sortedmap_int_remove(i64 t, i64 k) { return lm_bst_remove((lm_bst *)t, k); }
i64 lotus_sortedmap_int_min_key(i64 t) { return lotus_sortedset_int_min(t); }
i64 lotus_sortedmap_int_max_key(i64 t) { return lotus_sortedset_int_max(t); }
i64 lotus_sortedmap_int_len(i64 t) { return ((lm_bst *)t)->len; }
i64 lotus_sortedmap_int_free(i64 t) { return lotus_sortedset_int_free(t); }

/* ========================================================================
 * binary_search_int(base_ptr, len, target) -> index or -1
 * (base_ptr is a raw i64 data pointer, matching the old asm backend)
 * ======================================================================== */

i64 lotus_binary_search_int(i64 base, i64 len, i64 target) {
    const i64 *a = (const i64 *)base;
    i64 lo = 0, hi = len;
    while (lo < hi) {
        i64 mid = lo + (hi - lo) / 2;
        if (a[mid] == target) return mid;
        if (a[mid] < target) lo = mid + 1;
        else hi = mid;
    }
    return -1;
}

/* ========================================================================
 * SHA-256 — writes 32 raw big-endian bytes to out (old backend contract).
 * ======================================================================== */

static const uint32_t sha256_k[64] = {
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1,
    0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
    0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786,
    0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147,
    0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
    0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b,
    0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a,
    0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
    0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2};

#define ROTR32(x, n) (((x) >> (n)) | ((x) << (32 - (n))))

static void sha256_block(uint32_t h[8], const uint8_t *p) {
    uint32_t w[64];
    for (int i = 0; i < 16; i++)
        w[i] = (uint32_t)p[i * 4] << 24 | (uint32_t)p[i * 4 + 1] << 16 |
               (uint32_t)p[i * 4 + 2] << 8 | p[i * 4 + 3];
    for (int i = 16; i < 64; i++) {
        uint32_t s0 = ROTR32(w[i - 15], 7) ^ ROTR32(w[i - 15], 18) ^ (w[i - 15] >> 3);
        uint32_t s1 = ROTR32(w[i - 2], 17) ^ ROTR32(w[i - 2], 19) ^ (w[i - 2] >> 10);
        w[i] = w[i - 16] + s0 + w[i - 7] + s1;
    }
    uint32_t a = h[0], b = h[1], c = h[2], d = h[3];
    uint32_t e = h[4], f = h[5], g = h[6], hh = h[7];
    for (int i = 0; i < 64; i++) {
        uint32_t S1 = ROTR32(e, 6) ^ ROTR32(e, 11) ^ ROTR32(e, 25);
        uint32_t ch = (e & f) ^ (~e & g);
        uint32_t t1 = hh + S1 + ch + sha256_k[i] + w[i];
        uint32_t S0 = ROTR32(a, 2) ^ ROTR32(a, 13) ^ ROTR32(a, 22);
        uint32_t maj = (a & b) ^ (a & c) ^ (b & c);
        uint32_t t2 = S0 + maj;
        hh = g; g = f; f = e; e = d + t1;
        d = c; c = b; b = a; a = t1 + t2;
    }
    h[0] += a; h[1] += b; h[2] += c; h[3] += d;
    h[4] += e; h[5] += f; h[6] += g; h[7] += hh;
}

i64 lotus_sha256(i64 data, i64 len, i64 out) {
    const uint8_t *msg = (const uint8_t *)data;
    uint8_t *dst = (uint8_t *)out;
    uint32_t h[8] = {0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
                     0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19};
    i64 i = 0;
    for (; i + 64 <= len; i += 64) sha256_block(h, msg + i);
    uint8_t tail[128];
    i64 rem = len - i;
    memcpy(tail, msg + i, rem);
    tail[rem] = 0x80;
    i64 padlen = (rem < 56) ? 64 : 128;
    memset(tail + rem + 1, 0, padlen - rem - 1 - 8);
    uint64_t bits = (uint64_t)len * 8;
    for (int j = 0; j < 8; j++) tail[padlen - 1 - j] = (uint8_t)(bits >> (8 * j));
    sha256_block(h, tail);
    if (padlen == 128) sha256_block(h, tail + 64);
    for (int j = 0; j < 8; j++) {
        dst[j * 4] = (uint8_t)(h[j] >> 24);
        dst[j * 4 + 1] = (uint8_t)(h[j] >> 16);
        dst[j * 4 + 2] = (uint8_t)(h[j] >> 8);
        dst[j * 4 + 3] = (uint8_t)h[j];
    }
    return 0;
}

/* ========================================================================
 * MD5 — writes 16 raw little-endian bytes to out.
 * ======================================================================== */

static const uint32_t md5_k[64] = {
    0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee, 0xf57c0faf, 0x4787c62a,
    0xa8304613, 0xfd469501, 0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
    0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821, 0xf61e2562, 0xc040b340,
    0x265e5a51, 0xe9b6c7aa, 0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
    0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed, 0xa9e3e905, 0xfcefa3f8,
    0x676f02d9, 0x8d2a4c8a, 0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
    0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70, 0x289b7ec6, 0xeaa127fa,
    0xd4ef3085, 0x04881d05, 0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
    0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039, 0x655b59c3, 0x8f0ccc92,
    0xffeff47d, 0x85845dd1, 0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
    0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391};

static const int md5_r[64] = {7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7,
                              12, 17, 22, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14,
                              20, 5, 9, 14, 20, 4, 11, 16, 23, 4, 11, 16, 23, 4,
                              11, 16, 23, 4, 11, 16, 23, 6, 10, 15, 21, 6, 10,
                              15, 21, 6, 10, 15, 21, 6, 10, 15, 21};

#define ROTL32(x, n) (((x) << (n)) | ((x) >> (32 - (n))))

static void md5_block(uint32_t h[4], const uint8_t *p) {
    uint32_t m[16];
    for (int i = 0; i < 16; i++)
        m[i] = (uint32_t)p[i * 4] | (uint32_t)p[i * 4 + 1] << 8 |
               (uint32_t)p[i * 4 + 2] << 16 | (uint32_t)p[i * 4 + 3] << 24;
    uint32_t a = h[0], b = h[1], c = h[2], d = h[3];
    for (int i = 0; i < 64; i++) {
        uint32_t f;
        int g;
        if (i < 16) { f = (b & c) | (~b & d); g = i; }
        else if (i < 32) { f = (d & b) | (~d & c); g = (5 * i + 1) % 16; }
        else if (i < 48) { f = b ^ c ^ d; g = (3 * i + 5) % 16; }
        else { f = c ^ (b | ~d); g = (7 * i) % 16; }
        uint32_t tmp = d;
        d = c;
        c = b;
        b = b + ROTL32(a + f + md5_k[i] + m[g], md5_r[i]);
        a = tmp;
    }
    h[0] += a; h[1] += b; h[2] += c; h[3] += d;
}

i64 lotus_md5(i64 data, i64 len, i64 out) {
    const uint8_t *msg = (const uint8_t *)data;
    uint8_t *dst = (uint8_t *)out;
    uint32_t h[4] = {0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476};
    i64 i = 0;
    for (; i + 64 <= len; i += 64) md5_block(h, msg + i);
    uint8_t tail[128];
    i64 rem = len - i;
    memcpy(tail, msg + i, rem);
    tail[rem] = 0x80;
    i64 padlen = (rem < 56) ? 64 : 128;
    memset(tail + rem + 1, 0, padlen - rem - 1 - 8);
    uint64_t bits = (uint64_t)len * 8;
    for (int j = 0; j < 8; j++) tail[padlen - 8 + j] = (uint8_t)(bits >> (8 * j));
    md5_block(h, tail);
    if (padlen == 128) md5_block(h, tail + 64);
    for (int j = 0; j < 4; j++) {
        dst[j * 4] = (uint8_t)h[j];
        dst[j * 4 + 1] = (uint8_t)(h[j] >> 8);
        dst[j * 4 + 2] = (uint8_t)(h[j] >> 16);
        dst[j * 4 + 3] = (uint8_t)(h[j] >> 24);
    }
    return 0;
}

/* ========================================================================
 * Networking (IPv6 + DNS via getaddrinfo)
 * ======================================================================== */

/* connect_ipv6(ipv6_addr_ptr(16 bytes), port, socket_type 1=TCP 2=UDP)
 * -> fd on success; closes the socket and returns -1 on connect failure */
i64 lotus_connect_ipv6(i64 addr, i64 port, i64 socktype) {
    int fd = socket(AF_INET6, socktype == 2 ? SOCK_DGRAM : SOCK_STREAM, 0);
    if (fd < 0) return -1;
    struct sockaddr_in6 sa;
    memset(&sa, 0, sizeof(sa));
    sa.sin6_family = AF_INET6;
    sa.sin6_port = htons((uint16_t)port);
    memcpy(&sa.sin6_addr, (const void *)addr, 16);
    if (connect(fd, (struct sockaddr *)&sa, sizeof(sa)) < 0) {
        close(fd);
        return -1;
    }
    return fd;
}

/* bind_ipv6(fd, addr_ptr_or_0, port) -> bind() result */
i64 lotus_bind_ipv6(i64 fd, i64 addr, i64 port) {
    struct sockaddr_in6 sa;
    memset(&sa, 0, sizeof(sa)); /* zeroed addr == in6addr_any */
    sa.sin6_family = AF_INET6;
    sa.sin6_port = htons((uint16_t)port);
    if (addr) memcpy(&sa.sin6_addr, (const void *)addr, 16);
    return bind((int)fd, (struct sockaddr *)&sa, sizeof(sa));
}

/* sendto_ipv6(fd, buf, len, addr_ptr, port) -> bytes sent */
i64 lotus_sendto_ipv6(i64 fd, i64 buf, i64 len, i64 addr, i64 port) {
    struct sockaddr_in6 sa;
    memset(&sa, 0, sizeof(sa));
    sa.sin6_family = AF_INET6;
    sa.sin6_port = htons((uint16_t)port);
    memcpy(&sa.sin6_addr, (const void *)addr, 16);
    return sendto((int)fd, (const void *)buf, (size_t)len, 0,
                  (struct sockaddr *)&sa, sizeof(sa));
}

/* resolve(hostname, out_ipv4_ptr(4 bytes)) -> 1 on success, 0 on failure */
i64 lotus_resolve(i64 host, i64 out) {
    struct addrinfo hints, *res = NULL;
    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_INET;
    if (getaddrinfo((const char *)host, NULL, &hints, &res) != 0 || !res)
        return 0;
    struct sockaddr_in *sa = (struct sockaddr_in *)res->ai_addr;
    memcpy((void *)out, &sa->sin_addr, 4);
    freeaddrinfo(res);
    return 1;
}

/* resolve_ipv6(hostname, out_ipv6_ptr(16 bytes)) -> 1 on success, 0 on failure */
i64 lotus_resolve_ipv6(i64 host, i64 out) {
    struct addrinfo hints, *res = NULL;
    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_INET6;
    if (getaddrinfo((const char *)host, NULL, &hints, &res) != 0 || !res)
        return 0;
    struct sockaddr_in6 *sa = (struct sockaddr_in6 *)res->ai_addr;
    memcpy((void *)out, &sa->sin6_addr, 16);
    freeaddrinfo(res);
    return 1;
}

/* ========================================================================
 * File helpers
 * ======================================================================== */

i64 lotus_file_seek(i64 fd, i64 offset, i64 whence) {
    return lseek((int)fd, (off_t)offset, (int)whence);
}

/* stat(path, buf): buf receives the platform's raw struct stat (as the old
 * syscall-based backend did). Returns 0 on success, -1 on error. */
i64 lotus_file_stat(i64 path, i64 buf) {
    return stat((const char *)path, (struct stat *)buf);
}

i64 lotus_file_exists(i64 path) {
    return access((const char *)path, F_OK) == 0 ? 1 : 0;
}

/* ========================================================================
 * Time helpers — write the platform's struct tm into the caller's buffer
 * (allocate at least 64 bytes). Readable fields follow C's struct tm:
 * tm_sec, tm_min, tm_hour, tm_mday, tm_mon, tm_year, tm_wday, tm_yday
 * (each a 32-bit int, in that order).
 * ======================================================================== */

i64 lotus_gmtime(i64 ts, i64 buf) {
    time_t t = (time_t)ts;
    gmtime_r(&t, (struct tm *)buf);
    return 0;
}

i64 lotus_localtime(i64 ts, i64 buf) {
    time_t t = (time_t)ts;
    localtime_r(&t, (struct tm *)buf);
    return 0;
}

/* ========================================================================
 * Memory helpers
 * ======================================================================== */

i64 lotus_mmap(i64 size) {
    void *p = mmap(NULL, (size_t)size, PROT_READ | PROT_WRITE,
                   MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    return p == MAP_FAILED ? 0 : (i64)p;
}

i64 lotus_munmap(i64 ptr, i64 size) {
    return munmap((void *)ptr, (size_t)size);
}

/* ========================================================================
 * String helpers
 * ======================================================================== */

/* join(array_ptr, sep): array uses the collections array layout
 * {cap@0, len@8, elements@16}, elements are char* stored as i64.
 * Returns a newly malloc'd string. */
char *lotus_str_join(i64 arr, i64 sep) {
    const char *s = (const char *)sep;
    if (!arr) return strdup("");
    i64 len = *((i64 *)arr + 1);
    const i64 *elems = (const i64 *)arr + 2;
    size_t seplen = strlen(s), total = 1;
    for (i64 i = 0; i < len; i++) {
        if (elems[i]) total += strlen((const char *)elems[i]);
        if (i + 1 < len) total += seplen;
    }
    char *out = malloc(total);
    char *p = out;
    for (i64 i = 0; i < len; i++) {
        if (elems[i]) {
            size_t l = strlen((const char *)elems[i]);
            memcpy(p, (const char *)elems[i], l);
            p += l;
        }
        if (i + 1 < len) {
            memcpy(p, s, seplen);
            p += seplen;
        }
    }
    *p = '\0';
    return out;
}

/* printf-style formatting into a fresh malloc'd string (io::sprint family) */
char *lotus_asprintf(const char *fmt, ...) {
    char *s = NULL;
    va_list ap;
    va_start(ap, fmt);
    if (vasprintf(&s, fmt, ap) < 0) s = strdup("");
    va_end(ap);
    return s;
}

/* ========================================================================
 * JSON — real recursive-descent parser and tree.
 * Node types: 0 null, 1 bool, 2 number, 3 string, 4 array, 5 object.
 * ======================================================================== */

enum { LJ_NULL, LJ_BOOL, LJ_NUM, LJ_STR, LJ_ARR, LJ_OBJ };

typedef struct lj_val {
    int type;
    i64 num;        /* bool 0/1, or integer value of number */
    double dbl;     /* full-precision number */
    char *str;      /* LJ_STR */
    struct lj_val **items; /* LJ_ARR elements / LJ_OBJ values */
    char **keys;    /* LJ_OBJ keys */
    i64 count;
} lj_val;

typedef struct {
    const char *p;
    int err;
} lj_parser;

static void lj_skip_ws(lj_parser *ps) {
    while (*ps->p == ' ' || *ps->p == '\t' || *ps->p == '\n' || *ps->p == '\r')
        ps->p++;
}

static lj_val *lj_new(int type) {
    lj_val *v = calloc(1, sizeof(lj_val));
    v->type = type;
    return v;
}

static void lj_free(lj_val *v) {
    if (!v) return;
    free(v->str);
    for (i64 i = 0; i < v->count; i++) {
        lj_free(v->items[i]);
        if (v->keys) free(v->keys[i]);
    }
    free(v->items);
    free(v->keys);
    free(v);
}

static lj_val *lj_parse_value(lj_parser *ps);

static char *lj_parse_string_raw(lj_parser *ps) {
    if (*ps->p != '"') { ps->err = 1; return NULL; }
    ps->p++;
    size_t cap = 16, n = 0;
    char *out = malloc(cap);
    while (*ps->p && *ps->p != '"') {
        char c = *ps->p++;
        if (c == '\\') {
            char e = *ps->p++;
            switch (e) {
            case 'n': c = '\n'; break;
            case 't': c = '\t'; break;
            case 'r': c = '\r'; break;
            case 'b': c = '\b'; break;
            case 'f': c = '\f'; break;
            case '"': c = '"'; break;
            case '\\': c = '\\'; break;
            case '/': c = '/'; break;
            case 'u': {
                unsigned code = 0;
                for (int i = 0; i < 4; i++) {
                    char h = *ps->p;
                    if (h >= '0' && h <= '9') code = code * 16 + (h - '0');
                    else if (h >= 'a' && h <= 'f') code = code * 16 + (h - 'a' + 10);
                    else if (h >= 'A' && h <= 'F') code = code * 16 + (h - 'A' + 10);
                    else { ps->err = 1; free(out); return NULL; }
                    ps->p++;
                }
                /* UTF-8 encode */
                if (n + 4 >= cap) { cap *= 2; out = realloc(out, cap); }
                if (code < 0x80) out[n++] = (char)code;
                else if (code < 0x800) {
                    out[n++] = (char)(0xC0 | (code >> 6));
                    out[n++] = (char)(0x80 | (code & 0x3F));
                } else {
                    out[n++] = (char)(0xE0 | (code >> 12));
                    out[n++] = (char)(0x80 | ((code >> 6) & 0x3F));
                    out[n++] = (char)(0x80 | (code & 0x3F));
                }
                continue;
            }
            default: ps->err = 1; free(out); return NULL;
            }
        }
        if (n + 2 >= cap) { cap *= 2; out = realloc(out, cap); }
        out[n++] = c;
    }
    if (*ps->p != '"') { ps->err = 1; free(out); return NULL; }
    ps->p++;
    out[n] = '\0';
    return out;
}

static lj_val *lj_parse_object(lj_parser *ps) {
    lj_val *v = lj_new(LJ_OBJ);
    ps->p++; /* '{' */
    lj_skip_ws(ps);
    if (*ps->p == '}') { ps->p++; return v; }
    size_t cap = 4;
    v->keys = malloc(cap * sizeof(char *));
    v->items = malloc(cap * sizeof(lj_val *));
    for (;;) {
        lj_skip_ws(ps);
        char *key = lj_parse_string_raw(ps);
        if (ps->err) { lj_free(v); return NULL; }
        lj_skip_ws(ps);
        if (*ps->p != ':') { free(key); ps->err = 1; lj_free(v); return NULL; }
        ps->p++;
        lj_val *item = lj_parse_value(ps);
        if (ps->err) { free(key); lj_free(v); return NULL; }
        if ((size_t)v->count == cap) {
            cap *= 2;
            v->keys = realloc(v->keys, cap * sizeof(char *));
            v->items = realloc(v->items, cap * sizeof(lj_val *));
        }
        v->keys[v->count] = key;
        v->items[v->count] = item;
        v->count++;
        lj_skip_ws(ps);
        if (*ps->p == ',') { ps->p++; continue; }
        if (*ps->p == '}') { ps->p++; return v; }
        ps->err = 1;
        lj_free(v);
        return NULL;
    }
}

static lj_val *lj_parse_array(lj_parser *ps) {
    lj_val *v = lj_new(LJ_ARR);
    ps->p++; /* '[' */
    lj_skip_ws(ps);
    if (*ps->p == ']') { ps->p++; return v; }
    size_t cap = 4;
    v->items = malloc(cap * sizeof(lj_val *));
    for (;;) {
        lj_val *item = lj_parse_value(ps);
        if (ps->err) { lj_free(v); return NULL; }
        if ((size_t)v->count == cap) {
            cap *= 2;
            v->items = realloc(v->items, cap * sizeof(lj_val *));
        }
        v->items[v->count++] = item;
        lj_skip_ws(ps);
        if (*ps->p == ',') { ps->p++; continue; }
        if (*ps->p == ']') { ps->p++; return v; }
        ps->err = 1;
        lj_free(v);
        return NULL;
    }
}

static lj_val *lj_parse_value(lj_parser *ps) {
    lj_skip_ws(ps);
    switch (*ps->p) {
    case '{':
        return lj_parse_object(ps);
    case '[':
        return lj_parse_array(ps);
    case '"': {
        char *s = lj_parse_string_raw(ps);
        if (ps->err) return NULL;
        lj_val *v = lj_new(LJ_STR);
        v->str = s;
        return v;
    }
    case 't':
        if (strncmp(ps->p, "true", 4) == 0) {
            ps->p += 4;
            lj_val *v = lj_new(LJ_BOOL);
            v->num = 1;
            v->dbl = 1;
            return v;
        }
        ps->err = 1;
        return NULL;
    case 'f':
        if (strncmp(ps->p, "false", 5) == 0) {
            ps->p += 5;
            return lj_new(LJ_BOOL);
        }
        ps->err = 1;
        return NULL;
    case 'n':
        if (strncmp(ps->p, "null", 4) == 0) {
            ps->p += 4;
            return lj_new(LJ_NULL);
        }
        ps->err = 1;
        return NULL;
    default: {
        char *end = NULL;
        double d = strtod(ps->p, &end);
        if (end == ps->p) { ps->err = 1; return NULL; }
        lj_val *v = lj_new(LJ_NUM);
        v->dbl = d;
        v->num = (i64)d;
        ps->p = end;
        return v;
    }
    }
}

i64 lotus_json_parse(i64 str) {
    if (!str) return 0;
    lj_parser ps = {(const char *)str, 0};
    lj_val *v = lj_parse_value(&ps);
    if (ps.err) return 0;
    lj_skip_ws(&ps);
    if (*ps.p != '\0') { lj_free(v); return 0; }
    return (i64)v;
}

i64 lotus_json_is_valid(i64 str) {
    i64 v = lotus_json_parse(str);
    if (!v) return 0;
    lj_free((lj_val *)v);
    return 1;
}

i64 lotus_json_get(i64 vp, i64 key) {
    lj_val *v = (lj_val *)vp;
    if (!v || v->type != LJ_OBJ) return 0;
    for (i64 i = 0; i < v->count; i++)
        if (strcmp(v->keys[i], (const char *)key) == 0) return (i64)v->items[i];
    return 0;
}

i64 lotus_json_get_int(i64 vp, i64 key) {
    lj_val *f = (lj_val *)lotus_json_get(vp, key);
    if (!f) return 0;
    if (f->type == LJ_NUM || f->type == LJ_BOOL) return f->num;
    if (f->type == LJ_STR) return strtoll(f->str, NULL, 10);
    return 0;
}

char *lotus_json_get_string(i64 vp, i64 key) {
    lj_val *f = (lj_val *)lotus_json_get(vp, key);
    if (f && f->type == LJ_STR) return f->str;
    return "";
}

i64 lotus_json_get_bool(i64 vp, i64 key) {
    lj_val *f = (lj_val *)lotus_json_get(vp, key);
    return f && (f->type == LJ_BOOL || f->type == LJ_NUM) && f->num != 0;
}

i64 lotus_json_get_array(i64 vp, i64 key) {
    lj_val *f = (lj_val *)lotus_json_get(vp, key);
    if (f && f->type == LJ_ARR) return (i64)f;
    return 0;
}

i64 lotus_json_array_len(i64 vp) {
    lj_val *v = (lj_val *)vp;
    return (v && v->type == LJ_ARR) ? v->count : 0;
}

i64 lotus_json_array_get(i64 vp, i64 idx) {
    lj_val *v = (lj_val *)vp;
    if (!v || v->type != LJ_ARR || idx < 0 || idx >= v->count) return 0;
    return (i64)v->items[idx];
}

i64 lotus_json_is_null(i64 vp) {
    lj_val *v = (lj_val *)vp;
    return !v || v->type == LJ_NULL;
}

i64 lotus_json_free(i64 vp) {
    lj_free((lj_val *)vp);
    return 0;
}

/* -- stringify -- */

typedef struct {
    char *buf;
    size_t len, cap;
} lj_sb;

static void lj_sb_put(lj_sb *sb, const char *s, size_t n) {
    if (sb->len + n + 1 > sb->cap) {
        while (sb->len + n + 1 > sb->cap) sb->cap *= 2;
        sb->buf = realloc(sb->buf, sb->cap);
    }
    memcpy(sb->buf + sb->len, s, n);
    sb->len += n;
    sb->buf[sb->len] = '\0';
}

static void lj_sb_str(lj_sb *sb, const char *s) { lj_sb_put(sb, s, strlen(s)); }

static void lj_stringify_string(lj_sb *sb, const char *s) {
    lj_sb_str(sb, "\"");
    for (; *s; s++) {
        switch (*s) {
        case '"': lj_sb_str(sb, "\\\""); break;
        case '\\': lj_sb_str(sb, "\\\\"); break;
        case '\n': lj_sb_str(sb, "\\n"); break;
        case '\t': lj_sb_str(sb, "\\t"); break;
        case '\r': lj_sb_str(sb, "\\r"); break;
        default:
            if ((unsigned char)*s < 0x20) {
                char tmp[8];
                snprintf(tmp, sizeof(tmp), "\\u%04x", *s);
                lj_sb_str(sb, tmp);
            } else {
                lj_sb_put(sb, s, 1);
            }
        }
    }
    lj_sb_str(sb, "\"");
}

static void lj_stringify_value(lj_sb *sb, const lj_val *v) {
    if (!v) { lj_sb_str(sb, "null"); return; }
    char tmp[64];
    switch (v->type) {
    case LJ_NULL:
        lj_sb_str(sb, "null");
        break;
    case LJ_BOOL:
        lj_sb_str(sb, v->num ? "true" : "false");
        break;
    case LJ_NUM:
        if (v->dbl == (double)v->num)
            snprintf(tmp, sizeof(tmp), "%lld", (long long)v->num);
        else
            snprintf(tmp, sizeof(tmp), "%.17g", v->dbl);
        lj_sb_str(sb, tmp);
        break;
    case LJ_STR:
        lj_stringify_string(sb, v->str ? v->str : "");
        break;
    case LJ_ARR:
        lj_sb_str(sb, "[");
        for (i64 i = 0; i < v->count; i++) {
            if (i) lj_sb_str(sb, ",");
            lj_stringify_value(sb, v->items[i]);
        }
        lj_sb_str(sb, "]");
        break;
    case LJ_OBJ:
        lj_sb_str(sb, "{");
        for (i64 i = 0; i < v->count; i++) {
            if (i) lj_sb_str(sb, ",");
            lj_stringify_string(sb, v->keys[i]);
            lj_sb_str(sb, ":");
            lj_stringify_value(sb, v->items[i]);
        }
        lj_sb_str(sb, "}");
        break;
    }
}

char *lotus_json_stringify(i64 vp) {
    lj_sb sb;
    sb.cap = 64;
    sb.len = 0;
    sb.buf = malloc(sb.cap);
    sb.buf[0] = '\0';
    lj_stringify_value(&sb, (const lj_val *)vp);
    return sb.buf;
}
