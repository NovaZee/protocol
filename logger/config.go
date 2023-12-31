package logger

import "sync"

type Config struct {
	JSON  bool   `yaml:"json,omitempty"`
	Level string `yaml:"level,omitempty"`
	// true to enable log sampling, where the same log message and level will be throttled.
	// we have two layers of sampling
	// 1. global sampling - within a second, it will log the first SampleInitial, then every SampleInterval messages.
	// 2. per participant/track sampling - to be used with Logger.WithItemSampler(). This would be used to throttle
	//    the logs for a particular participant/track.
	Sample bool `yaml:"sample,omitempty"`

	ComponentLevels map[string]string `yaml:"component_levels,omitempty"`

	// global sampling per server
	// when sampling, the first N logs will be logged
	SampleInitial int `yaml:"sample_initial,omitempty"`
	// when sampling, every Mth log will be logged
	SampleInterval int `yaml:"sample_interval,omitempty"`

	// participant/track level sampling
	ItemSampleSeconds  int `yaml:"item_sample_seconds,omitempty"`
	ItemSampleInitial  int `yaml:"item_sample_initial,omitempty"`
	ItemSampleInterval int `yaml:"item_sample_interval,omitempty"`

	lock               sync.Mutex       `yaml:"-"`
	onUpdatedCallbacks []ConfigObserver `yaml:"-"`
}

type ConfigObserver func(*Config) error

func (c *Config) Update(o *Config) error {
	c.lock.Lock()
	c.JSON = o.JSON
	c.Level = o.Level
	c.Sample = o.Sample
	c.SampleInitial = o.SampleInitial
	c.SampleInterval = o.SampleInterval
	c.ItemSampleSeconds = o.ItemSampleSeconds
	c.ItemSampleInitial = o.ItemSampleInitial
	c.ItemSampleInterval = o.ItemSampleInterval
	c.ComponentLevels = o.ComponentLevels
	callbacks := c.onUpdatedCallbacks
	c.lock.Unlock()

	for _, cb := range callbacks {
		if err := cb(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) AddUpdateObserver(cb ConfigObserver) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.onUpdatedCallbacks = append(c.onUpdatedCallbacks, cb)
}
