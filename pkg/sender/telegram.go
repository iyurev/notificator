package sender

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iyurev/notificator/pkg/errors"
	"github.com/iyurev/notificator/pkg/logger"
	"github.com/iyurev/notificator/pkg/types"
	"go.uber.org/zap"
)

type Sender struct {
	tg  *tgbotapi.BotAPI
	log zap.Logger
}

func (s *Sender) Send(event types.Event) error {
	tgRecipient, err := globalConfig.GetTgProjectRecipient(event.Recipient().GetProjectName())
	if err != nil {
		if err == errors.NoSuchRecipient {
			s.log.Info(err.Error(), logger.ProjectName(event.Recipient().GetProjectName()))
		} else {
			return err
		}
	}
	if err == nil {
		msg, err := event.Msg(types.TelegramReceiverType())
		if err != nil {
			return err
		}
		tgMsg := tgbotapi.NewMessage(tgRecipient.ChatID, BytesToStr(msg))
		tgMsg.ParseMode = tgbotapi.ModeHTML
		if _, err := s.tg.Send(tgMsg); err != nil {
			return err
		}
		s.log.Info("message was sent to a telegram group", logger.TelegramChatID(tgRecipient.ChatID), logger.ProjectName(event.Recipient().GetProjectName()))
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
	senderLogger, err := logger.New()
	if err != nil {
		return nil, err
	}
	sender := &Sender{
		log: senderLogger,
	}
	tgbot, err := NewTgBot()
	if err != nil {
		return nil, err
	}
	sender.tg = tgbot
	return sender, nil
}
