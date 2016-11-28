package routes

var Routes = map[string]map[string]Handler{
	"GET": {
		"/api/v1/chaos": handleChaos,
		"/api/v1/shooter": handleShooter,
		"/api/v1/list": handleList,
	},
}



