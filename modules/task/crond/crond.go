package crond

import (
	"github.com/Cepave/open-falcon-backend/modules/task/crond/jobs"
	"github.com/Cepave/open-falcon-backend/modules/task/g"
	"github.com/bamzi/jobrunner"
)

func Start() {
	jobrunner.Start()
	for _, p := range g.Config().PingUrls {
		jobrunner.Schedule(p.Interval, jobs.PingGetObj{URL: p.URL})
	}
}
