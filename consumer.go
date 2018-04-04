package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	"sync"
)

/**
 * Configuration required for consumer to work
 */
var (
	serve = flag.Bool("serve", false, "Need a server on Server on 8000?")

	vhost    = flag.String("vhost", os.Getenv("RABBITMQ_VIRTUAL_HOST"), "RabbitMQ: Virtual host")
	user     = flag.String("user", os.Getenv("RABBITMQ_USER"), "RabbitMQ: User")
	password = flag.String("password", os.Getenv("RABBITMQ_PASSWORD"), "RabbitMQ: Password")
	port     = flag.String("port", os.Getenv("RABBITMQ_PORT"), "RabbitMQ: Port")
	host     = flag.String("host", os.Getenv("RABBITMQ_HOST"), "RabbitMQ: Host")

	exchange = flag.String("exchange", "", "RabbitMQ: Exchange")
	queue    = flag.String("queue", "", "RabbitMQ: Queue")
	tag      = flag.String("tag", "", "RabbitMQ: Tag")

	MessagesConsumed = []Message{}
)

/**
 * Message structure
 */
type Message struct {
	Application string
	Method      string
	Resource    string
}

/**
 * Parsing command line arguments
 */
func init() {
	flag.Parse()
}

/**
 * Main execution
 */
func main() {
	if *serve {
		Parallelize(server, consume)
	} else {
		consume()
	}
}

/**
 * Parallel execution
 */
func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

func consume() {
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

	Connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", *user, *password, *host, *port, *vhost))
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

	/**
	 * Looping through Deliveries
	 */
	for Delivery := range Deliveries {

		Body := []byte(Delivery.Body)
		Message := Message{}

		if err := json.Unmarshal(Body, &Message); err != nil {
			log.Println("Message cannot be decoded.")
			Delivery.Ack(false)
			continue
		}

		fmt.Println(fmt.Sprintf("Appliction: %s", Message.Application))
		fmt.Println(fmt.Sprintf("Resource: %s", Message.Resource))
		fmt.Println(fmt.Sprintf("Method: %s", Message.Method))

		fmt.Println("")

		MessagesConsumed = append(MessagesConsumed, Message)

		Delivery.Ack(false)
	}
}

/**
 * Setup server
 */
func server() {
	http.HandleFunc("/", homeAction)
	http.ListenAndServe(":8000", nil)
}

/**
 * Processing home page
 */
func homeAction(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, fmt.Sprintf(
		"<h1>For now we got %d messages in scope of this process</h1> "+
			"<h2>Refresh to see if more</h2>",
		len(MessagesConsumed)))

	for _, Message := range MessagesConsumed {
		fmt.Fprint(
			writer,
			fmt.Sprintf(
				"<li>Application: %s, Resource: %s, Method: %s,</li>",
				Message.Application,
				Message.Resource,
				Message.Method,
			),
		)
	}
}
