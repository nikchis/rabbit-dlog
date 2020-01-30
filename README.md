rabbit-dlog [![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/nikchis/rabbit-dlog/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/nikchis/rabbit-dlog?)](https://goreportcard.com/report/github.com/nikchis/rabbit-dlog)
==========

rabbit-dlog is application which allows you to route log text messages from RabbitMQ to Discord channels.
	
## Usage

~~~bash
#building image
docker build -t rabbit-dlog .
#gen config
docker run -i --rm --name rabbit-dlog --network skynet \
-v /var/lib/docker/volumes/rabbit-dlog/data:/app/data rabbit-dlog
#edit config
sudo vim /var/lib/docker/volumes/rabbit-dlog/data/.rabbit-dlog.conf
#run container
docker run --name rabbit-dlog -d --restart always --network skynet \
-v /var/lib/docker/volumes/rabbit-dlog/data:/app/data rabbit-dlog
~~~
