package routes

var RoutesAgent = map[string]map[string]Handler{
	"GET": {
		"/api/v1/list":    handleList,
	},
	"POST": {
		"/api/v1/valkiria": handleShooter,
	},
}

var RoutesOrchestrator = map[string]map[string]Handler{
	"GET": {
	},
}
