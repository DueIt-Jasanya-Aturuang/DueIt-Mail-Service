package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/DueIt-Jasanya-Aturuang/DueIt-Mail-Service/config"
	"github.com/DueIt-Jasanya-Aturuang/DueIt-Mail-Service/internal/logs"
	"github.com/DueIt-Jasanya-Aturuang/DueIt-Mail-Service/internal/template"
	"github.com/DueIt-Jasanya-Aturuang/DueIt-Mail-Service/modules/services"
)

func main() {
	dir, _ := os.Getwd()
	dir = fmt.Sprintf("%s/internal/logs", dir)
	logs.InitLogger(logs.Config{
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJson:      true,
		FileLoggingEnabled:    true,
		Directory:             dir,
		Filename:              "mail.log",
		MaxSize:               200000000,
		MaxBackups:            2000,
		MaxAge:                2000,
	})

	mechanism := plain.Mechanism{
		Username: config.Get().Application.Kafka.User,
		Password: config.Get().Application.Kafka.Pass,
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	config := kafka.ReaderConfig{
		Brokers:  []string{config.Get().Application.Kafka.Broker},
		GroupID:  config.Get().Application.Kafka.Group,
		Topic:    config.Get().Application.Kafka.Topic,
		MaxWait:  500 * time.Millisecond,
		MinBytes: 1,
		MaxBytes: 10e6,
		Dialer:   dialer,
	}
	r := kafka.NewReader(config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		log.Info().Msg("close kafka")
		if err := r.Close(); err != nil {
			log.Info().Msgf("Failed to close reader:%v", err)
		}
		cancel()
		log.Info().Msg("os exit")
		os.Exit(1)
	}()

	template := template.NewEmailTemplateImpl()
	mailSvc := services.NewEmailServiceImpl(template)

	log.Info().Msg("consumer start listen...")
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Err(err).Msgf("cannot read message kafka : %v", err)
			break
		}

		// go func() {
		if err := mailSvc.SendGOMAIL(m.Value); err != nil {
			log.Err(err).Msgf("cannot send mail with smtp : %v", err)
		}
		// }()

		format := fmt.Sprintf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		log.Info().Msg(format)
	}
}
