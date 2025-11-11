# Go Arena Memory

A high-performance, bump-pointer memory management implementation in Go, inspired by the C `Clay_Arena` from the [Clay](https://github.com/nicbarker/clay) library. This implementation provides a way to bypass Go's garbage collector for transient data, ensuring fast, predictable memory operations with O(1) frame reset capabilities.

## Overview

This Go implementation translates the high-performance memory management model of the C `Clay_Arena` into idiomatic Go. It's designed for performance-critical Go libraries that need to manage transient data efficiently, especially in immediate-mode UI frameworks where data is rebuilt every frame.

## Key Features

- **Bump-Pointer Allocation**: Fast, O(1) allocation from a pre-allocated memory block
- **Automatic Memory Alignment**: Handles CPU alignment requirements automatically
- **Configurable Options**: Customizable cache line alignment via functional options pattern
- **Dual Memory Regions**: Supports both persistent and ephemeral memory regions
- **O(1) Frame Reset**: Instantly "free" all ephemeral memory with a single pointer assignment
- **Error Handling**: Returns errors instead of panicking for graceful error handling
- **GC-Friendly**: Minimizes garbage collection pressure by using a single pre-allocated memory block

## The Arena Structure

The `Arena` struct manages a contiguous block of memory:

```go
type Arena struct {
    NextAllocation   uintptr  // Offset where the next allocation will begin
    Capacity         uintptr  // Total size of the memory block
    Memory           []byte   // The contiguous memory block
    ArenaResetOffset uintptr  // Boundary between persistent and ephemeral memory
}
```

## Usage

### Creating an Arena

```go
import "github.com/zodimo/go-arena-memory/mem"

// Allocate a memory block (e.g., 1MB)
memory := make([]byte, 1024*1024)

// Create a new arena with default options (64-byte cache line alignment)
arena, err := mem.NewArena(memory)
if err != nil {
    log.Fatal(err)
}

// Or create an arena with custom cache line size
arena, err := mem.NewArena(memory, mem.ArenaWithCacheLineSize(128))
if err != nil {
    log.Fatal(err)
}
```

### Allocating Raw Memory

The `Allocate` method allocates a raw byte slice of the specified size:

```go
// Allocate 1024 bytes
data, err := arena.Allocate(1024)
if err != nil {
    log.Fatal(err)
}

// Use the allocated memory
copy(data, []byte("Hello, Arena!"))
```

### Allocating Structs

The `AllocateStruct` function allocates space for any type with proper alignment:

```go
type MyStruct struct {
    ID    int64
    Value float64
}

// Allocate a struct from the arena
myStruct, err := mem.AllocateStruct[MyStruct](arena)
if err != nil {
    log.Fatal(err)
}

// Use the allocated struct
myStruct.ID = 42
myStruct.Value = 3.14
```

### Persistent and Ephemeral Memory

The arena supports two memory regions:

**Persistent Memory**: Data that persists across frames (e.g., hash maps, caches)

```go
// After allocating persistent structures, mark the boundary
arena.InitializePersistentMemory()
```

**Ephemeral Memory**: Transient data rebuilt every frame (e.g., layout buffers, render commands)

```go
// At the start of each frame, reset ephemeral memory (O(1) operation)
arena.ResetEphemeralMemory()
```

## How It Works

### Memory Alignment

The implementation automatically handles CPU alignment requirements. When allocating a struct, it:

1. Calculates the type's size and alignment requirements
2. Determines the required padding to align the current address
3. Bumps the allocation pointer forward (padding + size)
4. Returns a typed pointer to the aligned memory location

### Bump-Pointer Allocation

All allocations are sequential within the pre-allocated block. The `NextAllocation` pointer tracks the current position and is simply incremented for each allocation, making it extremely fast.

### O(1) Reset

The ephemeral memory region can be instantly "freed" by resetting `NextAllocation` back to `ArenaResetOffset`. This avoids individual deallocations and eliminates fragmentation.

## Why Use This?

- **Performance**: Predictable, near-zero overhead memory operations
- **No GC Pressure**: The pre-allocated block is rarely inspected by the Go GC
- **Immediate-Mode Friendly**: Perfect for frameworks that rebuild state every frame
- **Type Safety**: Uses Go generics to provide type-safe allocation

## Configuration Options

The arena supports configurable options using the functional options pattern:

- **Cache Line Size**: Control the initial cache line alignment (default: 64 bytes)
  ```go
  arena, err := mem.NewArena(memory, mem.ArenaWithCacheLineSize(128))
  ```

## Implementation Details

This implementation uses Go's `unsafe` package to:
- Calculate physical memory addresses for alignment
- Perform type-casting from raw memory addresses to typed pointers
- Handle memory arithmetic similar to C-style pointer operations

The use of `unsafe` is limited to these essential operations and is necessary to achieve the performance characteristics of the original C implementation.

### Error Handling

Both `Allocate` and `AllocateStruct` return errors instead of panicking when:
- The arena capacity is exceeded
- There's insufficient space for initial alignment

This allows for graceful error handling in production code.

## Attribution

This project is inspired by and based on the [Clay](https://github.com/nicbarker/clay) library's `Clay_Arena` implementation. The original C implementation provides the foundation for the memory management concepts used here.

**Original Repository**: [nicbarker/clay](https://github.com/nicbarker/clay)

## License

This project is licensed under the zlib/libpng License, matching the license of the original [Clay](https://github.com/nicbarker/clay) library that inspired this implementation.

See the [LICENSE](LICENSE) file for the full license text.
