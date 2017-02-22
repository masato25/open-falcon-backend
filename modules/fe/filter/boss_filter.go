package filter

import (
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/fe/dmodel"
)

func PlatformFilter(dat []dmodel.BossObj, filterTxt string) []dmodel.BossObj {
	res := []dmodel.BossObj{}
	for _, n := range dat {
		if strings.Contains(n.Platform, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func IspFilter(dat []dmodel.BossObj, filterTxt string) []dmodel.BossObj {
	res := []dmodel.BossObj{}
	for _, n := range dat {
		if strings.Contains(n.Isp, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func IdcFilter(dat []dmodel.BossObj, filterTxt string) []dmodel.BossObj {
	res := []dmodel.BossObj{}
	for _, n := range dat {
		if strings.Contains(n.Idc, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func IpFilter(dat []dmodel.BossObj, filterTxt string) []dmodel.BossObj {
	res := []dmodel.BossObj{}
	for _, n := range dat {
		if strings.Contains(n.Ip, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func ProvinceFilter(dat []dmodel.BossObj, filterTxt string) []dmodel.BossObj {
	res := []dmodel.BossObj{}
	for _, n := range dat {
		if strings.Contains(n.Province, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func HostNameFilter(dat []dmodel.BossObj, filterTxt string) []dmodel.BossObj {
	res := []dmodel.BossObj{}
	for _, n := range dat {
		if strings.Contains(n.Hostname, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}
