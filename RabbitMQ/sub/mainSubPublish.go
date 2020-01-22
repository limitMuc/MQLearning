package main

import (
	"MQLearning/RabbitMQ"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("SubModel")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("[RabbitMQ-Sub] 订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		fmt.Println("[RabbitMQ-Sub] 订阅模式生产第" +
			strconv.Itoa(i) + "条" + "数据")
		time.Sleep(1 * time.Second)
	}

}
