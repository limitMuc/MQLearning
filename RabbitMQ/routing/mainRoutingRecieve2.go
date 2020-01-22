package main

import "MQLearning/RabbitMQ"

func main()  {
	rabbitmqRoutingTwo :=RabbitMQ.NewRabbitMQRouting("RoutingModel","routing_two")
	rabbitmqRoutingTwo.RecieveRouting()
}
