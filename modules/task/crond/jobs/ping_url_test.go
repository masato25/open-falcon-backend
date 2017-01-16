package jobs

import (
	"testing"

	log "github.com/Sirupsen/logrus"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPingGet(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	p := PingGetObj{URL: "http://localhost:10080/api/v1/alarmadjust/ignoredrecovercond"}
	Convey("Get Endpoint Failed", t, func() {
		body, err := p.PingGet()
		if err != nil {
			log.Debugf("got error: %v", err.Error())
		}
		log.Debug(body)
		So(err, ShouldBeNil)
		So(body, ShouldNotBeEmpty)
	})
}
