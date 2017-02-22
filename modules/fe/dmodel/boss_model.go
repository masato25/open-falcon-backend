package dmodel

import "github.com/astaxie/beego/orm"

type BossObj struct {
	Platform string `json:"platform"`
	Province string `json:"province"`
	Isp      string `json:"isp"`
	Idc      string `json:"idc"`
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
}

func GetBossObjs() (res []BossObj) {
	res = []BossObj{}
	q := orm.NewOrm()
	q.Using("boss")
	q.Raw("select platform, province, isp, idc, ip, hostname from hosts where exist = 1").QueryRows(&res)
	return res
}
