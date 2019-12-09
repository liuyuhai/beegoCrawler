package models

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type MovieInfo struct {
	ID                  int32  `orm:"auto;column(序号);size(10)"`
	Move_id             int64  `orm:"column(电影编号);size(11)"`
	Move_name           string `orm:"unique;column(电影名称);size(100)"`
	Move_pic            string `orm:"column(电影图片);size(200)"`
	Move_director       string `orm:"column(导演);size(50)"`
	Move_writer         string `orm:"column(编剧);size(50)"`
	Move_country        string `orm:"column(电影产地);size(50)"`
	Move_language       string `orm:"column(语言);size(50)"`
	Move_main_character string `orm:"column(主演);size(300)"`
	Move_type           string `orm:"column(电影类型);size(300)"`
	Move_on_time        string `orm:"column(上映时间);size(50)"`
	Move_span           string `orm:"column(时长);size(50)"`
	Move_grade          string `orm:"column(评分);size(50)"`
	Move_create         string `orm:"column(创建时间);size(50)"`
}

func init() {
	//连接数据库
	orm.RegisterDataBase("default", "mysql", "root:LiuYuHai_235@/Move_Info?charset=utf8")

	//注册model
	orm.RegisterModel(new(MovieInfo))

	//建立数据库表
	orm.RunSyncdb("default", false, true)
}

