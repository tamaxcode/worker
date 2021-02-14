package worker

import "runtime"

// Config represent configuration setting for collector to send received task to worker
type Config struct {
	// Number of worker spawn to handle concurrent work. Default to "runtime.NumCPU"
	// giving value <= 0 will fallback to default value
	NoOfWorkers *int
}

// SetNoOfWorkers set numOfWorkers value
func (c *Config) SetNoOfWorkers(num int) {
	c.NoOfWorkers = &num
}

func mergeOrDefault(cfgs []*Config) *Config {
	cfg := &Config{}

	for _, c := range cfgs {
		if c == nil {
			continue
		}

		if c.NoOfWorkers != nil && *c.NoOfWorkers > 0 {
			cfg.NoOfWorkers = c.NoOfWorkers
		}
	}

	if cfg.NoOfWorkers == nil {
		cpu := runtime.NumCPU()
		cfg.NoOfWorkers = &cpu
	}

	return cfg
}
