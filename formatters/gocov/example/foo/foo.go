package foo

import "fmt"

type Foo struct {
	anotherFoo interface{}
}

func New(foo interface{}) *Foo {
	return &Foo{
		anotherFoo: foo,
	}
}

func (f *Foo) String() string {
	i, ok := f.anotherFoo.(interface{ String() string })
	if ok {
		return fmt.Sprintf("foo %s", i)
	}
	return "foo"
}