//GetMovieDirector function 获取导演的名字
func GetMovieDirector(movieHtml string) string {
	//写一个导演名字的正则
	//<a href="/celebrity/1274534/" rel="v:directedBy">曾国祥</a>
	//.*?表示非贪婪匹配，即匹配到第一个就返回；.*是贪婪匹配，将匹配所有满足条件的部分

	reg := regexp.MustCompile(`<a.*?rel="v:directedBy">(.*?)</a>`)

	//FindAllStringSubmatch函数返回二元切片，一元表示所有匹配的字符串，
	//二元第零个元素是每个匹配的字符串，其余是需要提取的字符串
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

//GetMovieID function 读取电影编号，从url中获取
func GetMovieID(movieurl string) int64 {
	//"https://movie.douban.com/subject/25827935/"
	reg := regexp.MustCompile(`https://movie.douban.com/subject/([0-9]+)/.*`)

	result := reg.FindAllStringSubmatch(movieurl, -1)

	if len(result) == 0 {
		return 0
	}

	movieID, err := strconv.ParseInt(result[0][1], 10, 64)

	if err != nil {
		return 0
	}

	return movieID
}

//GetMovieName function 函数获取电影的名字
func GetMovieName(movieHtml string) string {
	//<span property="v:itemreviewed">七月与安生</span>
	re := regexp.MustCompile(`<span.*?property="v:itemreviewed">(.*?)</span>`)

	result := re.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

//GetMoviePic function 获取电影海报地址
func GetMoviePic(movieHtml string) string {
	//str := `<img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p2378140502.jpg" title="点击看更多海报" alt="七月与安生" rel="v:image">`
	re := regexp.MustCompile(`<img\s*src="(.*?)"\s*title="点击看更多海报"\s*alt=".*"\s*rel="v:image"\s*/>`)

	result := re.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

//GetMovieWrite function 获取电影编剧的名字
func GetMovieWrite(movieHtml string) string {
	//<a href="/celebrity/1359881/">林咏琛</a>
	reg := regexp.MustCompile(`<a\s*href="/[a-zA-Z]+/[0-9]+/">(.*?)</a>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	var write string

	for i := 0; i < len(result); i++ {
		for j := 1; j < len(result[i]); j++ {
			write = write + result[i][j] + ";"
		}
	}
	return write
}

//GetMovieCountry function获取电影产地信息
func GetMovieCountry(movieHtml string) string {
	//<span class="pl">制片国家/地区:</span>
	// 中国大陆 / 中国香港<br>

	//有些标签网页上省略了/，但正则要写出来？比如image、br
	reg := regexp.MustCompile(`<span.*?>制片国家/地区:</span>(.*?)\s*<br/>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

//GetMovieLanguage function获取电影语言
func GetMovieLanguage(movieHtml string) string {
	//<span class="pl">语言:</span>
	//" 汉语普通话"
	//</br>
	reg := regexp.MustCompile(`span.*>语言:</span>(.*?)\s*<br/>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

//GetMovieMainCharacter function获取电影主演
func GetMovieMainCharacter(movieHtml string) string {
	//<a href="/celebrity/1274224/" rel="v:starring">周冬雨</a>
	reg := regexp.MustCompile(`<a\s*href="/[a-zA-Z]+/[0-9]+/"\s*rel="v:starring">(.*?)</a>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)

	if len(result) == 0 {
		return ""
	}

	var mainCharacter string

	for i := 0; i < len(result); i++ {
		for j := 1; j < len(result[i]); j++ {
			mainCharacter = mainCharacter + result[i][j] + ";"
		}
	}

	return mainCharacter
}

//GetMovieType function获取电影类型
func GetMovieType(movieHTML string) string {
	//<span property="v:genre">剧情</span>
	reg := regexp.MustCompile(`<span\s*property="v:genre">(.*?)</span>`)

	result := reg.FindAllStringSubmatch(movieHTML, -1)

	if len(result) == 0 {
		return ""
	}

	var movieType string
	for i := 0; i < len(result); i++ {
		for j := 1; j < len(result[i]); j++ {
			movieType = movieType + result[i][j] + ";"
		}
	}

	return movieType
}

//GetMovieOnTime function获取电影的上映时间
func GetMovieOnTime(movieHTML string) string {
	//<span property="v:initialReleaseDate" content="2016-09-14(中国大陆)">2016-09-14(中国大陆)</span>
	//<span property="v:initialReleaseDate" content="2016-10-27(中国香港)">2016-10-27(中国香港)</span>
	reg := regexp.MustCompile(`<span\s*property="v:initialReleaseDate"\s*content="(.*?)">`)

	result := reg.FindAllStringSubmatch(movieHTML, -1)

	if len(result) == 0 {
		return ""
	}

	var movieOnTime string

	for i := 0; i < len(result); i++ {
		for j := 1; j < len(result[i]); j++ {
			movieOnTime = movieOnTime + result[i][j] + "/"
		}
	}

	return movieOnTime
}

//GetMovieSpan function获取电影时长
func GetMovieSpan(movieHTML string) string {
	//<span property="v:runtime" content="109">109分钟</span>
	reg := regexp.MustCompile(`<span\s*property="v:runtime".*>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHTML, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

//GetMovieGrade function获取电影评分
func GetMovieGrade(movieHTML string) string {
	//<strong class="ll rating_num" property="v:average">7.6</strong>
	reg := regexp.MustCompile(`<strong.*property="v:average">(.*?)</strong>`)

	result := reg.FindAllStringSubmatch(movieHTML, -1)

	if len(result) == 0 {
		return ""
	}

	return result[0][1]
}

//GetMovieCreate function获取当前条目创建的时间
func GetMovieCreate() string {
	ti := time.Now()

	createTime := fmt.Sprintf("%d年%d月%d日%d时%d分", ti.Year(), ti.Month(), ti.Day(), ti.Hour(), ti.Minute())

	return createTime
}

//GetMoviePageURL function获取页面上其他电影的url
func GetMoviePageURL(movieHTML string) []string {
	// <a href="https://movie.douban.com/subject/4739952/?from=subject-page" class="" >初恋这件小事</a>

	reg := regexp.MustCompile(`<a\s*?href="(https://movie.douban.com/subject/[0-9.]+/.*?from.*?)"\s*?>`)

	result := reg.FindAllStringSubmatch(movieHTML, -1)
	var urls []string

	//fmt.Println(len(result))
	if len(result) == 0 {

		return urls
	}

	for _, v := range result {

		urls = append(urls, v[1])
	}

	return urls
}
