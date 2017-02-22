package filter

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type GrpHosts struct {
	Id       int    `json:"id";orm:"id"`
	GrpName  string `json:"grp_name";orm:"grp_name"`
	Hostname string `json:"hostname";orm:"hostname"`
}

func HostGroupFilter(filterTxt string) []GrpHosts {
	q := orm.NewOrm()
	q.Using("falcon_portal")
	sqlbuild := fmt.Sprintf(`select g2.id, g2.grp_name, h2.hostname from host h2 INNER JOIN (select g.id, g.grp_name, h.host_id from grp g INNER JOIN grp_host h
	on g.id = h.grp_id
	where g.grp_name regexp '%s') g2 on g2.host_id = h2.id`, filterTxt)
	res := []GrpHosts{}
	q.Raw(sqlbuild).QueryRows(&res)
	return res
}
