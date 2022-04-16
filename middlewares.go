package main

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

func lockedMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !locker.IsUnlocked(c.Sender().ID) {
			return c.Send("*Please unlock first\\.*",
				&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
				generateStartMenu(c))
		}

		user, err := storage.GetUser(c.Sender().ID)
		if err != nil {
			logger.Error(err)

			return c.Send("*User not found, please register before*",
				&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
		}

		if user.UnlockedAt.Before(time.Now().Add(-time.Hour * 1)) {
			locker.Lock(c.Sender().ID)

			return c.Send("*Session expired\\. Please unlock first\\.*",
				&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
				generateStartMenu(c))
		}

		return next(c)
	}
}

func lockedRemoveMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !locker.IsUnlocked(c.Sender().ID) {
			defer c.Delete()

			return c.Send("*Please unlock first\\.*",
				&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
				generateStartMenu(c))
		}

		user, err := storage.GetUser(c.Sender().ID)
		if err != nil {
			logger.Error(err)

			return c.Send("*User not found, please register before*",
				&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
		}

		if user.UnlockedAt.Before(time.Now().Add(-time.Hour * 1)) {
			locker.Lock(c.Sender().ID)

			return c.Send("*Session expired\\. Please unlock first\\.*",
				&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
				generateStartMenu(c))
		}

		return next(c)
	}
}

func registerRequiredMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !storage.IsUserExists(c.Sender().ID) {
			defer c.Delete()

			return c.Send(generateNonRegisteredMenu())
		}

		return next(c)
	}
}

func nonRegisterRequiredMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if storage.IsUserExists(c.Sender().ID) {
			defer c.Delete()

			return c.Send(generateRegisteredMenu())
		}

		return next(c)
	}
}
