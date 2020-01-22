package main

import "MQLearning/RabbitMQ"

func main()  {
	rabbitmqTopicOne :=RabbitMQ.NewRabbitMQTopic("TopicModel","#")
	rabbitmqTopicOne.RecieveTopic()
}
