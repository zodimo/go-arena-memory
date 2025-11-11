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

// typeName arrayName##_RemoveSwapback(arrayName *array, int32_t index) {                                          \
// 	if (Clay__Array_RangeCheck(index, array->length)) {                                                         \
// 		array->length--;                                                                                        \
// 		typeName removed = array->internalArray[index];                                                         \
// 		array->internalArray[index] = array->internalArray[array->length];                                      \
// 		return removed;                                                                                         \
// 	}                                                                                                           \
// 	return typeName##_DEFAULT;                                                                                  \
// }

// // The below functions define array bounds checking and convenience functions for a provided type.
// #define CLAY__ARRAY_DEFINE_FUNCTIONS(typeName, arrayName)                                                       \
//                                                                                                                 \
// typedef struct                                                                                                  \
// {                                                                                                               \
//     int32_t length;                                                                                             \
//     typeName *internalArray;                                                                                    \
// } arrayName##Slice;                                                                                             \
//                                                                                                                 \
// typeName typeName##_DEFAULT = CLAY__DEFAULT_STRUCT;                                                             \
//                                                                                                                 \
// arrayName arrayName##_Allocate_Arena(int32_t capacity, Clay_Arena *arena) {                                     \
//     return CLAY__INIT(arrayName){.capacity = capacity, .length = 0,                                             \
//         .internalArray = (typeName *)Clay__Array_Allocate_Arena(capacity, sizeof(typeName), arena)};            \
// }                                                                                                               \
//                                                                                                                 \
// typeName *arrayName##_Get(arrayName *array, int32_t index) {                                                    \
//     return Clay__Array_RangeCheck(index, array->length) ? &array->internalArray[index] : &typeName##_DEFAULT;   \
// }                                                                                                               \
//                                                                                                                 \
// typeName arrayName##_GetValue(arrayName *array, int32_t index) {                                                \
//     return Clay__Array_RangeCheck(index, array->length) ? array->internalArray[index] : typeName##_DEFAULT;     \
// }                                                                                                               \
//                                                                                                                 \
// typeName *arrayName##_Add(arrayName *array, typeName item) {                                                    \
//     if (Clay__Array_AddCapacityCheck(array->length, array->capacity)) {                                         \
//         array->internalArray[array->length++] = item;                                                           \
//         return &array->internalArray[array->length - 1];                                                        \
//     }                                                                                                           \
//     return &typeName##_DEFAULT;                                                                                 \
// }                                                                                                               \
//                                                                                                                 \
// typeName *arrayName##Slice_Get(arrayName##Slice *slice, int32_t index) {                                        \
//     return Clay__Array_RangeCheck(index, slice->length) ? &slice->internalArray[index] : &typeName##_DEFAULT;   \
// }                                                                                                               \
//                                                                                                                 \
//   \
//                                                                                                                 \
// void arrayName##_Set(arrayName *array, int32_t index, typeName value) {                                         \
// 	if (Clay__Array_RangeCheck(index, array->capacity)) {                                                       \
// 		array->internalArray[index] = value;                                                                    \
// 		array->length = index < array->length ? array->length : index + 1;                                      \
// 	}                                                                                                           \
// }                                                                                                               \
