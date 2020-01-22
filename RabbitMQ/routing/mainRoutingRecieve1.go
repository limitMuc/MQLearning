package main

import "MQLearning/RabbitMQ"

func main()  {
	rabbitmqRoutingOne :=RabbitMQ.NewRabbitMQRouting("RoutingModel","routing_one")
	rabbitmqRoutingOne.RecieveRouting()
}
