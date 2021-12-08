package main

import (
	"database/sql"
	"encoding/json"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mateuslima90/payment-gateway/adapter/broker/kafka"
	"github.com/mateuslima90/payment-gateway/adapter/factory"
	"github.com/mateuslima90/payment-gateway/adapter/presenter/transaction"
	"github.com/mateuslima90/payment-gateway/usecase/process_transaction"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//db
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal(err)
	}

	//repository
	repositoryFactory := factory.NewRepositoryDatabaseFactory(db)
	repository := repositoryFactory.CreateTransactionRepository()

	//configMapProducer
	configMapProducer := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	}

	//producer
	kafkaPresenter := transaction.NewTransactionKafkaPresenter()
	producer := kafka.NewKafkaProducer(configMapProducer, kafkaPresenter)

	//configMapConsumer
	var msgChan = make(chan *ckafka.Message)
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"client.id":"goapp",
		"group.id":"goapp",
	}

	//topic
	topics := []string{"trnsactions"}

	//consumer
	consumer := kafka.NewConsumer(configMapConsumer, topics)
	go consumer.Consume(msgChan)

	//usecase
	usecase := process_transaction.NewProcessTransaction(repository, producer, "transactions_result")

	for msg := range msgChan {
		var input process_transaction.TransactionDtoInput
		json.Unmarshal(msg.Value, &input)
		usecase.Execute(input)
	}
}
