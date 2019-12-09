package controllers

import (
	"beegoTest/webCrawler/models"
	_ "beegoTest/webCrawler/models"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

type MoveInfoController struct {
	beego.Controller
}

func (this *MoveInfoController) Crawler() {
	o := orm.NewOrm()
	var urlBuff []string
	//1、连接redis数据库
	models.ConnectionRedis("127.0.0.1:6379")

	//2、将爬虫的启动节点，也就是“七月与安生”的url添加到url_queue中
	url1 := "https://movie.douban.com/subject/25827935/"

	models.PutToQueue(url1)
	fmt.Println(models.GetQueueLen())
	//3、每次循环，判断url_queue的长度是否为零，为零则退出爬虫
	//   循环从url队列中取出url，先判断其是否在visit_queue（去重）
	//   如果不在，获取其HTML页面，判断是否是电影简介页面
	//   如果是，解析出电影信息存入数据库，并且解析这个页面的所有url存入url_queue中，将这个url放入visit_queue中

	for {
		//url_queue为空，跳出循环，结束爬虫
		if models.GetQueueLen() == 0 {
			fmt.Println(models.GetQueueLen())
			break
		}

		url := models.PopFromQueue()
		//已经访问过，继续下一次循环
		if models.IsVisit(url) {
			continue
		}

		rsp := httplib.Get(url)
		rsp.Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0")
		rsp.Header("Cookie", `bid=gFP9qSgGTfA; __utma=30149280.1124851270.1482153600.1483055851.1483064193.8; __utmz=30149280.1482971588.4.2.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; ll="118221"; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1483064193%2C%22https%3A%2F%2Fwww.douban.com%2F%22%5D; _pk_id.100001.4cf6=5afcf5e5496eab22.1482413017.7.1483066280.1483057909.; __utma=223695111.1636117731.1482413017.1483055857.1483064193.7; __utmz=223695111.1483055857.6.5.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; _vwo_uuid_v2=BDC2DBEDF8958EC838F9D9394CC5D9A0|2cc6ef7952be8c2d5408cb7c8cce2684; ap=1; viewed="1006073"; gr_user_id=e5c932fc-2af6-4861-8a4f-5d696f34570b; __utmc=30149280; __utmc=223695111; _pk_ses.100001.4cf6=*; __utmb=30149280.0.10.1483064193; __utmb=223695111.0.10.1483064193`)
		//判断是否是电影页面，是则解析所有信息
		htmlstr, err := rsp.String()

		if err != nil {
			beego.Info(err)
		}

		movieName := models.GetMovieName(htmlstr)
		if movieName != "" {
			this.Ctx.WriteString(url)

			movieInfo := models.MovieInfo{}
			movieInfo.Move_name = models.GetMovieName(htmlstr)
			movieInfo.Move_id = models.GetMovieID(url)
			movieInfo.Move_director = models.GetMovieDirector(htmlstr)
			movieInfo.Move_pic = models.GetMoviePic(htmlstr)
			movieInfo.Move_writer = models.GetMovieWrite(htmlstr)
			movieInfo.Move_country = models.GetMovieCountry(htmlstr)
			movieInfo.Move_language = models.GetMovieLanguage(htmlstr)
			movieInfo.Move_main_character = models.GetMovieMainCharacter(htmlstr)
			movieInfo.Move_type = models.GetMovieType(htmlstr)
			movieInfo.Move_on_time = models.GetMovieOnTime(htmlstr)
			movieInfo.Move_span = models.GetMovieSpan(htmlstr)
			movieInfo.Move_grade = models.GetMovieGrade(htmlstr)
			movieInfo.Move_create = models.GetMovieCreate()

			this.Ctx.WriteString(movieInfo.Move_name)
			fmt.Println(movieInfo.Move_name)
			o.Insert(&movieInfo)
		}
		//解析出页面的url，放入url_queue中
		urlBuff = models.GetMoviePageURL(htmlstr)

		for _, u := range urlBuff {
			models.PutToQueue(u)
		}
		//将已经爬过的页面url放入visit_queue中
		fmt.Println(url)
		models.PutToSet(url)

		//循环睡眠1S
		time.Sleep(time.Second)
	}

	fmt.Println("循环结束")
}

//Movie function 获取数据库表中的电影信息
func (this *MoveInfoController) Movie() {
	//1、获取需要得到信息的电影的名字
	movieName := this.GetString("name")

	//2、获取数据库对象
	o := orm.NewOrm()

	//3、新建一个数据库表结构对象，给主键赋值
	movie := models.MovieInfo{}
	movie.Move_name = movieName

	//根据主键已经被赋值的结构和主键的列名，读取指定行的数据
	//如此处，主键是“七月与安生”，主键的列名是“电影名称”
	err := o.Read(&movie, "电影名称")

	if err != nil {
		fmt.Println(err)
	} else {
		this.Ctx.WriteString("导演：" + movie.Move_director +
			"\n" + "电影编号：" + strconv.FormatInt(movie.Move_id, 10) +
			"\n" + "电影名称:" + movie.Move_name +
			"\n" + "电影海报地址：" + movie.Move_pic +
			"\n" + "编剧：" + movie.Move_writer +
			"\n" + "制片国家/地区：" + movie.Move_country +
			"\n" + "语言：" + movie.Move_language +
			"\n" + "主演：" + movie.Move_main_character +
			"\n" + "电影类型：" + movie.Move_type +
			"\n" + "上映时间：" + movie.Move_on_time +
			"\n" + "电影时长:" + movie.Move_span +
			"\n" + "评分：" + movie.Move_span +
			"\n" + "创建时间:" + movie.Move_create)
	}
}

//MovieURL function 显示页面上获取到的url链接
func (this *MoveInfoController) MovieURL() {
	//
	url := "https://movie.douban.com/subject/25827935/"
	rq := httplib.Get(url)

	html, err := rq.String()
	if err != nil {
		panic(err)
	}

	//this.Ctx.WriteString(html)
	urls := models.GetMoviePageURL(html)

	// if len(urls) > 0 {
	// 	this.Ctx.WriteString(urls[0])
	// } else {
	// 	this.Ctx.WriteString("没有采集到")
	// }
	for _, u := range urls {
		this.Ctx.WriteString(u + "\n")
	}
}
