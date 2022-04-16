package main

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func generateRegisteredMenu() *tele.ReplyMarkup {
	m := bot.NewMarkup()
	m.ResizeKeyboard = true

	m.Reply(
		m.Row(btnUnlock, btnHelp),
		m.Row(btnForgotKey),
		m.Row(btnCleanMyData),
	)

	return m
}

func generateNonRegisteredMenu() *tele.ReplyMarkup {
	m := bot.NewMarkup()
	m.ResizeKeyboard = true

	m.Reply(
		m.Row(btnRegisterKey),
		m.Row(btnHelp),
	)

	return m
}

func generateStartMenu(c tele.Context) *tele.ReplyMarkup {
	registered := storage.IsUserExists(c.Sender().ID)
	unlocked := locker.IsUnlocked(c.Sender().ID)

	if registered && unlocked {
		return generateMainMenu(c)
	}

	if registered {
		return generateRegisteredMenu()
	}

	return generateNonRegisteredMenu()
}

func generateMainMenu(c tele.Context) *tele.ReplyMarkup {
	m := bot.NewMarkup()
	m.ResizeKeyboard = true

	m.Reply(
		m.Row(btnLock, btnHelp),
		m.Row(btnPasswords, btnAddPassword),
		m.Row(btnPasswordsRemove, btnCleanMyPasswords),
	)

	return m
}

func showPasswords(c tele.Context, passwords []Password) *tele.ReplyMarkup {
	m := bot.NewMarkup()
	m.ResizeKeyboard = true

	m.InlineKeyboard = append(m.InlineKeyboard, []tele.InlineButton{*m.Data("Cancel ❌", btnCancel.Unique).Inline()})
	for _, p := range passwords {
		id := strconv.FormatUint(uint64(p.ID), 10)
		pBtn := m.Data(p.Name, btnPassword.Unique, p.Name, id).Inline()
		m.InlineKeyboard = append(m.InlineKeyboard, []tele.InlineButton{*pBtn})
	}

	return m
}

func showPasswordsRemove(c tele.Context, passwords []Password) *tele.ReplyMarkup {
	m := bot.NewMarkup()
	m.ResizeKeyboard = true

	m.InlineKeyboard = append(m.InlineKeyboard, []tele.InlineButton{*m.Data("Cancel ❌", btnCancel.Unique).Inline()})
	for _, p := range passwords {
		id := strconv.FormatUint(uint64(p.ID), 10)
		pBtn := m.Data(p.Name, btnPasswordRemove.Unique, p.Name, id).Inline()
		m.InlineKeyboard = append(m.InlineKeyboard, []tele.InlineButton{*pBtn})
	}

	return m
}

func showClean(c tele.Context) *tele.ReplyMarkup {
	m := bot.NewMarkup()
	m.ResizeKeyboard = true

	yesBtn := m.Data("Yes", btnYesClean.Unique).Inline()
	noBtn := m.Data("No", btnNoClean.Unique).Inline()
	m.InlineKeyboard = append(m.InlineKeyboard, []tele.InlineButton{*yesBtn, *noBtn})

	return m
}

func showCleanData(c tele.Context) *tele.ReplyMarkup {
	m := bot.NewMarkup()
	m.ResizeKeyboard = true

	yesBtn := m.Data("Yes", btnYesCleanData.Unique).Inline()
	noBtn := m.Data("No", btnNoCleanData.Unique).Inline()
	m.InlineKeyboard = append(m.InlineKeyboard, []tele.InlineButton{*yesBtn, *noBtn})

	return m
}
