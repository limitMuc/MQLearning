package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// 定义连接信息
// ur1 格式 amqp（golang规定）://账号：密码@rabbitmq服务器地址：端口号/vhost
const MQURL = "amqp://adminuser:adminuser@127.0.0.1:5672/admin"

//rabbitMQ结构体
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange  string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl     string
}

//创建结构体实例
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	return &RabbitMQ{QueueName:queueName,Exchange:exchange,Key:key,Mqurl:MQURL}
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("[RabbitMQ] %s:%s", message, err)
		panic(fmt.Sprintf("[RabbitMQ] %s:%s", message, err))
	}
}

//创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName,"","")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "[RabbitMQ] failed to connect rabb"+
		"itmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "[RabbitMQ] failed to open a channel")
	return rabbitmq
}

//简单模式下生产者
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建,保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化(当重启后所有数据都会消失，不会持久化保存)
		false,
		//是否自动删除(当最后一个消费者断开连接后，是否会将消息从队列中删除)
		false,
		//是否具有排他性(如果为true，其它消费者就不能访问，只有它自身创建的自身可见)
		false,
		//是否阻塞处理(发送消息以后，是否等待服务器响应)
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则,如果无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),		// 将发送的消息转成字节数组
		})
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,		// durable
		//是否自动删除
		false,		// autoDelete
		//是否具有排他性
		false,		// exclusive
		//是否阻塞处理
		false,		// noWait
		//额外的属性
		nil,			// args
	)
	if err != nil {
		fmt.Println(err)
	}

	//接收消息
	msgs, err :=r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"",     // consumer
		//是否自动应答(消费者将消息消费完后是否自动告诉rabbitmq服务器这个消息已经消费了，默认为true，消费完rabbitmq服务可以将这个消息删除，如果为false，那么autoAck的回调函数需要我们自己手动实现)
		true,   // autoAck
		//是否独有
		false,  // exclusive
		//如果设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false,  // noLocal
		//队列是否阻塞
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("[RabbitMQ] Received a message: %s", d.Body)

		}
	}()

	log.Printf("[RabbitMQ] Waiting for messages. To exit press CTRL+C")
	<-forever
}