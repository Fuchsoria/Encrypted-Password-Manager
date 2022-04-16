package main

import (
	"flag"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	storage = NewStorage()
	locker  = NewLocker()
	logger  = NewLogger(true)
	bot     *tele.Bot
)

var (
	DEBUG bool
	TOKEN string
)

func init() {
	flag.BoolVar(&DEBUG, "debug", false, "debug mode")
	flag.StringVar(&TOKEN, "token", "", "bot token")
	flag.Parse()

	err := storage.Connect()
	if err != nil {
		logger.Error(err)

		logger.Fatal(fmt.Errorf("cannot connect to storage: %w", err))
	}

	pref := tele.Settings{
		Token:  TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		logger.Error(err)
	}

	bot = b
}

func main() {
	defer storage.Disconnect()

	// buttons
	bot.Handle(&btnPasswords, btnPasswordsHandler, registerRequiredMiddleware, lockedMiddleware)
	bot.Handle(&btnPasswordsRemove, btnPasswordsRemoveHandler, registerRequiredMiddleware, lockedMiddleware)
	bot.Handle(&btnLock, btnLockHandler, registerRequiredMiddleware, lockedMiddleware)
	bot.Handle(&btnUnlock, btnUnlockHandler, registerRequiredMiddleware)
	bot.Handle(&btnPassword, btnPasswordHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle(&btnPasswordRemove, btnPasswordRemoveHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle(&btnAddPassword, btnAddPasswordHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle(&btnRegisterKey, btnRegisterKeyHandler, nonRegisterRequiredMiddleware)
	bot.Handle(&btnForgotKey, btnForgotKeyHandler, registerRequiredMiddleware)
	bot.Handle(&btnCleanMyData, btnCleanMyDataHandler, registerRequiredMiddleware)
	bot.Handle(&btnCleanMyPasswords, btnCleanMyPasswordsHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle(&btnYesClean, btnYesCleanHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle(&btnNoClean, btnNoCleanHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle(&btnYesCleanData, btnYesCleanDataHandler, registerRequiredMiddleware)
	bot.Handle(&btnNoCleanData, btnNoCleanDataHandler, registerRequiredMiddleware)
	bot.Handle(&btnCancel, btnCancelHandler, registerRequiredMiddleware)
	bot.Handle(&btnHelp, helpHandler)

	// commands
	bot.Handle("/unlock", unlockHandler, registerRequiredMiddleware)
	bot.Handle("/u", unlockHandler, registerRequiredMiddleware)
	bot.Handle("/start", startHandler)
	bot.Handle("/s", startHandler)
	bot.Handle("/help", helpHandler)
	bot.Handle("/h", helpHandler)
	bot.Handle("/passwords", btnPasswordsHandler, registerRequiredMiddleware, lockedMiddleware)
	bot.Handle("/p", btnPasswordsHandler, registerRequiredMiddleware, lockedMiddleware)
	bot.Handle("/add", addHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle("/a", addHandler, registerRequiredMiddleware, lockedRemoveMiddleware)
	bot.Handle("/register", registerHandler, nonRegisterRequiredMiddleware)
	bot.Handle("/r", registerHandler, nonRegisterRequiredMiddleware)

	bot.Start()
}
