package formatters

type Formatter interface {
	Search(...string) (string, error)
	Parse() error
	Format() (Report, error)
}
