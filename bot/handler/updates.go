package handler

import (
	"fmt"
	"tbot/controller"
	"tbot/controller/image"
	"tbot/model"

	tgram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update struct {
	responder      controller.Responder
	imageProcessor *image.Processor
}

func NewUpdate(responder controller.Responder, close chan struct{}) *Update {
	update := &Update{responder: responder}
	update.imageProcessor = image.NewProcessor(close, update)
	return update
}

func (u *Update) Handle(update tgram.Update) {

	//TODO: handle documents too
	if update.Message != nil && update.Message.Photo != nil {
		u.imageProcessor.Put(model.ImageSet{Images: update.Message.Photo, MessageID: update.Message.MessageID})
		u.responder.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		//TODO: "waiting" notification message
	}

	if update.CallbackQuery != nil {
		err := u.responder.HandleCallback(update.CallbackQuery)
		_ = err
	}
}

func (u *Update) GetImage(imageSet model.ImageSet) []byte {
	result := make([][]byte, 0, len(imageSet.Images))
	for _, image := range imageSet.Images {
		fmt.Printf("fileID: %s height: %d width: %d", image.FileID, image.Height, image.Width)
		bytes, _ := u.responder.GetImage(image.FileID)
		result = append(result, bytes)
	}
	return result[0]
}
