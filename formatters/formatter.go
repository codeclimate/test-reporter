package formatters

type Formatter interface {
	Parse() error
	Format() (Report, error)
}
