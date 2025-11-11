// Copyright (c) 2024 Nic Barker
// Copyright (c) 2024 go-arena-memory contributors
//
// This software is provided 'as-is', without any express or implied warranty.
// See LICENSE file for full license text.
package Arena

import (
	"testing"
	"unsafe"
)

func TestNewArena(t *testing.T) {
	t.Run("creates arena with valid memory", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if arena == nil {
			t.Fatal("expected arena to be non-nil")
		}
		if arena.Capacity != 1024 {
			t.Errorf("expected capacity 1024, got %d", arena.Capacity)
		}
		if len(arena.Memory) != 1024 {
			t.Errorf("expected memory length 1024, got %d", len(arena.Memory))
		}
		if arena.NextAllocation == 0 && len(memory) > 0 {
			// NextAllocation should be set to cache line alignment padding
			// It might be 0 if already aligned, but should be <= cache line size
			if arena.NextAllocation > 64 {
				t.Errorf("expected NextAllocation <= 64 (cache line size), got %d", arena.NextAllocation)
			}
		}
	})

	t.Run("creates arena with empty memory", func(t *testing.T) {
		memory := make([]byte, 0)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if arena.Capacity != 0 {
			t.Errorf("expected capacity 0, got %d", arena.Capacity)
		}
		if arena.NextAllocation != 0 {
			t.Errorf("expected NextAllocation 0, got %d", arena.NextAllocation)
		}
	})

	t.Run("handles custom cache line size", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, err := NewArena(memory, ArenaWithCacheLineSize(128))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if arena.NextAllocation > 128 {
			t.Errorf("expected NextAllocation <= 128 (custom cache line size), got %d", arena.NextAllocation)
		}
	})

	t.Run("handles arena too small for alignment", func(t *testing.T) {
		// Create a very small memory block that might not fit alignment
		memory := make([]byte, 1)
		// This might fail if alignment requires more than 1 byte
		// The exact behavior depends on the memory address alignment
		arena, err := NewArena(memory)
		// We can't reliably test this without knowing the exact memory address,
		// but we can verify it doesn't panic
		if err != nil && arena == nil {
			// This is acceptable - arena too small for alignment
			return
		}
		if arena != nil && arena.NextAllocation > arena.Capacity {
			t.Errorf("NextAllocation %d exceeds capacity %d", arena.NextAllocation, arena.Capacity)
		}
	})
}

func TestArena_Allocate(t *testing.T) {
	t.Run("allocates memory successfully", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		allocated, err := arena.Allocate(100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(allocated) != 100 {
			t.Errorf("expected allocated size 100, got %d", len(allocated))
		}
		if arena.NextAllocation < 100 {
			t.Errorf("expected NextAllocation >= 100, got %d", arena.NextAllocation)
		}
	})

	t.Run("allocates multiple blocks sequentially", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		initialOffset := arena.NextAllocation

		block1, err1 := arena.Allocate(50)
		if err1 != nil {
			t.Fatalf("expected no error on first allocation, got %v", err1)
		}

		block2, err2 := arena.Allocate(100)
		if err2 != nil {
			t.Fatalf("expected no error on second allocation, got %v", err2)
		}

		if len(block1) != 50 {
			t.Errorf("expected block1 size 50, got %d", len(block1))
		}
		if len(block2) != 100 {
			t.Errorf("expected block2 size 100, got %d", len(block2))
		}

		expectedOffset := initialOffset + 50 + 100
		if arena.NextAllocation != expectedOffset {
			t.Errorf("expected NextAllocation %d, got %d", expectedOffset, arena.NextAllocation)
		}
	})

	t.Run("returns error when capacity exceeded", func(t *testing.T) {
		memory := make([]byte, 100)
		arena, _ := NewArena(memory)

		// Try to allocate more than available
		_, err := arena.Allocate(200)
		if err == nil {
			t.Fatal("expected error when capacity exceeded, got nil")
		}
		if err.Error() != "arena capacity exceeded: cannot allocate required memory" {
			t.Errorf("expected specific error message, got %v", err)
		}
	})

	t.Run("allocates up to capacity limit", func(t *testing.T) {
		memory := make([]byte, 100)
		arena, _ := NewArena(memory)

		initialOffset := arena.NextAllocation
		available := arena.Capacity - initialOffset

		allocated, err := arena.Allocate(available)
		if err != nil {
			t.Fatalf("expected no error allocating up to capacity, got %v", err)
		}
		if len(allocated) != int(available) {
			t.Errorf("expected allocated size %d, got %d", available, len(allocated))
		}

		// Next allocation should fail
		_, err2 := arena.Allocate(1)
		if err2 == nil {
			t.Fatal("expected error when allocating beyond capacity, got nil")
		}
	})

	t.Run("allocated memory is writable", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		allocated, err := arena.Allocate(10)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Write to allocated memory
		for i := range allocated {
			allocated[i] = byte(i)
		}

		// Verify the data was written
		for i := range allocated {
			if allocated[i] != byte(i) {
				t.Errorf("expected allocated[%d] = %d, got %d", i, i, allocated[i])
			}
		}
	})
}

