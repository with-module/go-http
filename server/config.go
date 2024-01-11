package server

import "time"

type (
	Config struct {
		Addr    string        `json:"addr" yaml:"addr"`
		Timeout TimeoutConfig `json:"timeout" yaml:"timeout"`

		moreOptions struct {
			RuntimeErrorHandle func(error)
		}
	}

	TimeoutConfig struct {
		Shutdown time.Duration `json:"shutdown" yaml:"shutdown"`
		Read     time.Duration `json:"read" yaml:"read"`
		Write    time.Duration `json:"write" yaml:"write"`
		Idle     time.Duration `json:"idle" yaml:"idle"`
	}

	Option interface {
		use(Config) Config
	}

	useOption func(Config) Config
)

func (fn useOption) use(config Config) Config {
	return fn(config)
}

func UseRuntimeErrHandle(fn func(error)) Option {
	return useOption(func(config Config) Config {
		config.moreOptions.RuntimeErrorHandle = fn
		return config
	})
}
