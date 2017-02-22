package newgraph

import (
	"github.com/Cepave/open-falcon-backend/modules/fe/dmodel"
	"github.com/Cepave/open-falcon-backend/modules/fe/filter"
	"github.com/Cepave/open-falcon-backend/modules/fe/http/base"
	"github.com/emirpasic/gods/lists/arraylist"
)

type NewGraphController struct {
	base.BaseController
}

func (this *NewGraphController) FilterHostGroup() {
	AccpectTypes := arraylist.New()
	AccpectTypes.Add("platfrom", "idc", "isp", "province", "hostname", "hostgroup")
	baseResp := this.BasicRespGen()
	q := this.GetString("q", "--")
	ftype := this.GetString("filter_type", "all")
	if q == "--" {
		this.ResposeError(baseResp, "q is empty, please check it!")
	}
	if !(ftype == "all" || AccpectTypes.Contains(ftype)) {
		this.ResposeError(baseResp, "filter_type got error type, please check it!")
	}
	bossList := dmodel.GetBossObjs()
	res := map[string]interface{}{
		"platfrom":  filter.PlatformFilter(bossList, q),
		"idc":       filter.IdcFilter(bossList, q),
		"isp":       filter.IspFilter(bossList, q),
		"province":  filter.ProvinceFilter(bossList, q),
		"hostname":  filter.HostNameFilter(bossList, q),
		"hostgroup": filter.HostGroupFilter(q),
	}
	baseResp.Data["res"] = res
	this.ServeApiJson(baseResp)
	return
}
