package sender

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iyurev/notificator/pkg/types"
)

type Sender struct {
	tg *tgbotapi.BotAPI
}

func (s *Sender) Send(event types.Event) error {
	msg, err := event.Msg(types.TelegramReceiverType())
	if err != nil {
		return err
	}
	tgMsg := tgbotapi.NewMessage(666, BytesToStr(msg))
	tgMsg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := s.tg.Send(tgMsg); err != nil {
		return err
	}
	return nil
}

//func NewTgBot(token string) (*TgBot, error) {
//	tgbot, err := tgbotapi.NewBotAPI(token)
//	if err != nil {
//		return &TgBot{}, err
//	}
//	return &TgBot{BotAPI: tgbot}, nil
//}
