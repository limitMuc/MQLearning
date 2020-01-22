package main

import "MQLearning/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("SimpleModel")
	rabbitmq.ConsumeSimple()
}
