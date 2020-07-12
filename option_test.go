package gmicro

import (
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestStaticDir(t *testing.T) {
	s := NewService(WithStaticDir("/a/b/c"))
	assert.Equal(t, "/a/b/c", s.staticDir)
}

func TestAnnotator(t *testing.T) {
	s := NewService(
		WithAnnotator(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.New(nil)
			md.Set("key", "value")
			return md
		}),
	)

	assert.Len(t, s.annotators, 1)
}

func TestErrorHandler(t *testing.T) {
	s := NewService(WithErrorHandler(nil))
	assert.Nil(t, s.errorHandler)
}

func TestHTTPHandler(t *testing.T) {
	s := NewService(WithHTTPHandler(nil))
	assert.Nil(t, s.httpHandler)
}

func TestUnaryInterceptor(t *testing.T) {
	s := NewService(
		WithUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (resp interface{}, err error) {
			return nil, nil
		}),
	)

	assert.Len(t, s.unaryInterceptors, 3)
}

func TestStreamInterceptor(t *testing.T) {
	s := NewService(
		WithStreamInterceptor(func(srv interface{}, stream grpc.ServerStream,
			info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return nil
		}),
	)

	assert.Len(t, s.streamInterceptors, 3)
}

func TestInterruptSignal(t *testing.T) {
	s := NewService(
		WithInterruptSignal(syscall.SIGKILL),
	)

	assert.Len(t, s.interruptSignals, 7)
}

func TestGRPCServerOption(t *testing.T) {
	s := NewService(
		WithGRPCServerOption(grpc.ConnectionTimeout(10 * time.Second)),
	)

	assert.Len(t, s.gRPCServerOptions, 3)
}

func TestGRPCDialOption(t *testing.T) {
	s := NewService(
		WithGRPCDialOption(grpc.WithBlock()),
	)

	assert.Len(t, s.gRPCDialOptions, 1)
}

func TestWithHTTPServer(t *testing.T) {
	s := NewService(WithHTTPServer(&http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}))

	assert.NotNil(t, s.HTTPServer)
	assert.Equal(t, 5*time.Second, s.HTTPServer.ReadTimeout)
}

func TestMuxOption(t *testing.T) {
	s := NewService(
		WithMuxOption(runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{OrigName: true, EmitDefaults: true},
		)),
	)

	assert.Len(t, s.muxOptions, 3)
}
