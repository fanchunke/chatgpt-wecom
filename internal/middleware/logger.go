package middleware

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger injects log into requests context.
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(log.Logger.WithContext(ctx.Request.Context()))
		ctx.Next()
	}
}

// URLHandler adds the requested URL as a field to the context's logger
// using fieldKey as field key.
func URLHandler(fieldKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		log := zerolog.Ctx(r.Context())
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(fieldKey, r.URL.String())
		})
		ctx.Next()
	}
}

// MethodHandler adds the request method as a field to the context's logger
// using fieldKey as field key.
func MethodHandler(fieldKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		log := zerolog.Ctx(r.Context())
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(fieldKey, r.Method)
		})
		ctx.Next()
	}
}

type idKey struct{}

// IDFromRequest returns the unique id associated to the request if any.
func IDFromRequest(r *http.Request) (id xid.ID, ok bool) {
	if r == nil {
		return
	}
	return IDFromCtx(r.Context())
}

// IDFromCtx returns the unique id associated to the context if any.
func IDFromCtx(ctx context.Context) (id xid.ID, ok bool) {
	id, ok = ctx.Value(idKey{}).(xid.ID)
	return
}

// CtxWithID adds the given xid.ID to the context
func CtxWithID(ctx context.Context, id xid.ID) context.Context {
	return context.WithValue(ctx, idKey{}, id)
}

// RequestIDHandler returns a handler setting a unique id to the request which can
// be gathered using IDFromRequest(req). This generated id is added as a field to the
// logger using the passed fieldKey as field name. The id is also added as a response
// header if the headerName is not empty.
//
// The generated id is a URL safe base64 encoded mongo object-id-like unique id.
// Mongo unique id generation algorithm has been selected as a trade-off between
// size and ease of use: UUID is less space efficient and snowflake requires machine
// configuration.
func RequestIDHandler(fieldKey, headerName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		id, ok := IDFromRequest(r)
		if !ok {
			id = xid.New()
			r = r.WithContext(CtxWithID(r.Context(), id))
		}
		if fieldKey != "" {
			log := zerolog.Ctx(r.Context())
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str(fieldKey, id.String())
			})
		}
		if headerName != "" {
			ctx.Header(headerName, id.String())
		}
		ctx.Next()
	}
}

// CustomHeaderHandler adds given header from request's header as a field to
// the context's logger using fieldKey as field key.
func CustomHeaderHandler(fieldKey, header string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if val := r.Header.Get(header); val != "" {
				log := zerolog.Ctx(r.Context())
				log.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str(fieldKey, val)
				})
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AccessHandler returns a handler that call f after each request.
func AccessHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := ctx.Request.Context()

		// Request
		start := time.Now()
		var buf bytes.Buffer
		tee := io.TeeReader(ctx.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		ctx.Request.Body = ioutil.NopCloser(&buf)

		log.Ctx(c).Debug().Str("request", string(body)).Msg("收到请求")
		ww := &bodyLogWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = ww

		ctx.Next()

		// Response
		log.Ctx(c).Debug().
			Str("response", ww.body.String()).
			Int("status", ww.Status()).
			Dur("duration", time.Since(start)).
			Msg("返回响应")
	}
}
