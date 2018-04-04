package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/streadway/amqp"
	"log"
)

/**
 * Configuration required for consumer to work
 */
var (
	vhost    = flag.String("vhost", os.Getenv("RABBITMQ_VIRTUAL_HOST"), "RabbitMQ: Virtual host")
	user     = flag.String("user", os.Getenv("RABBITMQ_USER"), "RabbitMQ: User")
	password = flag.String("password", os.Getenv("RABBITMQ_PASSWORD"), "RabbitMQ: Password")
	port     = flag.String("port", os.Getenv("RABBITMQ_PORT"), "RabbitMQ: Port")
	host     = flag.String("host", os.Getenv("RABBITMQ_HOST"), "RabbitMQ: Host")

	exchange = flag.String("exchange", "", "RabbitMQ: Exchange")
	queue    = flag.String("queue", "", "RabbitMQ: Queue")
	tag      = flag.String("tag", "", "RabbitMQ: Tag")
)

/**
 * Parsing command line arguments
 */
func init() {
	flag.Parse()
}

func main() {

	if *queue == "" {
		log.Fatal("Queue should be defined by argument (-queue=my_queue)")
	}
	if *exchange == "" {
		log.Fatal("Exchange should be defined by argument (-exchange=my_exchange)")
	}
	if *tag == "" {
		log.Fatal("Tag should be defined by argument (-tag=my_tag)")
	}

	var err error

	Connection, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/%s",
			*user,
			*password,
			*host,
			*port,
			*vhost,
		),
	)

	if err != nil {
		log.Fatalf("%s", err)
	}

	Channel, err := Connection.Channel()

	if err != nil {
		log.Fatalf("%s", err)
	}

	Deliveries, err := Channel.Consume(*queue, *tag, false, false, false, false, nil)

	if err != nil {
		log.Fatalf("%s", err)
	}

	for Delivery := range Deliveries {
		log.Printf(
			"Tag: %v, Body: %s",
			Delivery.DeliveryTag,
			Delivery.Body,
		)

		Delivery.Ack(false)
	}
}
