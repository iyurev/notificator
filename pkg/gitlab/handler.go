package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xanzy/go-gitlab"
	"html/template"
	"log"
)

func WebHookHandler() gin.HandlerFunc {
	return func(req *gin.Context) {
		if req.GetHeader("X-Gitlab-Event") == "Push Hook" {
			var body bytes.Buffer
			_, err := body.ReadFrom(req.Request.Body)
			if err != nil {
				log.Println(err)
			}
			pushEvent := PushEvent{}
			if err := pushEvent.Unmarshal(body.Bytes()); err != nil {
				log.Println(err)
			}
			msg, err := pushEvent.Msg()
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("MSG: %s\n", msg)
		}
	}
}

type PushEvent struct {
	Event gitlab.PushEvent
}

func (pe *PushEvent) Unmarshal(raw []byte) error {
	if err := json.Unmarshal(raw, &pe.Event); err != nil {
		return err
	}
	return nil
}

func (pe *PushEvent) Msg() ([]byte, error) {
	var msg bytes.Buffer
	t := "GitLab project: {{ .Project.Name }}\n Description: {{ .Repository.Description }}\n  URL: {{ .Repository.GitHTTPURL }}\n REF: {{.Ref}}\n Before: {{ .Before }}\n After: {{ .After }}\n Commiter: {{ .UserName }}\n"
	tmpl, err := template.New("msg").Parse(t)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(&msg, pe.Event); err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}
