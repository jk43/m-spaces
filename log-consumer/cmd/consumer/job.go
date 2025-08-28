package main

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

func (app *Application) Job(msg *sarama.ConsumerMessage) {
	//json to map
	val := make(map[string]any)
	err := json.Unmarshal(msg.Value, &val)
	if err != nil {
		app.Logger.Error().Err(err).Send()
	}
	//get rid of value from json
	output, _ := json.Marshal(val["value"])
	fmt.Println(string(output))
}
