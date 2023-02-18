package app

import (
	"context"

	config "github.com/fanchunke/chatgpt-wecom/conf"
	"github.com/fanchunke/xgpt3/conversation/ent/chatent"
)

func Migrate(cfg *config.Config) error {
	dbConf := cfg.Database
	client, err := chatent.Open(dbConf.Driver, dbConf.DataSource)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		return err
	}
	return nil
}
