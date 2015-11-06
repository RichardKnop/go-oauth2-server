package api

import "github.com/ant0ine/go-json-rest/rest"

// DevelopmentStack for unit tests
var DevelopmentStack = []rest.Middleware{
	&rest.AccessLogApacheMiddleware{},
	&rest.TimerMiddleware{},
	&rest.RecorderMiddleware{},
	&rest.PoweredByMiddleware{},
	&rest.RecoverMiddleware{
		EnableResponseStackTrace: true,
	},
}

// ProductionStack for production
var ProductionStack = []rest.Middleware{
	&rest.AccessLogApacheMiddleware{
		Format: rest.CombinedLogFormat,
	},
	&rest.TimerMiddleware{},
	&rest.RecorderMiddleware{},
	&rest.PoweredByMiddleware{},
	&rest.RecoverMiddleware{},
	&rest.GzipMiddleware{},
}
