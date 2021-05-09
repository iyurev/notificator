package server

import (
	"github.com/gin-gonic/gin"
	"github.com/iyurev/notificator/pkg/gitlab"
	"github.com/iyurev/notificator/pkg/sender"
	"github.com/iyurev/notificator/pkg/types"
	"log"
)

type controller struct {
	S              types.Sender
	gitLabHooksSvc *gitlab.WHSvc
}

func newController() (*controller, error) {
	con := &controller{}
	s, err := sender.NewSender()
	if err != nil {
		return nil, err
	}
	con.S = s
	con.gitLabHooksSvc = gitlab.NewWHSvc(s)
	return con, nil
}

func StartServer() {
	con, err := newController()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/gitlab", con.gitLabHooksSvc.HookHandler())
	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}

}
