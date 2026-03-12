# Release Notes for Lotus 1.8.4

## Highlights

This release adds opt-in unqualified module imports, a complete LLVM HTTP module,
full iterative regex replace-all, UDP networking primitives, and an HTTP connection
pool — all backed by complete LLVM IR implementations with no stubs.

## New Features

### Language: Unqualified Module Imports

Modules imported with `use "module"` now export their functions into the file
scope so that the module qualifier becomes optional:

```lotus
use "regex";
use "http";

fn int main() {
    // unqualified – equivalent to regex::replace_all
    string out = replace_all("foo bar foo", "foo", "baz");

    // qualified form still works as before
    int status = http::parse_status(resp, resp_len);
    ret 0;
}
```

Functions whose bare name collides with a built-in continue to work only through
the qualified form; unqualified resolution never shadows existing builtins.

### LLVM Backend: Full HTTP Module

The LLVM backend now has complete, stub-free implementations for every `http::`
function:

| Function | Description |
|---|---|
| `http::get(fd, path, host)` | Write HTTP GET + read response |
| `http::post(fd, path, host, body, body_len)` | Write HTTP POST + read response |
| `http::parse_status(response, len)` | Extract HTTP status code |
| `http::get_header(response, len, name, buf)` | Extract a named header value |
| `http::get_body(response, len)` | Pointer to response body |
| `http::parse_headers(response, len, buf)` | Count response headers |
| `http::pool_new(capacity)` | Allocate a connection pool |
| `http::pool_get(pool)` | Check out a connection fd |
| `http::pool_put(pool, fd)` | Return a connection fd |
| `http::pool_close(pool)` | Destroy a connection pool |

### LLVM Backend: Full `regex::replace_all`

`regex::replace_all(input, pattern, replacement)` now runs a complete iterative
replacement loop over the input string using `strstr`/`malloc`/`memcpy` — all
occurrences of `pattern` are replaced without any stub delegation.

### LLVM Backend: UDP Networking

Two new low-level network functions are available:

- `net::sendto_ipv4(fd, buf, len, ip, port)` — send a UDP datagram to an IPv4 address
- `net::recvfrom(fd, buf, len)` — receive a UDP datagram, returns byte count

## Bug Fixes

- Fixed dispatch regression where `http::pool_*` functions were silently dropped
  after per-module dispatch prefixes were introduced.
- Fixed LLVM IR *"terminator found in middle of basic block"* in
  `generateHTTPGetHeader` — all `alloca` instructions are now hoisted above the
  first branch.
- Fixed LLVM IR *"void value has name"* error in `generateRegexReplaceAll` free call.
- Centralised the per-module dispatch-prefix logic into a single
  `moduleDispatchPrefix` helper, removing a duplicated inline switch.

## Testing

All Go unit and integration tests pass (`go test -count=1 ./...`).

New integration test files:
- `tests/http_pool_capabilities.lts` — validates all four pool operations
- `tests/http_regex_capabilities.lts` — validates HTTP parsing helpers + `replace_all`

## Installation

### From Source
```bash
git clone https://github.com/j-alexander3375/Lotus
cd Lotus
git checkout 1.8.4
cd src
go build -o ../lotus .
```

### Arch Linux (AUR)
```bash
yay -S lotus-lang
```

### Binary Downloads
Download pre-built binaries from the [releases page](https://github.com/j-alexander3375/Lotus/releases/tag/1.8.4).

## Checksums

SHA256 checksums for release artifacts are available in the release assets.
