package main

import (
	"MQLearning/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("WorkModel")

	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("[RabbitMQ-Work] Hello WorkModel!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