func TestArena_AllocateStruct(t *testing.T) {
	type TestStruct struct {
		X int64
		Y int64
	}

	t.Run("allocates struct successfully", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		ptr, err := AllocateStruct[TestStruct](arena)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}

		// Verify we can write to the struct
		ptr.X = 42
		ptr.Y = 100

		if ptr.X != 42 {
			t.Errorf("expected X = 42, got %d", ptr.X)
		}
		if ptr.Y != 100 {
			t.Errorf("expected Y = 100, got %d", ptr.Y)
		}
	})

	t.Run("allocates multiple structs with proper alignment", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		ptr1, err1 := AllocateStruct[TestStruct](arena)
		if err1 != nil {
			t.Fatalf("expected no error on first allocation, got %v", err1)
		}

		ptr2, err2 := AllocateStruct[TestStruct](arena)
		if err2 != nil {
			t.Fatalf("expected no error on second allocation, got %v", err2)
		}

		// Verify structs are distinct
		if ptr1 == ptr2 {
			t.Fatal("expected different pointers for different allocations")
		}

		// Verify we can write to both independently
		ptr1.X = 1
		ptr1.Y = 2
		ptr2.X = 3
		ptr2.Y = 4

		if ptr1.X != 1 || ptr1.Y != 2 {
			t.Errorf("ptr1 values corrupted: X=%d, Y=%d", ptr1.X, ptr1.Y)
		}
		if ptr2.X != 3 || ptr2.Y != 4 {
			t.Errorf("ptr2 values corrupted: X=%d, Y=%d", ptr2.X, ptr2.Y)
		}
	})

	t.Run("structs are properly aligned", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		ptr, err := AllocateStruct[TestStruct](arena)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check alignment
		address := uintptr(unsafe.Pointer(ptr))
		alignment := unsafe.Alignof(TestStruct{})
		if address%alignment != 0 {
			t.Errorf("struct not properly aligned: address %d, alignment %d", address, alignment)
		}
	})

	t.Run("returns error when capacity exceeded", func(t *testing.T) {
		// Use a size that's large enough to pass NewArena but too small for the struct
		// TestStruct is 16 bytes (2 int64s), plus alignment padding
		memory := make([]byte, 100)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error creating arena, got %v", err)
		}

		// Allocate most of the memory to leave just a small amount
		available := arena.Capacity - arena.NextAllocation
		_, err1 := arena.Allocate(available - 5) // Leave only 5 bytes
		if err1 != nil {
			t.Fatalf("expected no error allocating memory, got %v", err1)
		}

		// Now try to allocate struct - should fail
		_, err2 := AllocateStruct[TestStruct](arena)
		if err2 == nil {
			t.Fatal("expected error when capacity exceeded, got nil")
		}
		if err2.Error() != "arena capacity exceeded: cannot allocate struct" {
			t.Errorf("expected specific error message, got %v", err2)
		}
	})

	t.Run("allocates different struct types", func(t *testing.T) {
		type SmallStruct struct {
			X int8
		}
		type LargeStruct struct {
			X [100]int64
		}

		memory := make([]byte, 2048)
		arena, _ := NewArena(memory)

		small, err1 := AllocateStruct[SmallStruct](arena)
		if err1 != nil {
			t.Fatalf("expected no error allocating SmallStruct, got %v", err1)
		}

		large, err2 := AllocateStruct[LargeStruct](arena)
		if err2 != nil {
			t.Fatalf("expected no error allocating LargeStruct, got %v", err2)
		}

		small.X = 42
		large.X[0] = 100

		if small.X != 42 {
			t.Errorf("expected small.X = 42, got %d", small.X)
		}
		if large.X[0] != 100 {
			t.Errorf("expected large.X[0] = 100, got %d", large.X[0])
		}
	})
}

