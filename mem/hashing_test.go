package mem

import (
	"testing"
)

func TestNewHashBuilder(t *testing.T) {
	t.Run("creates builder with seed", func(t *testing.T) {
		seed := uint32(42)
		builder := NewHashBuilder(seed)

		if builder == nil {
			t.Fatal("expected builder to be non-nil")
		}
		if builder.hash != seed {
			t.Errorf("expected hash = %d, got %d", seed, builder.hash)
		}
		if builder.stringId != "" {
			t.Errorf("expected empty stringId, got %q", builder.stringId)
		}
	})

	t.Run("creates builder with zero seed", func(t *testing.T) {
		builder := NewHashBuilder(0)
		if builder.hash != 0 {
			t.Errorf("expected hash = 0, got %d", builder.hash)
		}
	})
}

func TestHashBuilder_AddByte(t *testing.T) {
	t.Run("adds single byte", func(t *testing.T) {
		builder := NewHashBuilder(0)
		result := builder.AddByte(65) // 'A'

		if result != builder {
			t.Error("expected AddByte to return the builder for chaining")
		}
		if builder.hash == 0 {
			t.Error("expected hash to be modified after adding byte")
		}
	})

	t.Run("adds multiple bytes sequentially", func(t *testing.T) {
		builder1 := NewHashBuilder(0)
		builder1.AddByte(65).AddByte(66).AddByte(67)

		builder2 := NewHashBuilder(0)
		builder2.AddByte(65)
		builder2.AddByte(66)
		builder2.AddByte(67)

		if builder1.hash != builder2.hash {
			t.Errorf("expected same hash for sequential adds, got %d vs %d", builder1.hash, builder2.hash)
		}
	})

	t.Run("produces different hashes for different bytes", func(t *testing.T) {
		builder1 := NewHashBuilder(0)
		builder1.AddByte(65)

		builder2 := NewHashBuilder(0)
		builder2.AddByte(66)

		if builder1.hash == builder2.hash {
			t.Error("expected different hashes for different bytes")
		}
	})
}

