package main

import "MQLearning/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("SubModel")
	rabbitmq.RecieveSub()
}
