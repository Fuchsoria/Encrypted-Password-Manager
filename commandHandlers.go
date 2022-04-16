package main

import (
	"fmt"
	"strings"

	"github.com/cespare/xxhash"
	tele "gopkg.in/telebot.v3"
)

// REGISTRATION REQUIRED
func unlockHandler(c tele.Context) error {
	defer c.Delete()

	if len(c.Args()) != 1 {
		return c.Send("*Unlock handler usage:*\n`/unlock <key>`",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	key := c.Args()[0]

	user, err := storage.GetUser(c.Sender().ID)
	if err != nil {
		logger.Error(err)

		return c.Send("*User not found, please register before unlocking*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		)
	}

	hashKey := fmt.Sprint(xxhash.Sum64String(key))

	if hashKey != user.HashKey {
		return c.Send("*Wrong key*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	err = locker.Unlock(c.Sender().ID, key)
	if err != nil {
		logger.Error(err)

		return err
	}

	return c.Send("*Successfully unlocked*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		generateMainMenu(c),
	)
}

func startHandler(c tele.Context) error {
	return c.Send(`*Welcome to Encrypted Wallet Bot\!*
	*Here is a Flow:*
	*1\) Register your Secret Key*
	*2\) Unlock your Wallet*
	*3\) Add and Use your Passwords\!*
	*4\) Profit\!*,
	*Use /help* or */h* \- to show help
	`,
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		generateStartMenu(c))
}

func helpHandler(c tele.Context) error {
	return c.Send(`*Encrypted Wallet Bot Help:*
	*Commands:*
	*/start* or */s* \- start bot
	*/help* or */h* \- show help
	*/register* or */r* \- register your secret key
	*/unlock* or */u* \- unlock your wallet \(register required\)
	*/passwords* or */p* \- show all passwords \(register required\)
	*/add* or */a* \- add password \(register required\)`,
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
}

// REGISTRATION REQUIRED
func addHandler(c tele.Context) error {
	defer c.Delete()

	if len(c.Args()) < 2 {
		return c.Send("*Add password usage:*\n`/add <password> <long name>`",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	args := c.Args()
	password := args[0]
	name := strings.Join(args[1:], " ")

	err := storage.SavePassword(c.Sender().ID, name, password)
	if err != nil {
		logger.Error(err)

		return err
	}

	return c.Send(fmt.Sprintf("*Password for* __*%s*__ *was added*", name),
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
}

// ONLY NOT REQUIRED
func registerHandler(c tele.Context) error {
	defer c.Delete()

	if len(c.Args()) < 3 {
		return c.Send("*Register key usage:*\n`/register <key> <repeat key> <key description>`",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	args := c.Args()
	key := args[0]
	repeatKey := args[1]
	keyDescription := strings.Join(args[2:], " ")

	if key != repeatKey {
		return c.Send("*Keys do not match*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	xxh := fmt.Sprint(xxhash.Sum64String(key))

	err := storage.SaveUser(c.Sender().ID, xxh, keyDescription)
	if err != nil {
		logger.Error(err)

		return c.Send("*Cannot create user, please try later*",
			&tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	}

	return c.Send("*User was successfully created*",
		&tele.SendOptions{ParseMode: tele.ModeMarkdownV2},
		generateRegisteredMenu(),
	)
}
