package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/nhan1603/IoTsystem/api/internal/appconfig/iam"
	"github.com/nhan1603/IoTsystem/api/internal/controller/auth"
	"github.com/nhan1603/IoTsystem/api/internal/controller/iot"
	authHandler "github.com/nhan1603/IoTsystem/api/internal/handler/rest/public/v1/auth"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/env"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/runner"
)

type router struct {
	ctx      context.Context
	authCtrl auth.Controller
	iotCtrol iot.Controller
}

func (rtr router) initKafkaConsumer() {
	batchSize, _ := strconv.Atoi(env.GetwithDefault("BATCH_SIZE", "100"))
	batchTimeout, _ := time.ParseDuration(env.GetwithDefault("BATCH_TIMEOUT", "5s"))
	// Inital consumer kafka
	consumer, err := kafka.NewBatchConsumer(
		rtr.ctx,
		os.Getenv("IOT_TOPIC"),
		os.Getenv("KAFKA_BROKER"),
		rtr.iotCtrol.HandleBatch,
		"iot",
		batchSize,
		batchTimeout,
	)
	if err != nil {
		log.Printf("Error when init consumer, %v", err)
		return
	}

	log.Printf("Kafka consumer start successfully\n")
	go runner.ExecParallel(rtr.ctx, consumer.Consume)
}

func (rtr router) routes(r chi.Router) {
	r.Group(rtr.authenticated)
	r.Group(rtr.public)
}

func (rtr router) authenticated(r chi.Router) {
	prefix := "/api/authenticated"

	r.Group(func(r chi.Router) {
		r.Use(iam.AuthenticateUserMiddleware(rtr.ctx))
		prefix = prefix + "/v1"

	})
}

func (rtr router) public(r chi.Router) {
	prefix := "/api/public"

	r.Use(middleware.Logger)
	r.Group(func(r chi.Router) {
		r.Get(prefix+"/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
	})

	// v1
	r.Group(func(r chi.Router) {
		prefix = prefix + "/v1"

		r.Group(func(r chi.Router) {
			authH := authHandler.New(rtr.authCtrl)
			r.Post(prefix+"/login", authH.AuthenticateOperationUser())
			r.Post(prefix+"/user", authH.CreateUser())
		})
	})
}
