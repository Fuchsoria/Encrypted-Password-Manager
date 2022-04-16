package main

import (
	"errors"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"
)

// REGISTRATION REQUIRED
func btnPasswordsHandler(c tele.Context) error {
	passwords, err := storage.GetPasswords(c.Sender().ID)
	if err != nil {
		logger.Error(err)

		return err
	}

	if len(passwords) == 0 {
		return c.Send("*You don't have any passwords yet*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	return c.Send("*Your Passwords:*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		showPasswords(c, passwords))
}

// REGISTRATION REQUIRED
func btnPasswordsRemoveHandler(c tele.Context) error {
	passwords, err := storage.GetPasswords(c.Sender().ID)
	if err != nil {
		logger.Error(err)

		return err
	}

	if len(passwords) == 0 {
		return c.Send("*You don't have any passwords yet*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	return c.Send("*Your Passwords for removing:*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		showPasswordsRemove(c, passwords))
}

// REGISTRATION REQUIRED
func btnLockHandler(c tele.Context) error {
	locker.Lock(c.Sender().ID)

	return c.Send("*Successfully locked*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		generateRegisteredMenu(),
	)
}

// REGISTRATION REQUIRED
func btnUnlockHandler(c tele.Context) error {
	return c.Send("*Unlock usage:*\n`/unlock <key>`",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
	)
}

// REGISTRATION REQUIRED
func btnPasswordHandler(c tele.Context) error {
	defer c.DeleteAfter(time.Second * 10)

	if len(c.Args()) != 2 {
		return errors.New("Wrong arguments in password")
	}

	name := c.Args()[0]
	passID := c.Args()[1]

	pass, err := storage.GetPassword(c.Sender().ID, passID)
	if err != nil {
		logger.Error(err)

		return c.Edit(
			"*Something went wrong with your password*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		)
	}

	return c.Edit(
		fmt.Sprintf("*Your password for %s is:*\n`%s`\n__Will be available for 10 seconds__", name, pass),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
	)
}

// REGISTRATION REQUIRED
func btnPasswordRemoveHandler(c tele.Context) error {
	defer c.DeleteAfter(time.Second * 10)

	if len(c.Args()) != 2 {
		return errors.New("Wrong arguments in password")
	}

	name := c.Args()[0]
	passID := c.Args()[1]

	err := storage.RemovePassword(c.Sender().ID, passID)
	if err != nil {
		logger.Error(err)

		return c.Edit(
			"*Something went wrong with your password removing*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	return c.Edit(
		fmt.Sprintf("*Your %s password successfully removed*", name),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
}

// REGISTRATION REQUIRED
func btnAddPasswordHandler(c tele.Context) error {
	return c.Send("*Add password usage:*\n`/add <password> <long name>`",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
}

// ONLY WITHOUT REGISTRATION
func btnRegisterKeyHandler(c tele.Context) error {
	return c.Send("*Register key usage:*\n`/register <key> <repeat key> <key description>`",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
}

// REGISTRATION REQUIRED
func btnCleanMyPasswordsHandler(c tele.Context) error {
	return c.Send("*Are you sure you want to clean your passwords?*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2}, showClean(c))
}

// REGISTRATION REQUIRED
func btnYesCleanHandler(c tele.Context) error {
	defer c.Delete()

	err := storage.CleanMyPasswords(c.Sender().ID)
	if err != nil {
		logger.Error(err)

		return c.Send("*Something went wrong on your data removing`",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	return c.Send("*All your passwords were removed*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
}

// REGISTRATION REQUIRED
func btnNoCleanHandler(c tele.Context) error {
	defer c.Delete()

	return nil
}

// REGISTRATION REQUIRED
func btnForgotKeyHandler(c tele.Context) error {
	user, err := storage.GetUser(c.Sender().ID)
	if err != nil {
		logger.Error(err)

		return c.Send("*User not found, please register before*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
			generateNonRegisteredMenu(),
		)
	}

	return c.Send(fmt.Sprintf("*Your key description is:*\n`%s`", user.HashKeyDescription),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
	)
}

// REGISTRATION REQUIRED
func btnCleanMyDataHandler(c tele.Context) error {
	return c.Send("*Are you sure you want to Clean your data?*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		showCleanData(c))
}

// REGISTRATION REQUIRED
func btnYesCleanDataHandler(c tele.Context) error {
	defer c.Delete()
	defer locker.Lock(c.Sender().ID)

	err := storage.CleanMyData(c.Sender().ID)
	if err != nil {
		logger.Error(err)

		return c.Send("*Something went wrong on your data removing`",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
			generateStartMenu(c),
		)
	}

	return c.Send("*All your data were removed*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		generateNonRegisteredMenu(),
	)
}

// REGISTRATION REQUIRED
func btnNoCleanDataHandler(c tele.Context) error {
	defer c.Delete()

	return nil
}

// REGISTRATION REQUIRED
func btnCancelHandler(c tele.Context) error {
	defer c.Delete()

	return nil
}
