// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2026 Pavel Tsayukov p.tsayukov@gmail.com

package cmdlog

import (
	"flag"
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options is options to create a new [zap.Logger].
type Options struct {
	level zap.AtomicLevel

	// These Extra options will be passed into [zap.Config.Build] when a logger
	// creation function is called right before the new [zap.Logger] is returned.
	Extra []zap.Option

	// Enable verbose output. The logging level will be set to [zap.DebugLevel].
	EnableVerbose bool

	// Enable colorful output.
	EnableColor bool
}

// NewOptions creates a new [Options] with the given logging level and optional
// [zap.Option]'s.
func NewOptions(level zap.AtomicLevel, opts ...zap.Option) Options {
	return Options{
		level: level,
		Extra: opts,
	}
}

// SetLogLevel sets the logging level to the given value.
func (opts *Options) SetLogLevel(value string) error {
	return atomicLevel(opts.level).Set(value)
}

// LogLevelFlag defines a new "loglevel" flag for the given [flag.FlagSet].
// The flag will set the logging level.
func (opts *Options) LogLevelFlag(set *flag.FlagSet) {
	const (
		name  = "loglevel"
		usage = "Logging level"
	)
	set.Var((*atomicLevel)(&opts.level), name, usage)
}

// VerboseFlag defines a new "verbose" and "v" flags for the given [flag.FlagSet].
// The flags will set the [Options.EnableVerbose] field.
func (opts *Options) VerboseFlag(set *flag.FlagSet) {
	const (
		name  = "verbose"
		value = false
		usage = "Enable verbose logging"
	)
	set.BoolVar(&opts.EnableVerbose, name, value, usage)
	set.BoolVar(&opts.EnableVerbose, string(name[0]), value, usage)
}

// ColorFlag defines a new "color" flag for the given [flag.FlagSet].
// The flag will set the [Options.EnableColor] field.
func (opts *Options) ColorFlag(set *flag.FlagSet) {
	const (
		name  = "color"
		usage = "Enable colorful logging"
	)
	set.BoolVar(&opts.EnableColor, name, defaultColorFlag(set.Output()), usage)
}

func defaultColorFlag(output io.Writer) bool {
	// See also: https://no-color.org.

	// The output character device must be a file.
	file, ok := output.(*os.File)
	if !ok {
		return false
	}

	// Command-line software which adds ANSI color to its output by default should
	// check for a NO_COLOR environment variable that, when present and not an empty
	// string (regardless of its value), prevents the addition of ANSI color.
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	fd := file.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

// atomicLevel is a wrapper over [zap.AtomicLevel] to implement the [flag.Value]
// interface.
type atomicLevel zap.AtomicLevel

var _ flag.Value = (*atomicLevel)(nil)

func (a atomicLevel) String() string {
	return zap.AtomicLevel(a).String()
}

func (a atomicLevel) Set(value string) error {
	level, err := zapcore.ParseLevel(value)
	if err != nil {
		return err
	}
	zap.AtomicLevel(a).SetLevel(level)
	return nil
}
