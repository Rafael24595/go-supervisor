package compact

type option func(*config)

type config struct {
	caller    string
	formatter Formatter
	processor Processor
}

func newConfig(formatter Formatter, options ...option) config {
	cfg := defaultConfig(formatter)
	for _, o := range options {
		o(&cfg)
	}
	return cfg
}

func defaultConfig(formatter Formatter) config {
	return config{
		caller:    defaultCaller,
		formatter: formatter,
		processor: EmptyProcessor(),
	}
}

// WithProcessor configures the Processor used to transform frame paths
// before they are formatted.
func WithProcessor(processor Processor) option {
	return func(c *config) {
		c.processor = processor
	}
}

func withCaller(caller string) option {
	return func(c *config) {
		c.caller = caller
	}
}
