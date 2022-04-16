package main

import (
	"errors"
)

type Locker struct {
	data map[int64]string
}

func (l *Locker) IsUnlocked(userID int64) bool {
	_, ok := l.data[userID]

	return ok
}

func (l *Locker) Lock(userID int64) {
	delete(l.data, userID)
}

func (l *Locker) Unlock(userID int64, hashKey string) error {
	l.data[userID] = hashKey

	return storage.UpdateUnlocked(userID)
}

func (l *Locker) GetKey(userID int64) (string, error) {
	if !l.IsUnlocked(userID) {
		return "", errors.New("user is locked")
	}

	return l.data[userID], nil
}

func NewLocker() *Locker {
	return &Locker{
		data: make(map[int64]string),
	}
}
