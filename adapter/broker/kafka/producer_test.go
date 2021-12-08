package kafka

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mateuslima90/payment-gateway/adapter/presenter/transaction"
	"github.com/mateuslima90/payment-gateway/domain/entity"
	"github.com/mateuslima90/payment-gateway/usecase/process_transaction"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducerPublish(t *testing.T)  {
	expectedOutput := process_transaction.TransactionDtoOutput{
		ID: "1",
		Status: entity.REJECTED,
		ErrorMessage: "you dont have limit for this transaction",
	}

	//outputJson, _ := json.Marshal(expectedOutput)

	configMap := ckafka.ConfigMap{
		"test.mock.num.brokers": 3,
	}

	producer := NewKafkaProducer(&configMap, transaction.NewTransactionKafkaPresenter())
	err := producer.Publish(expectedOutput, []byte("1"), "test")
	assert.Nil(t, err)
}