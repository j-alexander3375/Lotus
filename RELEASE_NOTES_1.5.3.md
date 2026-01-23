# Release Notes for Lotus v1.5.3

## New Features

### SDL3 Graphics Library Support
- **New**: Added comprehensive SDL3 standard library module (`sdl3`)
  - Window creation and management (SDL3 API: simplified window creation without x, y parameters)
  - Renderer operations (clear, present, draw primitives)
  - Drawing functions: lines, rectangles (outline and filled)
  - Color management
  - Event polling
  - Timing functions (delay, get ticks)
  - Proper cleanup and resource management

- **New**: 15 SDL3 functions available:
  - `sdl3::init(flags)` - Initialize SDL3 subsystems (returns bool)
  - `sdl3::quit()` - Cleanup SDL3
  - `sdl3::create_window(title, w, h, flags)` - Create window (SDL3: no x, y)
  - `sdl3::destroy_window(window*)` - Destroy window
  - `sdl3::create_renderer(window*, flags)` - Create renderer
  - `sdl3::destroy_renderer(renderer*)` - Destroy renderer
  - `sdl3::render_clear(renderer*)` - Clear renderer (returns bool)
  - `sdl3::render_present(renderer*)` - Present rendered frame
  - `sdl3::set_render_draw_color(renderer*, r, g, b, a)` - Set draw color (returns bool)
  - `sdl3::render_draw_line(renderer*, x1, y1, x2, y2)` - Draw line (returns bool)
  - `sdl3::render_draw_rect(renderer*, x, y, w, h)` - Draw rectangle outline (returns bool)
  - `sdl3::render_fill_rect(renderer*, x, y, w, h)` - Draw filled rectangle (returns bool)
  - `sdl3::poll_event(event*)` - Poll for events (returns bool)
  - `sdl3::delay(ms)` - Delay execution
  - `sdl3::get_ticks()` - Get milliseconds since init

- **New**: Test file `tests/test_sdl3.lts` demonstrating SDL3 usage

### Code Generation Improvements
- **Enhanced**: LLVM codegen now supports module-qualified function calls (e.g., `sdl3::init()`)
- **Enhanced**: Improved handling of imported stdlib functions in LLVM backend
- **Fixed**: SDL3 functions properly return bool values (converted to int64 0/1)

### Compiler Updates
- **Updated**: Linker now includes `-lSDL3` flag for SDL3 library linking
- **Updated**: SDL3 function declarations match SDL3 API (bool returns, simplified signatures)

## API Changes

### SDL3 Module
- All SDL3 render functions now return `bool` (SDL3 API change from SDL2)
- `create_window` signature changed: removed x, y parameters (SDL3 simplification)
- `init` returns `bool` instead of `int` (SDL3 API change)

## Technical Details

- SDL3 module fully integrated into standard library system
- Both LLVM and GCC backends support SDL3 (GCC uses stubs)
- SDL_Rect structures created on stack for rect drawing functions
- Proper type conversions (int64 → int32, int64 → pointers, bool → int64)

## Dependencies

- Requires SDL3 development libraries for compilation
- Linker flag: `-lSDL3`

## Example Usage

```lotus
use "sdl3";
use "io";

fn int main() {
    int init_result = sdl3::init(0x00000020); // SDL_INIT_VIDEO
    if (init_result == 0) {
        printf("Failed to initialize SDL3\n");
        ret 1;
    }
    
    int window = sdl3::create_window("My Game", 640, 480, 0x00000020);
    int renderer = sdl3::create_renderer(window, 0x00000002);
    
    sdl3::set_render_draw_color(renderer, 255, 0, 0, 255);
    sdl3::render_clear(renderer);
    sdl3::render_present(renderer);
    sdl3::delay(2000);
    
    sdl3::destroy_renderer(renderer);
    sdl3::destroy_window(window);
    sdl3::quit();
    ret 0;
}
```
