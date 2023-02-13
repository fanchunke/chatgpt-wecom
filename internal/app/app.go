package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/fanchunke/chatgpt-wecom/conf"
	"github.com/fanchunke/chatgpt-wecom/internal/api"
	"github.com/fanchunke/chatgpt-wecom/pkg/httpserver"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom"
	"github.com/fanchunke/xgpt3"
	"github.com/fanchunke/xgpt3/conversation/ent"
	"github.com/fanchunke/xgpt3/conversation/ent/chatent"

	"github.com/rs/zerolog/log"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func Run(cfg *config.Config) {
	log.Info().Msgf("Config: %v", cfg)

	// 初始化 gpt client
	gptClient := gogpt.NewClient(cfg.GPT.ApiKey)

	// 初始化数据库 client
	dbConf := cfg.Database
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DBName)
	chatentClient, err := chatent.Open(dbConf.Dialect, s)
	if err != nil {
		log.Fatal().Err(err).Msg("ent - open database failed")
	}

	// 初始化 xgpt3 client
	xgpt3Client := xgpt3.NewClient(gptClient, ent.New(chatentClient))

	// 初始化 wecom app
	wecomApp := wecom.NewApp(
		cfg.WeCom.Url,
		cfg.WeCom.CorpId,
		cfg.WeCom.AgentId,
		cfg.WeCom.CorpSecret,
		cfg.WeCom.EncodingAESKey,
		cfg.WeCom.Token,
		time.Duration(5)*time.Second,
	)

	handler, err := api.NewRouter(cfg, xgpt3Client, wecomApp)
	if err != nil {
		log.Fatal().Err(err).Msg("api - Router - api.Router failed")
	}
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	httpServer.Start()
	log.Info().Msg("Server Started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Msgf("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		log.Error().Err(err).Msg("app - Run - httpServer.Notify")
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("app - Run - httpServer.Shutdown")
	}

}
