// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2026 Pavel Tsayukov p.tsayukov@gmail.com

package cmdlog

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewDevelopmentLogger creates a new [zap.Logger] for logging development
// and deployment tools with optional [Options].
func NewDevelopmentLogger(options ...Options) (*zap.Logger, error) {
	logger, err := newLogger(options, func(options Options) *zap.Config {
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = options.level
		if options.EnableVerbose {
			cfg.Level.SetLevel(zapcore.DebugLevel)
		} else {
			cfg.DisableCaller = true
			cfg.DisableStacktrace = true
		}

		if options.EnableColor {
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		return &cfg
	})
	if err != nil {
		return nil, fmt.Errorf("logger.NewDevelopmentLogger: %w", err)
	}
	return logger, nil
}

// NewProductionLogger creates a new [zap.Logger] for logging production with
// optional [Options].
func NewProductionLogger(options ...Options) (*zap.Logger, error) {
	logger, err := newLogger(options, func(options Options) *zap.Config {
		cfg := zap.NewProductionConfig()

		cfg.Level = options.level
		if options.EnableVerbose {
			cfg.Level.SetLevel(zap.DebugLevel)
		}

		if options.EnableColor {
			cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		}

		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		return &cfg
	})
	if err != nil {
		return nil, fmt.Errorf("logger.NewProductionLogger: %w", err)
	}
	return logger, nil
}

func newLogger(
	zeroOrOneOptions []Options,
	cfgMaker func(o Options) *zap.Config,
) (*zap.Logger, error) {
	var options Options
	switch len(zeroOrOneOptions) {
	default:
		return nil, errors.New("zero or one Options is allowed")
	case 0:
	case 1:
		options = zeroOrOneOptions[0]
	}

	logger, err := cfgMaker(options).Build(options.Extra...)
	if err != nil {
		return nil, fmt.Errorf("logger creation failed: %w", err)
	}
	return logger, nil
}
