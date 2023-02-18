package api

import (
	"net/http"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom"
	"github.com/fanchunke/xgpt3"
	"github.com/rs/zerolog/log"

	"github.com/fanchunke/chatgpt-wecom/internal/middleware"

	config "github.com/fanchunke/chatgpt-wecom/conf"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type router struct {
	*gin.Engine
	cfg *config.Config
}

func NewRouter(cfg *config.Config, xgpt3Client *xgpt3.Client, wecomApp *wecom.WeComApp) (http.Handler, error) {
	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()
	pprof.Register(e, "debug/pprof")

	r := &router{Engine: e, cfg: cfg}
	r.Use(middleware.Logger())
	r.Use(middleware.URLHandler("url"))
	r.Use(middleware.MethodHandler("method"))
	r.Use(middleware.RequestIDHandler("requestId", "X-Request-Id"))
	r.Use(middleware.AccessHandler())
	r.GET("/healthz", r.Healthz)

	callback := NewCallbackHandler(cfg, xgpt3Client, wecomApp)
	msgHandler, err := wecomApp.RxMessageHandler(callback)
	if err != nil {
		log.Error().Err(err).Msgf("Init RxMessageHandler failed: %s", err)
		return nil, err
	}
	r.GET("/wecom/receive", gin.WrapH(msgHandler))
	r.POST("/wecom/receive", gin.WrapH(msgHandler))
	return r, nil
}
