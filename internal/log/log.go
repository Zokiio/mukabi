// Package log provides logging configuration and setup for the application
package log

import (
	"log/slog"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/topi314/tint"
)

// ANSI color codes for terminal output
const (
	ansiFaint         = "\033[2m"
	ansiWhiteBold     = "\033[37;1m"
	ansiYellowBold    = "\033[33;1m"
	ansiCyanBold      = "\033[36;1m"
	ansiCyanBoldFaint = "\033[36;1;2m"
	ansiRedBold       = "\033[31;1m"

	ansiRed     = "\033[31m"
	ansiYellow  = "\033[33m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
)

// Config defines logging configuration options
type Config struct {
	Level     slog.Level `toml:"level"`      // Minimum log level to output
	Format    string     `toml:"format"`     // Output format (json or text)
	AddSource bool       `toml:"add_source"` // Include source file and line in log output
	NoColor   bool       `toml:"no_color"`   // Disable colored output
}

// Setup configures the global logger based on the provided configuration
func Setup(cfg Config) {
	var handler slog.Handler

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: cfg.AddSource,
			Level:     cfg.Level,
		})

	case "text":
		handler = tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
			AddSource: cfg.AddSource,
			Level:     cfg.Level,
			NoColor:   cfg.NoColor,
			LevelColors: map[slog.Level]string{
				slog.LevelDebug: ansiMagenta,
				slog.LevelInfo:  ansiGreen,
				slog.LevelWarn:  ansiYellow,
				slog.LevelError: ansiRed,
			},
			Colors: map[tint.Kind]string{
				tint.KindTime:            ansiYellowBold,
				tint.KindSourceFile:      ansiCyanBold,
				tint.KindSourceSeparator: ansiCyanBoldFaint,
				tint.KindSourceLine:      ansiCyanBold,
				tint.KindMessage:         ansiWhiteBold,
				tint.KindKey:             ansiFaint,
				tint.KindSeparator:       ansiFaint,
				tint.KindValue:           ansiWhiteBold,
				tint.KindErrorKey:        ansiRedBold,
			},
		})

	default:
		slog.Error("Unsupported log format, defaulting to text",
			slog.String("format", cfg.Format))

		// Set up default text handler without colors
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: cfg.AddSource,
			Level:     cfg.Level,
		})
	}

	slog.SetDefault(slog.New(handler))
}
