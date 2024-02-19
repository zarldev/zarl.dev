package repo

import (
	"log/slog"
	"os"
)

var _ AdminsRepository = (*Logger)(nil)

type Logger struct {
	*slog.Logger
	delegate AdminsRepository
}

func WithLogging(d AdminsRepository) *Logger {
	return &Logger{
		Logger:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		delegate: d,
	}
}

// Delete implements AdminsRepository.
func (l *Logger) Delete(username string) error {
	l.Info("deleting admin",
		"username", username)

	err := l.delegate.Delete(username)
	if err != nil {
		l.Error("deleting admin",
			"username", username,
			"error", err)
	}
	return err
}

// Get implements AdminsRepository.
func (l *Logger) Get(username string) (string, error) {
	l.Info("getting admin",
		"username", username)
	p, err := l.delegate.Get(username)
	if err != nil {
		l.Error("getting admin",
			"username", username,
			"error", err)
	}
	return p, err
}

// ListAdmins implements AdminsRepository.
func (l *Logger) ListAdmins() ([]string, error) {
	l.Info("listing admins")
	list, err := l.delegate.ListAdmins()
	if err != nil {
		l.Error("listing admins",
			"error", err)
	}
	return list, err
}

// Set implements AdminsRepository.
func (l *Logger) Set(username string, password string) error {
	l.Info("setting admin",
		"username", username)
	err := l.delegate.Set(username, password)
	if err != nil {
		l.Error("setting admin",
			"username", username,
			"error", err)
	}
	return err
}
