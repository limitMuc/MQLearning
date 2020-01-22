package main

import (
	"MQLearning/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqRoutingOne := RabbitMQ.NewRabbitMQRouting("RoutingModel", "routing_one")
	rabbitmqRoutingTwo := RabbitMQ.NewRabbitMQRouting("RoutingModel", "routing_two")
	for i := 0; i <= 10; i++ {
		rabbitmqRoutingOne.PublishRouting("[RabbitMQ-Routing] Hello RoutingModel one!" + strconv.Itoa(i))
		rabbitmqRoutingTwo.PublishRouting("[RabbitMQ-Routing] Hello RoutingModel Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
