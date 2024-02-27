package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn *Connection
}

func NewProducer(conn *Connection) *Producer {
	return &Producer{conn}
}

func (p *Producer) Produce(exchange string, routingKey string, msg amqp.Publishing) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.conn.ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		msg,
	)

	if err != nil {
		return fmt.Errorf("p.conn.ch.PublishWithContext: %w", err)
	}

	return nil
}

// func main() {
// 	q, err := ch.QueueDeclare(
// 		"hello", // name
// 		true,    // устойчивость очереди к рестарту брокера
// 		false,   // delete when unused
// 		false,   // exclusive
// 		false,   // no-wait
// 		nil,     // arguments
// 	)
// 	if err != nil {
// 		log.Fatalf("failed to declare a queue. Error: %s", err)
// 	}
// }
