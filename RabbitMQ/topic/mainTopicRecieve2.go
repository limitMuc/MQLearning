package main

import "MQLearning/RabbitMQ"

func main()  {
	rabbitmqTopicTwo :=RabbitMQ.NewRabbitMQTopic("TopicModel","exrabbitmq.*.two")
	rabbitmqTopicTwo.RecieveTopic()
}
