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
	//创建RabbitMQ实例,在 RabbitMQ 中，消息是不能直接发送到队列，它需要发送到交换机（exchange）中，它可以使用一个空字符串来标识。交换机允许我们指定某条消息需要投递到哪个队列。
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
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建,保证队列存在，消息能发送到队列中, 如果我们把消息发送到一个不存在的队列，RabbitMQ 会丢弃这条消息
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化(当重启后所有数据都会消失，不会持久化保存)
		false,		// durable
		//是否自动删除(当最后一个消费者断开连接后，是否会将消息从队列中删除)
		false,		// autoDelete
		//是否具有排他性(如果为true，其它消费者就不能访问，只有它自身创建的自身可见)
		false,		// exclusive
		//是否阻塞处理(发送消息以后，是否等待服务器响应)
		false,		// noWait
		//额外的属性
		nil,			// args
	)
	if err != nil {
		fmt.Println(err)
	}
	//调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则,如果无法找到符合条件的队列会把消息返还给发送者
		false,		// mandatory
		//如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
		false,		// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),		// 将发送的消息转成字节数组
		})
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	// 为什么要重复声明队列呢 —— 我们已经在前面的代码中声明过它了。如果我们确定了队列是已经存在的，那么我们可以不这么做，可是我们并不确定哪个程序会首先运行。这种情况下，在程序中重复将队列重复声明一下是种值得推荐的做法.
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
		false,  	// noWait
		nil,    	// args
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


//订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("",exchangeName,"")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err,"[RabbitMQ] failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "[RabbitMQ] failed to open a channel")
	return rabbitmq
}

//订阅模式生产者
func (r *RabbitMQ) PublishPub(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型，订阅模式下要将类型设置为fanout广播类型
		"fanout",		// kind
		true,		// durable
		false,		// autoDelete
		//如果true, 表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,		// internal
		false,		// noWait
		nil,			// args
	)

	r.failOnErr(err, "[RabbitMQ] Failed to declare an exchange")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",				// key
		false,		// mandatory
		false,		// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//订阅模式消费者
func (r *RabbitMQ) RecieveSub() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"fanout",		// kind
		true,		// durable
		false,		// autoDelte
		//如果true, 表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,		// internal
		false,		//noWait
		nil,			// args
	)
	r.failOnErr(err, "[RabbitMQ] Failed to declare an exchange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"",		 //留空让其随机生产队列名称
		false,		// durable
		false,		// autoDelete
		true,		// exclusive
		false,		// noWait
		nil,			// args
	)
	r.failOnErr(err, "[RabbitMQ] Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key必须要为空
		"",		// key
		r.Exchange,
		false,		// noWait
		nil,
		)			// args

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",		// consumer
		true,		// autoAck
		false,		// exclusive
		false,		// noLocal
		false,		// noWait
		nil,			// args
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("[RabbitMQ] Received a message: %s", d.Body)
		}
	}()

	fmt.Println("[RabbitMQ] Waiting for messages. To exit press CTRL+C")
	<-forever
}


//
//创建RabbitMQ路由模式实例
func NewRabbitMQRouting(exchangeName string,routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err,"[RabbitMQ] failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "[RabbitMQ] failed to open a channel")
	return rabbitmq
}

//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string )  {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//要改成direct
		"direct",		// kind
		true,		// durable
		false,		// autoDelete
		false,		// internal
		false,		// noWait
		nil,			// args
	)

	r.failOnErr(err, "[RabbitMQ] Failed to declare an exchange")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//一定要设置key
		r.Key,
		false,		// mandatory
		false,		// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
//路由模式接受消息
func (r *RabbitMQ) RecieveRouting() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"direct",		// kind
		true,		// durable
		false,		// autoDelete
		false,		// internal
		false,		// noWait
		nil,			// args
	)
	r.failOnErr(err, "[RabbitMQ] Failed to declare an exchange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,		// durable
		false,		// autoDelete
		true,		// exclusive
		false,		// noWait
		nil,			// args
	)
	r.failOnErr(err, "[RabbitMQ] Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//需要绑定key
		r.Key,
		r.Exchange,
		false,		// noWait
		nil)			// args

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",		// consumer
		true,		// autoAck
		false,		// exclusive
		false,		// noLocal
		false,		// noWait
		nil,			// args
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("[RabbitMQ] Received a message: %s", d.Body)
		}
	}()

	fmt.Println("[RabbitMQ] Waiting for messages. To exit press CTRL+C")
	<-forever
}



//创建RabbitMQ话题模式实例
func NewRabbitMQTopic(exchangeName string,routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err,"[RabbitMQ] failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "[RabbitMQ] failed to open a channel")
	return rabbitmq
}
//话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string )  {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//要改成topic
		"topic",			// kind
		true,			// durable
		false,			// autoDelete
		false,			// internal
		false,			// noWait
		nil,				// args
	)

	r.failOnErr(err, "[RabbitMQ] Failed to declare an exchange")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//要设置key
		r.Key,
		false,			// mandatory
		false,			// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 exrabbitmq.* 表示匹配 exrabbitmq.hello, exrabbitmq.hello.one需要用exrabbitmq.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"topic",			// kind
		true,			// durable
		false,			// autoDelete
		false,			// internal
		false,			// noWait
		nil,				// args
	)
	r.failOnErr(err, "[RabbitMQ] Failed to declare an exchange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,		// durable
		false,		// autoDelete
		true,		// exclusive
		false,		// noWait
		nil,			// args
	)
	r.failOnErr(err, "[RabbitMQ] Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,		// noWait
		nil)			// args

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",		// consumer
		true,		// autoAck
		false,		// exclusive
		false,		// noLocal
		false,		// noWait
		nil,			// args
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("[RabbitMQ] Received a message: %s", d.Body)
		}
	}()

	fmt.Println("[RabbitMQ] Waiting for messages. To exit press CTRL+C")
	<-forever
}