func TestArena_PersistentEphemeralMemory(t *testing.T) {
	t.Run("initializes persistent memory boundary", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate some memory
		_, err := arena.Allocate(100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Mark persistent boundary
		arena.InitializePersistentMemory()

		if arena.ArenaResetOffset != arena.NextAllocation {
			t.Errorf("expected ArenaResetOffset = NextAllocation (%d), got %d",
				arena.NextAllocation, arena.ArenaResetOffset)
		}
	})

	t.Run("reset ephemeral memory preserves persistent region", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate persistent memory
		persistent, err1 := arena.Allocate(100)
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}

		// Write to persistent memory
		for i := range persistent {
			persistent[i] = byte(i)
		}

		// Mark persistent boundary
		arena.InitializePersistentMemory()
		resetOffset := arena.NextAllocation

		// Allocate ephemeral memory
		ephemeral, err2 := arena.Allocate(50)
		if err2 != nil {
			t.Fatalf("expected no error, got %v", err2)
		}

		// Write to ephemeral memory
		for i := range ephemeral {
			ephemeral[i] = 0xFF
		}

		// Reset ephemeral memory
		arena.ResetEphemeralMemory()

		// Verify NextAllocation was reset
		if arena.NextAllocation != resetOffset {
			t.Errorf("expected NextAllocation = %d after reset, got %d",
				resetOffset, arena.NextAllocation)
		}

		// Verify persistent memory is still intact
		for i := range persistent {
			if persistent[i] != byte(i) {
				t.Errorf("persistent memory corrupted at index %d: expected %d, got %d",
					i, i, persistent[i])
			}
		}
	})

	t.Run("can allocate after reset", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate persistent memory
		_, err1 := arena.Allocate(100)
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}

		arena.InitializePersistentMemory()

		// Allocate and reset ephemeral memory multiple times
		for i := 0; i < 5; i++ {
			_, err2 := arena.Allocate(50)
			if err2 != nil {
				t.Fatalf("expected no error on allocation %d, got %v", i, err2)
			}
			arena.ResetEphemeralMemory()
		}

		// Should still be able to allocate
		_, err3 := arena.Allocate(50)
		if err3 != nil {
			t.Fatalf("expected no error after multiple resets, got %v", err3)
		}
	})

	t.Run("reset before initialization sets to zero", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate some memory
		_, err := arena.Allocate(100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Reset without initializing persistent memory
		arena.ResetEphemeralMemory()

		// Should reset to 0 (or initial alignment offset)
		if arena.NextAllocation != arena.ArenaResetOffset {
			t.Errorf("expected NextAllocation = ArenaResetOffset (%d), got %d",
				arena.ArenaResetOffset, arena.NextAllocation)
		}
		if arena.ArenaResetOffset != 0 {
			t.Errorf("expected ArenaResetOffset = 0 before initialization, got %d",
				arena.ArenaResetOffset)
		}
	})
}

func TestArena_Integration(t *testing.T) {
	t.Run("complex allocation scenario", func(t *testing.T) {
		memory := make([]byte, 2048)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error creating arena, got %v", err)
		}

		// Allocate persistent data
		persistent1, err1 := arena.Allocate(200)
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}
		persistent1[0] = 0xAA

		persistent2, err2 := arena.Allocate(100)
		if err2 != nil {
			t.Fatalf("expected no error, got %v", err2)
		}
		persistent2[0] = 0xBB

		// Mark persistent boundary
		arena.InitializePersistentMemory()

		// Allocate ephemeral data
		ephemeral1, err3 := arena.Allocate(150)
		if err3 != nil {
			t.Fatalf("expected no error, got %v", err3)
		}
		ephemeral1[0] = 0xCC

		// Allocate struct
		type DataStruct struct {
			Value int64
		}
		structPtr, err4 := AllocateStruct[DataStruct](arena)
		if err4 != nil {
			t.Fatalf("expected no error, got %v", err4)
		}
		structPtr.Value = 12345

		// Reset ephemeral
		arena.ResetEphemeralMemory()

		// Verify persistent data intact
		if persistent1[0] != 0xAA {
			t.Error("persistent1 data corrupted")
		}
		if persistent2[0] != 0xBB {
			t.Error("persistent2 data corrupted")
		}

		// Allocate new ephemeral data
		ephemeral2, err5 := arena.Allocate(50)
		if err5 != nil {
			t.Fatalf("expected no error after reset, got %v", err5)
		}
		ephemeral2[0] = 0xDD

		// Verify we can still use the arena
		if arena.NextAllocation < arena.ArenaResetOffset {
			t.Error("NextAllocation should be >= ArenaResetOffset")
		}
	})
}
