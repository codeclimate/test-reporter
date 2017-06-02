package formatters

// Formatter needs to be implemented for each new test system
type Formatter interface {
	// Search for the both the "standard" paths for the formatter,
	// plus any additional paths, for a file that can be parsed
	// by the formatter.
	Search(...string) (string, error)
	// Format the information for Parse into a standardized "Report".
	// Returns an error if there was a problem formatting the results.
	Format() (Report, error)
}
