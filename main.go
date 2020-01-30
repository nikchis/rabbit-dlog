// Copyright (c) 2020 Nikita Chisnikov <chisnikov@gmail.com>
// Distributed under the MIT/X11 software license

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/streadway/amqp"
)

type logLevel int

const (
	logInfo logLevel = iota
	logWarn
	logError
)

func checkErrFatal(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err.Error())
	}
}

func checkErrPrint(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err.Error())
	}
}

func getEnvFatal(env string) (res string) {
	var exist bool
	if res, exist = os.LookupEnv(env); !exist {
		log.Fatalf("%s: %s", env, msgErrEnv)
	}
	return res
}

func main() {
	var wg sync.WaitGroup

	fDir := flag.String("dir", "", "set custom app files directory")
	flag.Parse()

	cfg, err := newConfig(*fDir)
	checkErrFatal(err, msgErrConfig)

	c, err := amqp.Dial(cfg.Url)
	checkErrFatal(err, msgErrDial)
	defer c.Close()
	ch, err := c.Channel()
	checkErrFatal(err, msgErrChannel)
	defer ch.Close()
	err = ch.ExchangeDeclare(cfg.Exchange, "direct", true, true, false, false, nil)
	checkErrFatal(err, msgErrExchange)

	rInfo := newRoute(ch, cfg.Exchange, cfg.QueInfo, cfg.KeyInfo, cfg.ConInfo,
		cfg.Id, cfg.Token, logInfo, &wg)
	rWarn := newRoute(ch, cfg.Exchange, cfg.QueWarn, cfg.KeyWarn, cfg.ConWarn,
		cfg.Id, cfg.Token, logWarn, &wg)
	rError := newRoute(ch, cfg.Exchange, cfg.QueError, cfg.KeyError, cfg.ConError,
		cfg.Id, cfg.Token, logError, &wg)

	rInfo.run()
	rWarn.run()
	rError.run()

	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ossig

	rInfo.stop()
	rWarn.stop()
	rError.stop()
	wg.Wait()
}
