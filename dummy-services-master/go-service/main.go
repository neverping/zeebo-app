package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/joeshaw/envdecode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/sybogames/dummy-services/helloworld"
)

type config struct {
	Port            uint32        `env:"PORT,default=4458"`
	ServiceEndpoint string        `env:"SERVICE_ENDPOINT,default=localhost:50051"`
	RequestTimeout  time.Duration `env:"REQUEST_TIMEOUT,default=2"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT,default=4"`
}

type HelloRequest struct {
	Name string `json:"name"`
}

func main() {
	var cfg config
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log, _ := config.Build()

	defer log.Sync()
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
	}
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	}

	svconn, err := grpc.Dial(cfg.ServiceEndpoint, dialOpts...)
	svclient := helloworld.NewGreeterClient(svconn)
	if err != nil {
		log.Fatal("failed to dial service", zap.Error(err))
	}

	router := chi.NewRouter()

	public := chi.Middlewares{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Heartbeat("/"),
		middleware.Heartbeat("/heartbeat"),
		render.SetContentType(render.ContentTypeJSON),
		middleware.DefaultCompress,
		middleware.Timeout(60 * time.Second),
	}

	router.Use(
		middleware.Heartbeat("/"),
		middleware.Heartbeat("/heartbeat"),
	)

	router.Group(func(r chi.Router) {
		r.Use(public...)
		r.Post("/dummy/hello", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, struct {
				Response string `json:"response"`
			}{Response: "OK"})
		})
		r.Post("/dummy/service", func(w http.ResponseWriter, r *http.Request) {
			var hrq HelloRequest
			err := render.Decode(r, &hrq)
			if err != nil {
				log.Error("Error decoding request", zap.Error(err))
				render.Status(r, 500)
				return
			}
			reply, err := svclient.SayHello(context.Background(), &helloworld.HelloRequest{Name: hrq.Name})
			if err != nil {
				log.Error("Error calling GRPC service", zap.Error(err))
				render.Status(r, 500)
				return
			}
			render.JSON(w, r, struct {
				Response string `json:"response"`
			}{Response: reply.Message})
		})
	})

	// Setup interruption signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	defer signal.Stop(interrupt)

	// Create a context base context we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup an error group to run our servers in
	group, ctx := errgroup.WithContext(ctx)

	var httpSrv *http.Server

	group.Go(func() error {
		httpSrv = &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			ReadTimeout:  cfg.RequestTimeout * time.Second,
			WriteTimeout: cfg.RequestTimeout * time.Second,
			IdleTimeout:  cfg.IdleTimeout * time.Second,
			Handler:      router,
		}

		log.Info(fmt.Sprintf("listenting to port :%d...", cfg.Port))
		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Fatal("unable to start http server", zap.Error(err))
			return err
		}
		log.Info("Stopped serving Http on GDPR")
		return nil
	})

	// Wait for interrupt or errorgroup context cancel
	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Info("starting graceful shutdown")
	cancel()

	if httpSrv != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = httpSrv.Shutdown(shutdownCtx)
	}

	// wait for error group to complete
	if err := group.Wait(); err != nil {
		log.Fatal("shutdown failed", zap.Error(err))
	}

	log.Info("shutdown complete")

	log.Fatal("server failed to start", zap.Error(err))
}
