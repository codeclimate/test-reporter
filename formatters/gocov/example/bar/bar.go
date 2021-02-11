package bar

type Bar struct{}

func New() *Bar {
	return &Bar{}
}

func (b *Bar) String() string {
	return "bar"
}

func (b *Bar) Bytes() []byte {
	return []byte("bar")
}
