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
-vhost=myvhost -user=myuser -password=mypassword -port=myport -host=myhost
```

# Execution
To run the consumer enter the directory and hit:
```bash
go run consumer.go -exchange=myexchange -queue=myqueue -tag=myconsumer
```

# Server
You can run http server on 8000 port by adding `-server=true`:
<br /><br />
Example:
```bash
go run consumer.go -exchange=myexchange -queue=myqueue -tag=myconsumer -server=true
```

In this case you can see the current amount of consumed messages visiting `http://127.0.0.1:8000`