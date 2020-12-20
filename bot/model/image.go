package model

import tgram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ImageSet struct {
	Images    []tgram.PhotoSize
	MessageID int
}
