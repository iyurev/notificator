package handlers

//https://github.com/goharbor/harbor/blob/master/src/pkg/notifier/model/event.go

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/goharbor/harbor/src/pkg/notifier/model"
	"github.com/iyurev/notificator/pkg/errors"
	"github.com/iyurev/notificator/pkg/logger"
	"github.com/iyurev/notificator/pkg/types"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
)

type HarborSvc struct {
	S   types.Sender
	log *zap.Logger
}

//go:embed harbor_event_tg.tmpl
var harborEventTmpl string

func (svc *HarborSvc) HookHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		if !ContentTypeIsJSON(context) {
			svc.log.Info("wrong content type")
			context.JSON(http.StatusBadRequest, gin.H{"reason": "wrong content type"})
			return
		}
		//Read request body
		var body bytes.Buffer
		_, err := body.ReadFrom(context.Request.Body)
		if err != nil {
			svc.log.Error(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"reason": "error while reading request body"})
			return
		}
		harborEvent, err := NewHarborEvent(body.Bytes())
		if err != nil {
			svc.log.Error(err.Error())
			context.JSON(http.StatusBadRequest, gin.H{"reason": "failed to make event from request's body"})
			return
		}
		log.Println(context.Request.Header)
		if err := svc.S.Send(harborEvent); err != nil {
			svc.log.Error(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"reason": "failed to send message"})
			return
		}
	}
}

type HarborEvent struct {
	Event model.Payload
}

func (he *HarborEvent) Recipient() *types.RecipientRef {
	var projectName string
	if he.Event.EventData.Repository != nil {
		projectName = he.Event.EventData.Repository.Namespace
	}
	recipientRef := &types.RecipientRef{
		Project: projectName,
	}
	return recipientRef

}

func (he *HarborEvent) TgMsg() ([]byte, error) {
	var msg bytes.Buffer
	tmpl, err := template.New("msg").Parse(harborEventTmpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(&msg, he.Event); err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}

func (he *HarborEvent) Msg(rt types.ReceiverType) ([]byte, error) {
	switch rt {
	case types.ReceiverTypeTelegram:
		msg, err := he.TgMsg()
		if err != nil {
			return nil, err
		}
		return msg, nil
	default:
		return nil, errors.UnknownReceiverType
	}

}

func NewHarborEvent(raw []byte) (*HarborEvent, error) {
	he := &HarborEvent{}
	if err := json.Unmarshal(raw, &he.Event); err != nil {
		return nil, err
	}
	return he, nil
}

func NewHarborSvc(sender types.Sender) (*HarborSvc, error) {
	harborLogger, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &HarborSvc{
		S:   sender,
		log: harborLogger.Named("harborHooksHandler"),
	}, nil
}
