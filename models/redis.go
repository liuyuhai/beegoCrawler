package models

/*
	用redis来实现对于url的生产者和消费者的模型
	步骤：
	1、解析出当前页面上所有的超链接，将其加入redis的一个队列中url_queue
	2、从url_queue中取出url，解析出想要的内容，并将url存入visit_queue中
	3、每次从url_queue中去除的url，都需要判断是否在visit_queue中，如果是，则这个url已经访问过，不用再访问
	4、一旦url_queue队列长度为0，则当前页面的所有子页面全部被访问，爬取结束
*/
import (
	"github.com/astaxie/goredis"
)

const (
	URL_QUEUE = "url_queue"
	VISIT_SET = "visit_set"
)

var (
	redisClient goredis.Client
)

//ConnectionRedis function连接redis服务器
func ConnectionRedis(addr string) {
	redisClient.Addr = addr
}

//PutToQueue function将url添加到队列中
func PutToQueue(url string) {
	redisClient.Lpush(URL_QUEUE, []byte(url))
}

//PopFromQueue function从url_queue队列中取出url
func PopFromQueue() string {
	url, err := redisClient.Lpop(URL_QUEUE)

	if err != nil {
		panic(err)
	}
	return string(url)
}

//GetQueueLen function获取url_queue队列的长度
func GetQueueLen() int {
	len, err := redisClient.Llen(URL_QUEUE)
	if err != nil {
		panic(err)
	}
	return len
}

//PushToSet function将每一个操作过的url加入visit_set集合中
func PutToSet(url string) {
	// redisClient.Lpush(VISIT_SET, []byte(url))
	redisClient.Sadd(VISIT_SET, []byte(url)) //向名称为VISIT_SET的集合中添加元素
}

//IsVisit function判断url是否存在visit_set中
func IsVisit(url string) bool {
	re, err := redisClient.Sismember(VISIT_SET, []byte(url))
	if err != nil {
		return false
	}

	return re
}
