package exporter

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type KafkaExporter struct {
	producer sarama.SyncProducer
	topic    string
}

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

func (k *KafkaExporter) Close() error {
	return k.producer.Close()
}
