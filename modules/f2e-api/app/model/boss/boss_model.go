package boss

import (
	"fmt"
	"strings"

	con "github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

type BossHost struct {
	Platform string `json:"platform" gorm:"column:platform"`
	Province string `json:"province" gorm:"column:province"`
	Isp      string `json:"isp"  gorm:"column:isp"`
	Idc      string `json:"idc" gorm:"column:idc"`
	Ip       string `json:"ip" gorm:"column:ip"`
	Hostname string `json:"hostname" gorm:"column:hostname"`
	RawText  string `json:"-"`
}

func (this BossHost) TableName() string {
	return "hosts"
}

func (this BossHost) MatchString(patt string) bool {
	return strings.Contains(this.RawText, patt)
}

func (this BossHost) MatchStrings(patt []string) bool {
	matched := true
	for _, pa := range patt {
		matched = matched && strings.Contains(this.RawText, pa)
	}
	return matched
}

func GetBossObjs() (res []BossHost) {
	db := con.Con()
	res = []BossHost{}
	db.Boss.Select("platform, province, isp, idc, ip, hostname").Table("hosts").Where("exist = 1 and activate = 1").Scan(&res)
	for i, rs := range res {
		res[i].RawText = fmt.Sprintf("%s %s %s %s %s", rs.Hostname, rs.Isp, rs.Idc, rs.Platform, rs.Province)
	}
	return res
}
