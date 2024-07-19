package utils

import (
	"time"
)

// DoWithTries will attempt to run the provided function until it returns nil or the number of attempts is exhausted.
func DoWithTries(fn func() error, attempts int, timeout time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(timeout)
			attempts--
			continue
		}
		return nil
	}
	return
}

//TODO: Есть альтернативный вариант реализации повторяющегося подключения к базе данных, через Dockerize