func TestHashBuilder_AddBytes(t *testing.T) {
	t.Run("adds bytes from slice", func(t *testing.T) {
		builder := NewHashBuilder(0)
		data := []byte{65, 66, 67}
		builder.AddBytes(data, 3)

		if builder.hash == 0 {
			t.Error("expected hash to be modified after adding bytes")
		}
	})

	t.Run("adds partial bytes", func(t *testing.T) {
		builder1 := NewHashBuilder(0)
		data := []byte{65, 66, 67, 68}
		builder1.AddBytes(data, 2)

		builder2 := NewHashBuilder(0)
		builder2.AddByte(65).AddByte(66)

		if builder1.hash != builder2.hash {
			t.Errorf("expected same hash for partial bytes, got %d vs %d", builder1.hash, builder2.hash)
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		builder := NewHashBuilder(42)
		initialHash := builder.hash
		data := []byte{}
		builder.AddBytes(data, 0)

		if builder.hash != initialHash {
			t.Errorf("expected hash to remain unchanged for empty slice, got %d vs %d", builder.hash, initialHash)
		}
	})
}

func TestHashBuilder_AddString(t *testing.T) {
	t.Run("adds string to hash", func(t *testing.T) {
		builder := NewHashBuilder(0)
		result := builder.AddString("test")

		if result != builder {
			t.Error("expected AddString to return the builder for chaining")
		}
		if builder.hash == 0 {
			t.Error("expected hash to be modified after adding string")
		}
		if builder.stringId != "test" {
			t.Errorf("expected stringId = %q, got %q", "test", builder.stringId)
		}
	})

	t.Run("adds multiple strings", func(t *testing.T) {
		builder := NewHashBuilder(0)
		builder.AddString("hello").AddString("world")

		if builder.stringId != "helloworld" {
			t.Errorf("expected stringId = %q, got %q", "helloworld", builder.stringId)
		}
	})

	t.Run("produces different hashes for different strings", func(t *testing.T) {
		builder1 := NewHashBuilder(0)
		builder1.AddString("test1")

		builder2 := NewHashBuilder(0)
		builder2.AddString("test2")

		if builder1.hash == builder2.hash {
			t.Error("expected different hashes for different strings")
		}
	})

	t.Run("handles empty string", func(t *testing.T) {
		builder := NewHashBuilder(42)
		initialHash := builder.hash
		builder.AddString("")

		if builder.hash != initialHash {
			t.Errorf("expected hash to remain unchanged for empty string, got %d vs %d", builder.hash, initialHash)
		}
		if builder.stringId != "" {
			t.Errorf("expected empty stringId, got %q", builder.stringId)
		}
	})

	t.Run("uses custom joiner", func(t *testing.T) {
		builder := NewHashBuilder(0)
		customJoiner := func(a, b string) string {
			return a + "-" + b
		}
		option := func(opts *HashingOptions) {
			opts.StringIdJoiner = customJoiner
		}
		builder.AddString("hello", option).AddString("world", option)

		// First call: joiner("", "hello") = "-hello"
		// Second call: joiner("-hello", "world") = "-hello-world"
		if builder.stringId != "-hello-world" {
			t.Errorf("expected stringId = %q, got %q", "-hello-world", builder.stringId)
		}
	})
}

func TestHashBuilder_AddNumber(t *testing.T) {
	t.Run("adds number to hash", func(t *testing.T) {
		builder := NewHashBuilder(0)
		result := builder.AddNumber(42)

		if result != builder {
			t.Error("expected AddNumber to return the builder for chaining")
		}
		if builder.hash == 0 {
			t.Error("expected hash to be modified after adding number")
		}
		if builder.stringId != "42" {
			t.Errorf("expected stringId = %q, got %q", "42", builder.stringId)
		}
	})

	t.Run("adds multiple numbers", func(t *testing.T) {
		builder := NewHashBuilder(0)
		builder.AddNumber(1).AddNumber(2).AddNumber(3)

		if builder.stringId != "123" {
			t.Errorf("expected stringId = %q, got %q", "123", builder.stringId)
		}
	})

	t.Run("produces different hashes for different numbers", func(t *testing.T) {
		builder1 := NewHashBuilder(0)
		builder1.AddNumber(1)

		builder2 := NewHashBuilder(0)
		builder2.AddNumber(2)

		if builder1.hash == builder2.hash {
			t.Error("expected different hashes for different numbers")
		}
	})

	t.Run("handles zero", func(t *testing.T) {
		builder := NewHashBuilder(0)
		builder.AddNumber(0)

		if builder.stringId != "0" {
			t.Errorf("expected stringId = %q, got %q", "0", builder.stringId)
		}
		if builder.hash == 0 {
			t.Error("expected hash to be modified even for zero")
		}
	})

	t.Run("uses custom joiner", func(t *testing.T) {
		builder := NewHashBuilder(0)
		customJoiner := func(a, b string) string {
			return a + "|" + b
		}
		option := func(opts *HashingOptions) {
			opts.StringIdJoiner = customJoiner
		}
		builder.AddNumber(1, option).AddNumber(2, option)

		// First call: joiner("", "1") = "|1"
		// Second call: joiner("|1", "2") = "|1|2"
		if builder.stringId != "|1|2" {
			t.Errorf("expected stringId = %q, got %q", "|1|2", builder.stringId)
		}
	})
}

func TestHashBuilder_AddNumbers(t *testing.T) {
	t.Run("adds multiple numbers", func(t *testing.T) {
		builder := NewHashBuilder(0)
		numbers := []uint32{1, 2, 3, 4, 5}
		result := builder.AddNumbers(numbers)

		if result != builder {
			t.Error("expected AddNumbers to return the builder for chaining")
		}
		if builder.stringId != "12345" {
			t.Errorf("expected stringId = %q, got %q", "12345", builder.stringId)
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		builder := NewHashBuilder(42)
		initialHash := builder.hash
		initialStringId := builder.stringId
		builder.AddNumbers([]uint32{})

		if builder.hash != initialHash {
			t.Errorf("expected hash to remain unchanged for empty slice, got %d vs %d", builder.hash, initialHash)
		}
		if builder.stringId != initialStringId {
			t.Errorf("expected stringId to remain unchanged, got %q vs %q", builder.stringId, initialStringId)
		}
	})

	t.Run("produces same result as sequential AddNumber calls", func(t *testing.T) {
		numbers := []uint32{10, 20, 30}

		builder1 := NewHashBuilder(0)
		builder1.AddNumbers(numbers)

		builder2 := NewHashBuilder(0)
		builder2.AddNumber(10).AddNumber(20).AddNumber(30)

		if builder1.hash != builder2.hash {
			t.Errorf("expected same hash, got %d vs %d", builder1.hash, builder2.hash)
		}
		if builder1.stringId != builder2.stringId {
			t.Errorf("expected same stringId, got %q vs %q", builder1.stringId, builder2.stringId)
		}
	})

	t.Run("uses custom joiner", func(t *testing.T) {
		builder := NewHashBuilder(0)
		customJoiner := func(a, b string) string {
			return a + "," + b
		}
		option := func(opts *HashingOptions) {
			opts.StringIdJoiner = customJoiner
		}
		builder.AddNumbers([]uint32{1, 2, 3}, option)

		// First: joiner("", "1") = ",1"
		// Second: joiner(",1", "2") = ",1,2"
		// Third: joiner(",1,2", "3") = ",1,2,3"
		if builder.stringId != ",1,2,3" {
			t.Errorf("expected stringId = %q, got %q", ",1,2,3", builder.stringId)
		}
	})
}

func TestHashBuilder_build(t *testing.T) {
	t.Run("builds HashElementId", func(t *testing.T) {
		builder := NewHashBuilder(0)
		builder.AddString("test")
		result := builder.build()

		if result.Id == 0 {
			t.Error("expected non-zero Id")
		}
		if result.Offset != 0 {
			t.Errorf("expected Offset = 0, got %d", result.Offset)
		}
		if result.BaseId != result.Id {
			t.Errorf("expected BaseId = Id (%d), got %d", result.Id, result.BaseId)
		}
		if result.StringId != "test" {
			t.Errorf("expected StringId = %q, got %q", "test", result.StringId)
		}
	})

	t.Run("produces consistent results", func(t *testing.T) {
		builder1 := NewHashBuilder(0)
		builder1.AddString("test")
		result1 := builder1.build()

		builder2 := NewHashBuilder(0)
		builder2.AddString("test")
		result2 := builder2.build()

		if result1.Id != result2.Id {
			t.Errorf("expected consistent Id, got %d vs %d", result1.Id, result2.Id)
		}
		if result1.StringId != result2.StringId {
			t.Errorf("expected consistent StringId, got %q vs %q", result1.StringId, result2.StringId)
		}
	})

	t.Run("Id and BaseId are set correctly", func(t *testing.T) {
		builder := NewHashBuilder(0)
		builder.AddString("test")
		result := builder.build()

		if result.Id == 0 {
			t.Error("expected non-zero Id")
		}
		if result.BaseId != result.Id {
			t.Errorf("expected BaseId = Id, got %d vs %d", result.BaseId, result.Id)
		}
		if result.Id != builder.hash+1 {
			t.Errorf("expected Id = hash + 1, got %d (hash was %d)", result.Id, builder.hash)
		}
	})
}

func TestHashString(t *testing.T) {
	t.Run("hashes string with seed", func(t *testing.T) {
		result := HashString("test", 0)

		if result.Id == 0 {
			t.Error("expected non-zero Id")
		}
		if result.StringId != "test" {
			t.Errorf("expected StringId = %q, got %q", "test", result.StringId)
		}
	})

	t.Run("produces consistent results", func(t *testing.T) {
		result1 := HashString("hello", 42)
		result2 := HashString("hello", 42)

		if result1.Id != result2.Id {
			t.Errorf("expected consistent Id, got %d vs %d", result1.Id, result2.Id)
		}
		if result1.StringId != result2.StringId {
			t.Errorf("expected consistent StringId, got %q vs %q", result1.StringId, result2.StringId)
		}
	})

	t.Run("produces different results for different seeds", func(t *testing.T) {
		result1 := HashString("test", 0)
		result2 := HashString("test", 1)

		if result1.Id == result2.Id {
			t.Error("expected different Ids for different seeds")
		}
	})

	t.Run("produces different results for different strings", func(t *testing.T) {
		result1 := HashString("test1", 0)
		result2 := HashString("test2", 0)

		if result1.Id == result2.Id {
			t.Error("expected different Ids for different strings")
		}
	})

	t.Run("uses custom joiner", func(t *testing.T) {
		customJoiner := func(a, b string) string {
			return a + "::" + b
		}
		option := func(opts *HashingOptions) {
			opts.StringIdJoiner = customJoiner
		}
		result := HashString("test", 0, option)

		// Single string: joiner("", "test") = "::test"
		if result.StringId != "::test" {
			t.Errorf("expected StringId = %q, got %q", "::test", result.StringId)
		}
	})
}

func TestHashNumber(t *testing.T) {
	t.Run("hashes number with seed", func(t *testing.T) {
		result := HashNumber(42, 0)

		if result.Id == 0 {
			t.Error("expected non-zero Id")
		}
		if result.StringId != "42" {
			t.Errorf("expected StringId = %q, got %q", "42", result.StringId)
		}
	})

	t.Run("produces consistent results", func(t *testing.T) {
		result1 := HashNumber(100, 10)
		result2 := HashNumber(100, 10)

		if result1.Id != result2.Id {
			t.Errorf("expected consistent Id, got %d vs %d", result1.Id, result2.Id)
		}
		if result1.StringId != result2.StringId {
			t.Errorf("expected consistent StringId, got %q vs %q", result1.StringId, result2.StringId)
		}
	})

	t.Run("produces different results for different numbers", func(t *testing.T) {
		result1 := HashNumber(1, 0)
		result2 := HashNumber(2, 0)

		if result1.Id == result2.Id {
			t.Error("expected different Ids for different numbers")
		}
	})

	t.Run("handles zero", func(t *testing.T) {
		result := HashNumber(0, 0)

		if result.StringId != "0" {
			t.Errorf("expected StringId = %q, got %q", "0", result.StringId)
		}
		if result.Id == 0 {
			t.Error("expected non-zero Id even for zero input")
		}
	})
}

func TestHashManyNumbers(t *testing.T) {
	t.Run("hashes multiple numbers", func(t *testing.T) {
		numbers := []uint32{1, 2, 3, 4, 5}
		result := HashManyNumbers(0, numbers)

		if result.Id == 0 {
			t.Error("expected non-zero Id")
		}
		if result.StringId != "12345" {
			t.Errorf("expected StringId = %q, got %q", "12345", result.StringId)
		}
	})

	t.Run("produces consistent results", func(t *testing.T) {
		numbers := []uint32{10, 20, 30}
		result1 := HashManyNumbers(5, numbers)
		result2 := HashManyNumbers(5, numbers)

		if result1.Id != result2.Id {
			t.Errorf("expected consistent Id, got %d vs %d", result1.Id, result2.Id)
		}
		if result1.StringId != result2.StringId {
			t.Errorf("expected consistent StringId, got %q vs %q", result1.StringId, result2.StringId)
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		result := HashManyNumbers(42, []uint32{})

		if result.Id == 0 {
			t.Error("expected non-zero Id even for empty slice")
		}
		if result.StringId != "" {
			t.Errorf("expected empty StringId, got %q", result.StringId)
		}
	})

	t.Run("produces same result as sequential HashNumber calls", func(t *testing.T) {
		numbers := []uint32{5, 10, 15}

		builder1 := NewHashBuilder(0)
		builder1.AddNumbers(numbers)
		result1 := builder1.build()

		result2 := HashManyNumbers(0, numbers)

		if result1.Id != result2.Id {
			t.Errorf("expected same Id, got %d vs %d", result1.Id, result2.Id)
		}
		if result1.StringId != result2.StringId {
			t.Errorf("expected same StringId, got %q vs %q", result1.StringId, result2.StringId)
		}
	})
}

func TestHashingOptionsWithJoiner(t *testing.T) {
	t.Run("creates options with custom joiner", func(t *testing.T) {
		customJoiner := func(a, b string) string {
			return a + "|" + b
		}
		opts := HashingOptionsWithJoiner(customJoiner)

		if opts.StringIdJoiner == nil {
			t.Fatal("expected StringIdJoiner to be set")
		}

		result := opts.StringIdJoiner("a", "b")
		if result != "a|b" {
			t.Errorf("expected joiner result = %q, got %q", "a|b", result)
		}
	})
}

func TestDefaultHashingOptions(t *testing.T) {
	t.Run("has default joiner that concatenates", func(t *testing.T) {
		result := DefaultHashingOptions.StringIdJoiner("hello", "world")
		if result != "helloworld" {
			t.Errorf("expected default joiner to concatenate, got %q", result)
		}
	})
}
