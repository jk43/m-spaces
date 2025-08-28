package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"log"

	"github.com/Shopify/sarama"
)

var kafkaConn sarama.SyncProducer

type TopicName = string

const (
	TopicTest    TopicName = "test"
	TopicLog     TopicName = "log"
	TopicEmail   TopicName = "email"
	TopicHistory TopicName = "history"
)

// Base struct for sending to kafka
type Message struct {
	Topic TopicName `json:"-"`
	Value any       `json:"value"`
}

func (km *Message) GetTopic() TopicName {
	return km.Topic
}

// Message to send AWS SES
type SESTemplate struct {
	Template     string            `json:"template"`
	TemplateData map[string]string `json:"templateData"`
	ToAddresses  []string          `json:"toAddresses"`
	Source       string            `json:"source"`
}

func SendEmail(data *SESTemplate) {
	msg := NewMessage(data, TopicEmail)
	SendMessage(msg)
}

// Log to kafka
type MessageInterface interface {
	GetTopic() TopicName
}

// for singleton
var kafkaLogMessage *LogMessage
var once sync.Once

// sending general log to kafka
// KafkaLog().Level("info").Msg("Hello").TraceInfo(traceInfo).Str("hello", "gogo").Str("nana", "meme").Send()
func KafkaLog() *LogMessage {
	once.Do(func() {
		kafkaLogMessage = &LogMessage{}
	})

	_, file, no, ok := runtime.Caller(1)
	if ok {
		kafkaLogMessage.Caller = fmt.Sprintf("%s:%d", file, no)
	}
	kafkaLogMessage.ServiceName = os.Getenv("SERVICE_NAME")
	kafkaLogMessage.LogStr = make(map[string]string)
	kafkaLogMessage.Time = time.Now().Format(time.RFC3339)
	return kafkaLogMessage
}

type LogLevel string

const (
	LogInfoLevel  LogLevel = "Info"
	LogDebugLevel LogLevel = "Debug"
	LogPanicLevel LogLevel = "Panic"
	LogFatalLevel LogLevel = "Fatal"
	LogErrorLevel LogLevel = "Error"
	LogWarnLevel  LogLevel = "Warn"
	LogTraceLevel LogLevel = "Trace"
)

type LogMessage struct {
	LogLevel     LogLevel          `json:"level"`
	LogMessage   string            `json:"message"`
	LogTraceInfo any               `json:"traceInfo"`
	LogStr       map[string]string `json:"str"`
	Caller       string            `json:"caller"`
	ServiceName  string            `json:"service"`
	Time         string            `json:"logTime"`
}

func (m *LogMessage) Level(input LogLevel) *LogMessage {

	m.LogLevel = input
	return m
}

func (m *LogMessage) Message(input string) *LogMessage {
	m.LogMessage = input
	return m
}

func (m *LogMessage) TraceInfo(input any) *LogMessage {
	m.LogTraceInfo = input
	return m
}

func (m *LogMessage) Str(k, v string) *LogMessage {
	m.LogStr[k] = v
	return m
}

func (m *LogMessage) Send() {
	SendLog(m)
}

func SendLog(data any) {
	msg := NewMessage(data, TopicLog)
	SendMessage(msg)
}

// Create a new kafka message
func NewMessage(data any, topic TopicName) *Message {
	msg := Message{
		Topic: topic,
		Value: data,
	}
	return &msg
}

func SendMessage(msg MessageInterface) error {
	if kafkaConn == nil {
		kafkaAddr := kafkaAddrs()
		conn, err := sarama.NewSyncProducer(kafkaAddr, nil)
		if err != nil {
			return err
		}
		kafkaConn = conn
	}

	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	go func(payload []byte) {
		_, _, err := kafkaConn.SendMessage(&sarama.ProducerMessage{
			Topic: msg.GetTopic(),
			Value: sarama.StringEncoder(payload),
		})
		if err != nil {
			log.Print(err)
		}
	}(json)

	return nil
}

func kafkaAddrs() []string {
	kafkaAddr := os.Getenv("KAFKA_LB_ADDR")
	return strings.Split(kafkaAddr, ",")
}

//
//kafka consumer
//

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

var (
	brokers  []string
	version  string
	group    string
	topics   string
	assignor string
	oldest   bool
	verbose  bool
)

// Initialize the Kafka consumer setup
func initKafkaConsumerSetup() {
	brokers = kafkaAddrs()                                     // Get the Kafka addresses
	version = os.Getenv("KAFKA_VERSION")                       // Get the Kafka version from the environment variable
	group = os.Getenv("KAFKA_GROUP")                           // Get the Kafka group from the environment variable
	topics = os.Getenv("KAFKA_TOPIC")                          // Get the Kafka topics from the environment variable
	assignor = os.Getenv("KAFKA_ASSIGNOR")                     // Get the Kafka assignor from the environment variable
	oldest, _ = strconv.ParseBool(os.Getenv("KAFKA_OLDEST"))   // Get the Kafka oldest flag from the environment variable
	verbose, _ = strconv.ParseBool(os.Getenv("KAFKA_VERBOSE")) // Get the Kafka verbose flag from the environment variable
	log.Println("kafka config -", "brokers:", brokers, "version:", version, "group:", group, "topics:", topics, "assignor:", assignor, "oldest:", oldest, "verbose:", verbose)
}

type ConsumerWithChan interface {
	GetChan() *chan bool
	sarama.ConsumerGroupHandler
}

func ConsumeMessage(consumer ConsumerWithChan) {
	initKafkaConsumerSetup()
	keepRunning := true

	KafkaLog().Level("info").Message("Starting a new Sarama consumer").Send()
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumerChan := consumer.GetChan()

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split(topics, ","), consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			*consumerChan = make(chan bool)
		}
	}()

	<-*consumerChan // Await till the consumer has been set up
	KafkaLog().Level("info").Message("Sarama consumer up and running!...").Send()

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			KafkaLog().Level("info").Message("terminating: context cancelled").Send()
			keepRunning = false
		case <-sigterm:
			KafkaLog().Level("info").Message("terminating: via signal").Send()
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
