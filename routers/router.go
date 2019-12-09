package routers

import (
	"beegoTest/webCrawler/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/crawler", &controllers.MoveInfoController{}, "*:Crawler")
	beego.Router("/sqlData", &controllers.MoveInfoController{}, "*:Movie")
	beego.Router("/GetURL", &controllers.MoveInfoController{}, "*:MovieURL")
}
