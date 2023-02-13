package main

import (
	"context"
	"flag"
	"fmt"

	config "github.com/fanchunke/chatgpt-wecom/conf"
	"github.com/fanchunke/chatgpt-wecom/internal/app"
	"github.com/fanchunke/chatgpt-wecom/pkg/logger"
	"github.com/fanchunke/xgpt3/conversation/ent/chatent"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
	_ "go.uber.org/automaxprocs"
)

func main() {
	conf := flag.String("conf", "conf/online.conf", "配置文件")
	initEnt := flag.Bool("init-ent", false, "是否初始化数据库")
	flag.Parse()
	cfg, err := config.New(*conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Config Load Failed")
	}
	logger.Configure(logger.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Level:                 cfg.Logger.Level,
		Filename:              cfg.Logger.Filename,
	})

	// 数据库迁移
	if *initEnt {
		if err := migrate(cfg); err != nil {
			log.Fatal().Err(err).Msg("failed creating schema resources")
		}
		return
	}

	app.Run(cfg)
}

func migrate(cfg *config.Config) error {
	dbConf := cfg.Database
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DBName)
	client, err := chatent.Open(dbConf.Dialect, s)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
