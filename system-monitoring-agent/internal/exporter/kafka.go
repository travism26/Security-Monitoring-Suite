package exporter

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

// mtravis notes: I dont think we need this I dont plan on directly connecting to kafka
// I plan on sending data to a endpoint that is connected to kafka and then passing it to kafka
// I will keep this for now but we may want to remove it...

// KafkaExporter handles the export of monitoring data to Kafka topics
type KafkaExporter struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaExporter creates a new Kafka exporter instance
// brokers is a list of Kafka broker addresses
// topic is the Kafka topic to publish messages to
// Returns the new exporter and any error encountered during setup
func NewKafkaExporter(brokers []string, topic string) (*KafkaExporter, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaExporter{
		producer: producer,
		topic:    topic,
	}, nil
}

// Export sends the provided data to Kafka as a JSON message
// data is a map of values to be exported
// Returns an error if marshaling or sending fails
func (k *KafkaExporter) Export(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(jsonData),
	}

	_, _, err = k.producer.SendMessage(msg)
	return err
}

// Close cleanly shuts down the Kafka producer
// Returns any error encountered during shutdown
func (k *KafkaExporter) Close() error {
	return k.producer.Close()
}
