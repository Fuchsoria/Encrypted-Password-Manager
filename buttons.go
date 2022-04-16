package main

import tele "gopkg.in/telebot.v3"

var (
	_menu = &tele.ReplyMarkup{
		ResizeKeyboard: true,
	}

	btnPasswords        = _menu.Text("🔑 Passwords")
	btnAddPassword      = _menu.Text("➕ Add Password")
	btnLock             = _menu.Text("🔒 Lock")
	btnUnlock           = _menu.Text("🔓 Unlock")
	btnRegisterKey      = _menu.Text("🔑 Register Encryption Key")
	btnForgotKey        = _menu.Text("🔑 Forgot Encryption Key")
	btnCleanMyData      = _menu.Text("🗑 Clean My Full Data")
	btnPasswordsRemove  = _menu.Text("🗑 Remove Password")
	btnCleanMyPasswords = _menu.Text("🗑 Clean My Passwords")
	btnHelp             = _menu.Text("📖 Help")
	btnPassword         = tele.InlineButton{Unique: "password"}
	btnPasswordRemove   = tele.InlineButton{Unique: "password_remove"}
	btnYesClean         = tele.InlineButton{Unique: "yes_clean"}
	btnNoClean          = tele.InlineButton{Unique: "no_clean"}
	btnYesCleanData     = tele.InlineButton{Unique: "yes_clean_data"}
	btnNoCleanData      = tele.InlineButton{Unique: "no_clean_data"}
	btnCancel           = tele.InlineButton{Unique: "cancel"}
)
