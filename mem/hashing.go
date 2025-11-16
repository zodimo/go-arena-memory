package mem

type HashElementId struct {
	Id       uint32
	Offset   uint32
	BaseId   uint32
	StringId string // to recover the string from the hash
}

type HashBuilder struct {
	hash     uint32
	stringId string
}

func NewHashBuilder(seed uint32) *HashBuilder {
	return &HashBuilder{hash: seed, stringId: ""}
}

func (h *HashBuilder) AddBytes(data []byte, length int32) {
	for _, charByte := range data[:length] {
		h.AddByte(charByte)
	}
}
func (h *HashBuilder) AddByte(data byte) *HashBuilder {
	h.hash += uint32(data)
	h.hash += (h.hash << 10)
	h.hash ^= (h.hash >> 6)
	return h
}
func (h *HashBuilder) AddString(key string) *HashBuilder {
	stringBytes := []byte(key)
	for _, charByte := range stringBytes {
		h.AddByte(charByte)
	}
	h.stringId += key
	return h
}
func (h *HashBuilder) AddNumber(number uint32) *HashBuilder {
	h.hash += (number + 48)
	h.hash += (h.hash << 10)
	h.hash ^= (h.hash >> 6)
	return h
}

func (h *HashBuilder) AddNumbers(numbers ...uint32) *HashBuilder {
	for _, number := range numbers {
		h.AddNumber(number)
	}
	return h
}

func (h *HashBuilder) build() HashElementId {

	hash := h.hash
	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)

	return HashElementId{
		Id:       h.hash + 1,
		Offset:   0,
		BaseId:   h.hash + 1,
		StringId: h.stringId,
	}
}

func HashString(key string, seed uint32) HashElementId {
	return NewHashBuilder(seed).AddString(key).build()
}

func HashNumber(number uint32, seed uint32) HashElementId {
	return NewHashBuilder(seed).AddNumber(number).build()
}

func HashManyNumbers(seed uint32, numbers ...uint32) HashElementId {
	return NewHashBuilder(seed).AddNumbers(numbers...).build()
}
