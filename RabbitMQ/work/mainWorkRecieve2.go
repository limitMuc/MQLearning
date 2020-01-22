package main

import "MQLearning/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("WorkModel")
	rabbitmq.ConsumeSimple()
}
