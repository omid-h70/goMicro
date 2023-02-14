package event

import(
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go" // For sake Of Compatibility With Old Versions
	"helpers"
)

type Consumer struct{
	conn *amqp.Connection
}

func NewConsumer(conn *amqp.Connection)(Consumer, error){
	consumer := Consumer{
		conn,
	}

	err := consumer.setup()
	if err!=nil{
		return Consumer{}, err
	}

	return consumer, nil
}

func (c*Consumer) setup() error{
	return Setup(c.conn)
}

func (c *Consumer) Listen(topics[] string) error{
	ch, err := c.conn.Channel()
	if err!=nil{
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err!=nil{
		return err
	}

	for _,topic := range topics{
		err := ch.QueueBind(
			q.Name,
			topic,
			"logs_topic",
			false,
			nil,
			)

		if err!= nil{
			break
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil{
		return  err
	}

	/*Handle Them in a Separate Routine*/
	go func() {
		for message := range messages {
			var payLoad helpers.PayLoad
			_ = json.Unmarshal(message.Body, &payLoad)

			go handlePayLoad(payLoad)
		}
	}()
	/*Handle Them in a Separate Routine*/
	fmt.Printf("Ready For Exchange Messages")

	stayForever := make(chan bool)
	<-stayForever

	return nil
}

func handlePayLoad(p helpers.PayLoad){
	switch p.Name {
	case "auth":
		//Call Auth Service

	default:
		err := logEvent(p)
		if err!= nil{
			fmt.Println(err)
		}
	}
}

func logEvent(p helpers.PayLoad) error{
	return  nil
}