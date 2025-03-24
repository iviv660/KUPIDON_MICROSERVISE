package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	Writer *kafka.Writer
}

func New(brokers []string, topic string) (*Client, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}

	c := Client{}

	c.Writer = &kafka.Writer{
		Addr:                   kafka.TCP(brokers[0]),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &c, nil
}

func (c *Client) SendMessage(ctx context.Context, message string) error {
	if message == "" {
		return errors.New("пустое сообщение не может быть отправлено в Kafka")
	}

	fmt.Println(message)
	msg := kafka.Message{
		Value: []byte(message),
	}

	err := c.Writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		return err
	}

	log.Printf("Сообщение отправлено в Kafka: %s", message)
	return nil
}

// import (
// 	"github.com/IBM/sarama"
// )

// type KafkaProducer struct {
// 	producer sarama.SyncProducer
// }

// func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
// 	config := sarama.NewConfig()
// 	config.Producer.Return.Successes = true
// 	producer, err := sarama.NewSyncProducer(brokers, config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &KafkaProducer{producer: producer}, nil
// }

// func (kp *KafkaProducer) SendMessage(topic string, message string) error {
// 	kafkaMessage := &sarama.ProducerMessage{
// 		Topic: topic,
// 		Value: sarama.StringEncoder(message),
// 	}
// 	_, _, err := kp.producer.SendMessage(kafkaMessage)
// 	return err
// }

// func (kp *KafkaProducer) Close() error {
// 	return kp.producer.Close()
// }
