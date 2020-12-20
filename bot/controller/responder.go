package controller

import (
	"io/ioutil"
	"net/http"
	"tbot/model"

	tgram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ImageMessage struct {
	ID     int
	FileID string
}

type Responder struct {
	bot *tgram.BotAPI
	//TODO: extract to separate cache object
	buttonToImage map[int]ImageMessage
}

func NewResponder(bot *tgram.BotAPI) Responder {
	return Responder{bot, make(map[int]ImageMessage)}
}

func (r *Responder) SendMessage(m *tgram.MessageConfig) error {
	resultMessage, err := r.bot.Send(m)
	_, _ = err, resultMessage
	if err != nil {
		return err
	}

	return nil
}

func (r *Responder) SendStickerPreview(s *model.Sticker) error {
	btnMessage, err := r.bot.Send(s.Msg)
	if err != nil {
		return err
	}

	photoMsg, err := r.bot.Send(s.Photo)
	if err != nil {
		return err
	}

	//TODO: resolve idx
	r.buttonToImage[btnMessage.MessageID] = ImageMessage{photoMsg.MessageID, photoMsg.Photo[0].FileID}
	return nil
}

func (r *Responder) DeleteMessage(chatID int64, msgID int) (tgram.APIResponse, error) {
	m := tgram.NewDeleteMessage(chatID, msgID)
	return r.bot.Request(m)

}

func (r *Responder) DeleteStickerPreview(chatID int64, msgID int) error {
	responce, err := r.DeleteMessage(chatID, msgID)
	_, _ = responce, err

	imageMsg := r.buttonToImage[msgID]
	responce, err = r.DeleteMessage(chatID, imageMsg.ID)

	return nil
}

func (r *Responder) NotifyDelete(callbackID string) error {
	responce, err := r.bot.Send(tgram.NewCallback(callbackID, "deleted"))

	_, _ = responce, err

	return nil
}

func (r *Responder) HandleCallback(callback *tgram.CallbackQuery) error {
	switch callback.Data {
	case deleteData:
		callbackMsg := callback.Message
		err := r.DeleteStickerPreview(callbackMsg.Chat.ID, callbackMsg.MessageID)
		_ = err
		//TODO: copy-paste 1
		responce, err := r.bot.Send(tgram.NewCallback(callback.ID, "Deleted"))
		_, _ = responce, err
	case yesData:
		for _, imageMsg := range r.buttonToImage {
			println(imageMsg.FileID)
			url, _ := r.bot.GetFileDirectURL(imageMsg.FileID)
			println(url)
		}
		//TODO: continue, create stickerpack
	case addMoreData:
		responce, err := r.DeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		_, _ = responce, err
		msg, err := r.bot.Send(tgram.NewCallback(callback.ID, "Add more pictures"))
		_, _ = msg, err

	case cancelData:
		responce, err := r.DeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		_, _ = responce, err

		for btnID := range r.buttonToImage {
			err := r.DeleteStickerPreview(callback.Message.Chat.ID, btnID)
			_ = err
			_, _ = responce, err
		}

		//TODO: copy-paste 1
		msg, err := r.bot.Send(tgram.NewCallback(callback.ID, "Deleted"))
		_, _ = msg, err
	}

	return nil
}

func (r *Responder) GetImage(fileID string) ([]byte, error) {
	fileURL, _ := r.bot.GetFileDirectURL(fileID)
	response, err := http.Get(fileURL)
	_ = err
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
