## Kafka

### Send basic message

```

//create basic message struct
dataToKafka := struct {
  Name string `json:"name"`
} {
  Name: "John"
}
msg := NewMessage(dataToKafka, TopicTest)
err := SendMessage(msg)

```

### Send Log

```
KafkaLog().Level("info").Msg("Hello").TraceInfo(traceInfo).Str("hello", "gogo").Str("nana", "meme").Send()

```

### Send History

```
history := kafka.NewMessage(row.User, kafka.TopicHistory)
kafka.SendMessage(history)

```

### Send SES Email

### Send SES Template Email
