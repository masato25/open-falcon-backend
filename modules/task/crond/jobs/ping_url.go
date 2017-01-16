package jobs

import (
	log "github.com/Sirupsen/logrus"
	"github.com/masato25/resty"
)

type PingGetObj struct {
	URL string
}

func (this PingGetObj) Run() {
	this.PingGet()
}

func (this PingGetObj) PingGet() (body string, err error) {
	rt := resty.New()
	resp, err := rt.R().Get(this.URL)
	body = resp.String()
	log.Debugf("body: %v", body)
	if err != nil {
		log.Errorf("ping %s, got error with: %v", this.URL, err.Error())
	}
	return
}
