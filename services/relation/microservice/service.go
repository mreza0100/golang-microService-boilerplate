package microservice

import (
	"microServiceBoilerplate/services/relation/domain"
	"microServiceBoilerplate/services/relation/handlers"
	"microServiceBoilerplate/services/relation/instances"
	relationNats "microServiceBoilerplate/services/relation/nats"

	"github.com/mreza0100/golog"
)

func NewRelationService(lgr *golog.Core) instances.Handlers {
	nc := relationNats.Connection(lgr)

	publishers := relationNats.NewPublishers(nc, lgr)

	services := domain.New(&domain.NewOpts{
		Lgr: lgr,
	})

	relationNats.InitSubs(&relationNats.InitSubsOpts{
		Lgr: lgr,
		Srv: services,
		Nc:  nc,
	})

	return handlers.New(&handlers.NewOpts{
		Srv:        services,
		Publishers: publishers,
		Lgr:        lgr,
	})
}
