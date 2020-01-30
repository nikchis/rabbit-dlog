// Copyright (c) 2020 Nikita Chisnikov <chisnikov@gmail.com>
// Distributed under the MIT/X11 software license

package main

const (
	msgErrDial     = `Failed to connect to RabbitMQ`
	msgErrConfig   = `Failed to read a config file`
	msgErrChannel  = `Failed to open a channel`
	msgErrExchange = `Failed to declare an exchange`
	msgErrQueue    = `Failed to declare a queue`
	msgErrBind     = `Failed to bind a queue`
	msgErrConsumer = `Failed to register a consumer`
	msgErrCancel   = `Failed to cancel a consumer`
	msgErrMessage  = `Failed to receive message`
	msgErrEnv      = `Failed to obtain environment variable`
)
