package kafka

import (
	"os"
	"testing"
)

func Test_kafkaAddrs(t *testing.T) {
	os.Setenv("KAFKA_LB_ADDR", "192.168.1.55:9094")
	addrs := kafkaAddrs()
	if addrs[0] != "192.168.1.55:9094" {
		t.Errorf("Unexpected address")
	}
	os.Setenv("KAFKA_LB_ADDR", "192.168.1.55:9094,192.168.1.52:9094")
	addrs = kafkaAddrs()
	if addrs[0] != "192.168.1.55:9094" || addrs[1] != "192.168.1.52:9094" {
		t.Errorf("Unexpected address")
	}
}

func Test_KafkaLog(t *testing.T) {
	ti := struct {
		Name string `json:"string"`
	}{
		"golang",
	}
	var level LogLevel = "info"
	message := "message"
	klm := KafkaLog()
	KafkaLog().Level(level).Message(message).TraceInfo(ti).Str("name", "alex")
	if klm.LogLevel != level {
		t.Errorf("Expect string '%s' in level but got a %s", level, klm.LogLevel)
	}
	if klm.LogMessage != message {
		t.Errorf("Expect string '%s' in message but got a %s", message, klm.LogMessage)
	}
	if klm.LogStr["name"] != "alex" {
		t.Errorf("Expect string '%s' in str but got a %s", "alex", klm.LogStr["name"])
	}
}
