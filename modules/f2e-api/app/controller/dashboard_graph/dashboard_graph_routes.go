package dashboard_graph

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed
const TMP_GRAPH_FILED_DELIMITER = "|"

func Routes(r *gin.Engine) {
	db = config.Con()
	authapi := r.Group("/api/v1/dashboard")
	authapi.Use(utils.AuthSessionMidd)
	authapi.POST("/tmpgraph", DashboardTmpGraphCreate)
	authapi.GET("/tmpgraph/:id", DashboardTmpGraphQuery)
	authapi.POST("/graph", DashboardGraphCreate)
	authapi.PUT("/graph/:id", DashboardGraphUpdate)
	authapi.GET("/graph/:id", DashboardGraphGet)
	authapi.DELETE("/graph/:id", DashboardGraphDelete)
	authapi.GET("/graphs/screen/:screen_id", DashboardGraphGetsByScreenID)
}
