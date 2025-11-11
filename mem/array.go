package mem

type MemArray[T any] struct {
	Capacity      int32
	Length        int32
	InternalArray []T
}

func NewMemArray[T any](capacity int32) MemArray[T] {
	return MemArray[T]{
		Capacity:      capacity,
		Length:        0,
		InternalArray: make([]T, capacity),
	}
}

func rangeCheck(index int32, length int32) bool {
	return index < length && index >= 0
}

func MemArray_Get[T any](array *MemArray[T], index int32) *T {
	if !rangeCheck(index, int32(len(array.InternalArray))) {
		return nil
	}
	return &array.InternalArray[index]
}
func MArray_GetValue[T any](array *MemArray[T], index int32) T {
	zero := new(T)
	if !rangeCheck(index, int32(len(array.InternalArray))) {
		return *zero
	}
	return array.InternalArray[index]
}

func MArray_Add[T any](array *MemArray[T], item T) *T {
	if array.Length == array.Capacity-1 {
		return nil
	}
	array.InternalArray[array.Length] = item
	array.Length++
	return &array.InternalArray[array.Length-1]
}

func MArray_Set[T any](array *MemArray[T], index int32, item T) {
	if index < 0 || index >= int32(len(array.InternalArray)) {
		return
	}
	array.InternalArray[index] = item
}

func MArray_RemoveSwapback[T any](array *MemArray[T], index int32) T {
	zero := new(T)
	if !rangeCheck(index, array.Length) {
		return *zero
	}
	array.Length--
	removed := array.InternalArray[index]
	array.InternalArray[index] = array.InternalArray[array.Length]
	return removed
}
