// Copyright (c) 2020 Nikita Chisnikov <chisnikov@gmail.com>
// Distributed under the MIT/X11 software license

package main

import (
	"log"
	"sync"

	"github.com/streadway/amqp"

	discordlog "github.com/nikchis/discord-log"
)

type route struct {
	ch   *amqp.Channel
	que  *amqp.Queue
	cons string
	msgs <-chan amqp.Delivery
	dl   *discordlog.Webhook
	lev  logLevel
	done chan bool
	wg   *sync.WaitGroup
}

func newRoute(ch *amqp.Channel, exchange, queue, key, consumer, id, token string,
	lev logLevel, wg *sync.WaitGroup) *route {
	q, err := ch.QueueDeclare(queue, true, true, false, false, nil)
	checkErrFatal(err, msgErrQueue)
	err = ch.QueueBind(q.Name, key, exchange, false, nil)
	checkErrFatal(err, msgErrBind)
	messages, err := ch.Consume(q.Name, consumer, false, false, false, false, nil)
	checkErrFatal(err, msgErrConsumer)
	done := make(chan bool, 1)
	dlog := discordlog.NewWebhook(id, token)
	r := &route{
		ch:   ch,
		que:  &q,
		msgs: messages,
		cons: consumer,
		dl:   dlog,
		lev:  lev,
		done: done,
		wg:   wg,
	}
	return r
}

func (r *route) run() {
	r.wg.Add(1)
	go func() {
		for {
			select {
			case msg, ok := <-r.msgs:
				if !ok {
					log.Printf("%s: %s", r.cons, msgErrMessage)
				} else {
					switch r.lev {
					case logInfo:
						if err := r.dl.PrintInfo(string(msg.Body)); err == nil {
							msg.Ack(false)
						}
					case logWarn:
						if err := r.dl.PrintWarning(string(msg.Body)); err == nil {
							msg.Ack(false)
						}
					case logError:
						if err := r.dl.PrintError(string(msg.Body)); err == nil {
							msg.Ack(false)
						}
					}
				}
			case <-r.done:
				log.Printf("%s: stopped", r.que.Name)
				r.wg.Done()
				close(r.done)
				return
			}
		}
	}()
}

func (r *route) stop() {
	err := r.ch.Cancel(r.cons, false)
	checkErrPrint(err, msgErrCancel)
	r.done <- true
}
