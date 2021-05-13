package server

import (
	"github.com/gin-gonic/gin"
	"github.com/iyurev/notificator/pkg/handlers"
	"github.com/iyurev/notificator/pkg/logger"
	"github.com/iyurev/notificator/pkg/sender"
	"github.com/iyurev/notificator/pkg/types"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	startLogger zap.Logger
)

func init() {
	var err error
	startLogger, err = logger.New()
	if err != nil {
		log.Fatal(err)
	}
}

type controller struct {
	S              types.Sender
	gitLabHooksSvc *handlers.GitLabSvc
	harborSvc      *handlers.HarborSvc
	log            zap.Logger
}

func newController() (*controller, error) {
	conLogger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	con := &controller{}
	s, err := sender.NewSender()
	if err != nil {
		return nil, err
	}
	con.S = s
	con.log = conLogger
	gitlabHooksSvc, err := handlers.NewGitLabSvc(s)
	if err != nil {
		return nil, err
	}
	harborSvc, err := handlers.NewHarborSvc(s)
	if err != nil {
		return nil, err
	}
	con.gitLabHooksSvc = gitlabHooksSvc
	con.harborSvc = harborSvc
	return con, nil
}

func StartServer() {
	con, err := newController()
	if err != nil {
		startLogger.Fatal(err.Error())
		os.Exit(1)
	}
	r := gin.Default()
	r.POST("/gitlab", con.gitLabHooksSvc.HookHandler())
	r.POST("/harbor", con.harborSvc.HookHandler())
	if err := r.Run(); err != nil {
		startLogger.Fatal(err.Error())
	}

}
