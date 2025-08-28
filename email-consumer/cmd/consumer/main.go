// kafka-topics.sh --bootstrap-server 192.168.1.55:9094 --create --topic log --partitions 3
package main

import (
	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/kafka"
	"github.com/rs/zerolog"
)

// func init() {
// 	os.Setenv("KAFKA_VERSION", "2.1.1")
// 	os.Setenv("KAFKA_GROUP", "main")
// 	os.Setenv("KAFKA_TOPIC", "logs")
// 	os.Setenv("KAFKA_ASSIGNOR", "range")
// 	os.Setenv("KAFKA_OLDEST", "true")
// 	os.Setenv("KAFKA_VERBOSE", "false")
// 	os.Setenv("KAFKA_LB_ADDR", "192.168.1.55:9094")
// }

type Application struct {
	Logger *zerolog.Logger
}

var app Application

func main() {
	//logger
	app.Logger = molylibs.Logger()

	consumer := Consumer{
		ready: make(chan bool),
		job:   app.Job,
	}
	kafka.ConsumeMessage(&consumer)
}
