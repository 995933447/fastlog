package fmts

import "github.com/995933447/fastlog/logger"

type Format int

const (
	FormatText Format = iota
)

var levelToStdoutColorMap = map[logger.Level]logger.Color{
	logger.LevelDebug:     logger.ColorLightGreen,
	logger.LevelInfo:      logger.ColorLightGreen,
	logger.LevelImportant: logger.ColorBlue,
	logger.LevelWarn:      logger.ColorGreen,
	logger.LevelError:     logger.ColorRed,
	logger.LevelPanic:     logger.ColorRed,
	logger.LevelFatal:     logger.ColorPurple,
}
