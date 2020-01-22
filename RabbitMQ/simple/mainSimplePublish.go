package main

import (
	"MQLearning/RabbitMQ"
	"fmt"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("SimpleModel")
	rabbitmq.PublishSimple("[RabbitMQ-Simple] Hello SimpleModel!")
	fmt.Println("[RabbitMQ-Simple] 发送成功！")
}
