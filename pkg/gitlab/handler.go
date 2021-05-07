package gitlab

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	eRR "github.com/iyurev/notificator/pkg/errors"
	"github.com/iyurev/notificator/pkg/types"
	"github.com/xanzy/go-gitlab"
	"html/template"
	"log"
)

const (
	eventTypeHeader        = "X-Gitlab-Event"
	administratorUserName  = "Administrator"
	administratorUserEmail = "Administrator@local"
)

//go:embed push_event_tg.tmpl
var pushEventTmpl string

type PushEvent struct {
	Event gitlab.PushEvent
	types.ReceiverRef
}

func WebHookHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.GetHeader(eventTypeHeader) != "" {
			//Read request body
			var body bytes.Buffer
			_, err := body.ReadFrom(context.Request.Body)
			if err != nil {
				log.Println(err)
			}
			switch gitlab.HookEventType(context.Request) {
			case gitlab.EventTypePush:
				pushEvent, err := NewPushEvent(body.Bytes())
				if err != nil {
					log.Println(err)
				}
				msg, err := pushEvent.Msg(types.TelegramReceiverType())
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("MSG: %s\n", msg)
			default:
				log.Println("Unknown webhook type")
			}

		}
	}
}

func NewPushEvent(raw []byte) (*PushEvent, error) {
	pushEvent := &PushEvent{}
	if err := json.Unmarshal(raw, &pushEvent.Event); err != nil {
		return nil, err
	}
	return pushEvent, nil
}

func (pe *PushEvent) TgMsg() ([]byte, error) {
	var msg bytes.Buffer
	tmpl, err := template.New("msg").Parse(pushEventTmpl)
	if err != nil {
		return nil, err
	}
	if pe.Event.UserName == administratorUserName {
		pe.Event.UserEmail = administratorUserEmail
	}
	if err := tmpl.Execute(&msg, pe.Event); err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}

func (pe *PushEvent) Msg(rt types.ReceiverType) ([]byte, error) {
	switch rt {
	case types.ReceiverTypeTelegram:
		msg, err := pe.TgMsg()
		if err != nil {
			return nil, err
		}
		return msg, nil
	default:
		return nil, eRR.UnknownReceiverType
	}
}
