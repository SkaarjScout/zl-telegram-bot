package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/spf13/pflag"

	"github.com/SkaarjScout/zl-telegram-bot/bothandler"
)

func main() {
	configFileName := pflag.StringP("config", "c", "config.yaml", "Configuration file path")
	pflag.Parse()

	configYaml, err := ioutil.ReadFile(*configFileName)
	if err != nil {
		log.Panicf("Error on config file read: %v", err)
	}
	configYamlExpanded := os.ExpandEnv(string(configYaml))
	config := Config{}
	if err := yaml.Unmarshal([]byte(configYamlExpanded), &config); err != nil {
		log.Panicf("Error on config file read: %v", err)
	}

	db := GetPostgres(config.PostgresConfig)

	ctx, cancelFunc := context.WithCancel(context.Background())
	bot := bothandler.New(ctx, config.TelegramBotConfig, db)

	var wg sync.WaitGroup
	bot.StartServe(&wg)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
	cancelFunc()
	wg.Wait()
}
