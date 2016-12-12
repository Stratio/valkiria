package routes

var RoutesAgent = map[string]map[string]Handler{
	"GET": {
		"/api/v1/chaos":   handleChaos,
		"/api/v1/shooter": handleShooter,
		"/api/v1/list":    handleList,
	},
}

var RoutesOrchestrator = map[string]map[string]Handler{
	"GET": {
		"/api/v1/chaos":   handleChaos,
		"/api/v1/shooter": handleShooter,
		"/api/v1/list":    handleList,
	},
}
