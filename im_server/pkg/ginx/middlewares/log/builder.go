package log

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type MiddlewaresBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewMiddlewaresLoggerBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewaresBuilder {
	return &MiddlewaresBuilder{
		loggerFunc: fn,
	}
}

func (b *MiddlewaresBuilder) AllowReqBody() *MiddlewaresBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewaresBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader
			ctx.Request.GetBody = func() (io.ReadCloser, error) {
				return reader, nil
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(ctx, al)
		}()

		// 执行到业务逻辑
		ctx.Next()
	}
}

func (b *MiddlewaresBuilder) AllowRespBody() *MiddlewaresBuilder {
	b.allowRespBody = true
	return b
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) WriteHeader(status int) {
	w.al.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.al.RespBody = string(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.WriteString(data)
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// Url 整个请求
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}
