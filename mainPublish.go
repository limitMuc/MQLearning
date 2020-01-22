package main

import (
	"MQLearning/RabbitMQ/simple"
	"fmt"
)

func main() {
	rabbitmq := simple.NewRabbitMQSimple("SimpleModel")
	rabbitmq.PublishSimple("[RabbitMQ-Simple] Hello SimpleModel!")
	fmt.Println("[RabbitMQ-Simple] 发送成功！")
}
