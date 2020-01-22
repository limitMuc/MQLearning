package main

import (
	"MQLearning/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqTopicOne := RabbitMQ.NewRabbitMQTopic("TopicModel", "exrabbitmq.topic.one")
	rabbitmqTopicTwo := RabbitMQ.NewRabbitMQTopic("TopicModel", "exrabbitmq.topic.two")
	for i := 0; i <= 10; i++ {
		rabbitmqTopicOne.PublishTopic("[RabbitMQ-Topic] Hello TopicModel one!" + strconv.Itoa(i))
		rabbitmqTopicTwo.PublishTopic("[RabbitMQ-Topic] Hello TopicModel Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
