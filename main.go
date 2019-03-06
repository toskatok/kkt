package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Toskatok/kkt/actions"
	"github.com/Toskatok/kkt/kafka"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("kkt")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Banner Message
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	// Logrus logger setup
	logrus.SetLevel(logrus.Level(viper.GetInt("logger.level")))
	if viper.GetString("env") == "DEBUG" {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// echo server
	e := actions.App(viper.GetString("env") == "DEBUG")
	go func() {
		if err := e.Start(viper.GetString("server.address")); err != http.ErrServerClosed {
			log.Fatalf("API Service failed with %s", err)
		}
	}()

	// kafka consumer
	cg, err := kafka.NewConsumerGroup(
		viper.GetString("kafka.topic"),
		viper.GetStringSlice("kafka.brokers"),
		viper.GetString("kafka.client-id"),
	)
	if err != nil {
		log.Fatalf("Kafka Consumer Service failed with %s", err)
	}
	go cg.Run()

	// Detect kill signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Banner Message
	fmt.Println("18.20 As always ... left me alone")

	// kafka consumer
	cg.Exit()

	// echo server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("API Service failed on exit: %s", err)

	}
}
