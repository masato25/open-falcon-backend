package dashboard_screen

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	authapi := r.Group("/api/v1/dashboard")
	authapi.Use(utils.AuthSessionMidd)
	authapi.POST("/screen", ScreenCreate)
	authapi.GET("/screen/:screen_id", ScreenGet)
	authapi.GET("/screens/pid/:pid", ScreenGetsByPid)
	authapi.GET("/screens", ScreenGetsAll)
	authapi.DELETE("/screen/:screen_id", ScreenDelete)
	authapi.PUT("/screen/:screen_id", ScreenUpdate)
}
