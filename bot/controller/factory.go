package controller

import (
	"tbot/model"

	tgram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	deleteData  = "delete"
	yesData     = "Yes"
	addMoreData = "addMore"
	cancelData  = "cancelData"
)

func CreateResultImage(chatID int64, file tgram.FileReader) *model.Sticker {
	dummy := "."
	delButton := tgram.NewMessage(chatID, dummy)
	delButton.ReplyMarkup = tgram.NewInlineKeyboardMarkup([]tgram.InlineKeyboardButton{{Text: "Delete", CallbackData: &deleteData}})
	image := tgram.NewPhotoUpload(chatID, file)

	return &model.Sticker{Msg: delButton, Photo: image}
}

func CreateConfirmKeyboard(chatID int64) *tgram.MessageConfig {
	m := tgram.NewMessage(chatID, "Create sticker pack?")
	m.ReplyMarkup = tgram.NewInlineKeyboardMarkup([]tgram.InlineKeyboardButton{
		{Text: "Yes", CallbackData: &yesData},
		{Text: "Add more...", CallbackData: &addMoreData},
		{Text: "Cancel", CallbackData: &cancelData}})
	return &m
}
