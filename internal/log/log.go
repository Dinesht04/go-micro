package log

import (
	"io"
	"log/slog"
	"os"
)

func CreateLogger() (*slog.Logger, *os.File, error) {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	handler := slog.NewJSONHandler(multiWriter, nil)

	logger := slog.New(handler)

	logger.Info("Instantiating Logger")

	return logger, file, nil
}
