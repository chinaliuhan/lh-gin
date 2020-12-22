package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"lh-gin/tools"
)

/**
qmqp:// 协议名,go语言中定死的
root:root RabbitMQ的账号密码
localhost:5672 ip地址和端口
test  Virtual host名称
*/
//协议名://账号:密码@IP地址:端口号/Virtual host
//const MQURL = "amqp://root:root@localhost:5672/test"

type RabbitMQUtil struct {
	conn      *amqp.Connection //保存的连接
	channel   *amqp.Channel    //通信信息
	QueueName string           //队列名称
	Exchange  string           //交换机
	Key       string           //key
	Mqurl     string           //连接信息
}

//创建结构体实例
func NewRabbitMQUtil(queueName string, exchange string, key string, mqurl string) *RabbitMQUtil {
	//rabbitmq := &RabbitMQUtil{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	rabbitmq := &RabbitMQUtil{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: mqurl}
	//创建连接
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接失败")
	//创建channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

//手动断开channel和connection,如果不断开会被一直占用
func (r *RabbitMQUtil) Destory() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

//错误处理逻辑
func (r *RabbitMQUtil) failOnErr(err error, message string) {
	if err != nil {
		tools.NewLogUtil().Error(message, err)
	}
}

//sample模式step: 1. 创建RabbitMQ实例
func NewRabbitMQSimpleUtil(queueName string) *RabbitMQUtil {

	//实例化,准备发送连接
	//queueName在生产者和消费者中必须是一致的,否则将得不到消息
	//exchange为空: 采用默认交换机,direct模式的交换机此种模式表示直接发送到队列
	//key为空: 这里是没有Key
	//"amqp://root:root@localhost:5672/test"
	mqurl := getMQConfig()

	return NewRabbitMQUtil(queueName, "", "", mqurl)
}

//sample模式step: 2. 生产者
func (r *RabbitMQUtil) ProducerSimple(message string) {
	//1.生成队列(固定用法),如果队列不存在会自动创建,如果存在则跳过创建
	//保证一定能入列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //当最后一个消费者断开连接,是否自动删除
		false, //是否有排他性,如果为true,会创建一个只有自己能访问的队列
		false, //是否阻塞,发送消息以后是否等待服务器的响应
		nil,   //额外属性
	)
	if err != nil {
		tools.NewLogUtil().Error("sample模式-创建队列失败:", err)
	}

	//2.发送消息到队列中
	err = r.channel.Publish(
		r.Exchange,  //交换机名称
		r.QueueName, //队列名称
		false,       //如果为true会根据Exchange类型和routekey规则自动寻找符合要求的队列,找不到就发还给生产者
		false,       //如果为true,发送消息到队列发现队列上没有绑定消费者,则会把消息发还给生产者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		tools.NewLogUtil().Error("sample模式-发送消息失败: %s", err)
	}
	tools.NewLogUtil().Info("sample模式-发送消息成功")
}

//sample模式step: 3. 消费者
func (r *RabbitMQUtil) ConsumeSample() {
	//1.生成队列(固定用法),如果队列不存在会自动创建,如果存在则跳过创建
	//保证一定能入列
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //当最后一个消费者断开连接,是否自动删除
		false, //是否有排他性,如果为true,会创建一个只有自己能访问的队列
		false, //是否阻塞,发送消息以后是否等待服务器的响应
		nil,   //额外属性
	)
	if err != nil {
		tools.NewLogUtil().Error("sample模式-创建队列失败: %s", err)
	}

	//2. 接收消息
	message, err := r.channel.Consume(
		r.QueueName, //队列名称
		"",          //用来区分多个消费者,这里不区分
		true,        //是否自动应答,为true时消费者用完消息是否自动告诉rabbitMQ服务器我们已经用完了,让他删除.为false时需要我们自己写回调了
		false,       //是否具有排他性
		false,       //如果设置为true不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,       //是否阻塞,消费完毕之后下一个再进来,注意:false为阻塞.
		nil,
	)

	if err != nil {
		tools.NewLogUtil().Error("sample模式-接收消息失败: %s", err)
	}

	//3. 消费接收到的消息,
	forever := make(chan int) //利用无buffer的channel造成死循环
	//启动协程处理消息
	go func() {
		for d := range message {
			//这里可以写我们的消息处理逻辑的代码
			tools.NewLogUtil().Info("sample模式-从队列中获取到消息:%s", d.Body)
		}
	}()
	tools.NewLogUtil().Info("[**********]sample模式-等待队列消息中")

	<-forever
}

func getMQConfig() string {
	mqConfig := tools.NewConfigUtil().GetRabbitMQConfig("rabbitmq")
	mqurl := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", mqConfig.UserName, mqConfig.Password, mqConfig.Host, mqConfig.Port, mqConfig.Virtual)
	return mqurl
}

//路由模式step1: 创建实例, 与订阅模式不同的是这里要指定routingkey
func NewRabbitMQRoutingUtil(exchangeName string, routingKey string) *RabbitMQUtil {
	//实例化
	mqurl := getMQConfig()
	rabbitmq := NewRabbitMQUtil("", exchangeName, routingKey, mqurl)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接失败")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "打开channel失败")
	return rabbitmq
}

//路由模式step2: 发送消息
func (r *RabbitMQUtil) ProducerRouting(message string) {

	//创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", //和订阅模式唯一不同的地方是这里改为direct模式
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "创建交换机失败")

	//发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key, //这里必须要设置key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	r.failOnErr(err, "发送消息失败")
	tools.NewLogUtil().Info("消息发送成功")
}

//路由模式step3:接收消息
func (r *RabbitMQUtil) ReceiveRouting() []byte {
	//尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", //和订阅模式不同的地方
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "尝试创建交换机失败")

	//尝试创建队列
	q, err := r.channel.QueueDeclare(
		"", //要系统随机生成
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "创建队列失败")

	//绑定队列
	err = r.channel.QueueBind(
		q.Name,
		r.Key, //设置key
		r.Exchange,
		false,
		nil,
	)

	//消费者流控
	err = r.channel.Qos(
		1,     //当前消费者一次能接受的最大消息数量
		0,     //服务器传递的最大容量（以八位字节为单位）
		false, //如果设置为true 对channel可用
	)
	if err != nil {
		tools.NewLogUtil().Error("消费者限流失败:", err)
	}

	//消费消息
	message, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for d := range message {
		//执行, 进一步处理接收到的消息
		return d.Body
	}
	return nil
}
