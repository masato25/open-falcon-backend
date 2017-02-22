package newgraph

import "github.com/astaxie/beego"

func ConfigRoutes() {

	newgraph := beego.NewNamespace("/api/v1/graph",
		beego.NSRouter("/keywordSearch", &NewGraphController{}, "get:FilterHostGroup;post:FilterHostGroup"),
	)
	beego.AddNamespace(newgraph)
}
