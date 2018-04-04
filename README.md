AMQP Experiments with GO
------------------------

# Configuration
It is almost always better to store configuration inside environment variables:

    - RABBITMQ_VIRTUAL_HOST
    - RABBITMQ_USER
    - RABBITMQ_PASSWORD
    - RABBITMQ_PORT
    - RABBITMQ_HOST

<br />

If this is not possible for you, use passing command line arguments:
``` bash
-vhost=my-vhost -user=my-user -password=my-password -port=my-port -host=my-host
```

# Execution
To run the consumer enter the directory and hit:
```bash
go run consumer.go -exchange=my-exchange -queue=my-queue -tag=my-consumer
```