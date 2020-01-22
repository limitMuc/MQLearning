package main

import "MQLearning/RabbitMQ/simple"

func main() {
	rabbitmq := simple.NewRabbitMQSimple("SimpleModel")
	rabbitmq.ConsumeSimple()
}
