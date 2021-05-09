package sender

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iyurev/notificator/pkg/types"
)

type Sender struct {
	tg *tgbotapi.BotAPI
}

func (s *Sender) Send(event types.Event) error {
	recipient, err := globalConfig.GetTgProjectRecipient(event.Recipient().GetProjectName())
	if err != nil {
		return err
	}
	msg, err := event.Msg(types.TelegramReceiverType())
	if err != nil {
		return err
	}
	tgMsg := tgbotapi.NewMessage(recipient.ChatID, BytesToStr(msg))
	tgMsg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := s.tg.Send(tgMsg); err != nil {
		return err
	}
	return nil
}

func NewTgBot() (*tgbotapi.BotAPI, error) {
	tgbot, err := tgbotapi.NewBotAPI(tgConfig.ApiToken())
	if err != nil {
		return nil, err
	}
	return tgbot, nil
}

func NewSender() (*Sender, error) {
	sender := &Sender{}
	tgbot, err := NewTgBot()
	if err != nil {
		return nil, err
	}
	sender.tg = tgbot
	return sender, nil
}
