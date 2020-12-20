package model

import (
	tgram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sticker struct {
	Msg   tgram.MessageConfig
	Photo tgram.PhotoConfig
}
