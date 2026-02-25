package search

import (
	"context"
	"database/sql"
	"mall/customers/customerspb"
	"mall/internal/am"
	"mall/internal/di"
	"mall/internal/jetstream"
	pg "mall/internal/postgres"
	"mall/internal/registry"
	"mall/internal/system"
	"mall/internal/tm"
	"mall/ordering/orderingpb"
	"mall/search/internal/application"
	"mall/search/internal/constants"
	"mall/search/internal/domain"
	"mall/search/internal/grpc"
	"mall/search/internal/handlers"
	"mall/search/internal/logging"
	"mall/search/internal/postgres"
	"mall/search/internal/rest"
	"mall/stores/storespb"
)

type Module struct{}

func (*Module) Startup(ctx context.Context, mono system.Service) error {
	return Root(ctx, mono)
}

func Root(ctx context.Context, service system.Service) error {
	container := di.New()

	// setup driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()

		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := customerspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := storespb.Registrations(reg); err != nil {
			return nil, err
		}

		return reg, nil
	})

	stream := jetstream.NewStream(service.Config().Nats.Stream, service.JS(), service.Logger())

	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return service.DB().Begin()
	})

	container.AddScoped(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
		), nil
	})

	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	container.AddScoped(constants.CustomersRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		fallback := grpc.NewCustomerRepository(service.Config().Rpc.Service(constants.CustomersServiceName))

		return postgres.NewCustomerCacheRepository(constants.CustomersCacheTableName, tx, fallback), nil
	})

	container.AddScoped(constants.StoresRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		fallback := grpc.NewStoreRepository(service.Config().Rpc.Service(constants.StoresServiceName))

		return postgres.NewStoreCacheRepository(constants.StoresCacheTableName, tx, fallback), nil
	})

	container.AddScoped(constants.ProductsRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		fallback := grpc.NewProductRepository(service.Config().Rpc.Service(constants.StoresServiceName))

		return postgres.NewProductCacheRepository(constants.ProductsCacheTableName, tx, fallback), nil
	})

	container.AddScoped(constants.OrdersRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)

		return postgres.NewOrderRepository(constants.OrdersTableName, tx), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		orders := c.Get(constants.OrdersRepoKey).(domain.OrderRepository)

		app := application.New(orders)

		return logging.LogApplicationAccess(app, service.Logger()), nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		orders := c.Get(constants.OrdersRepoKey).(domain.OrderRepository)
		customers := c.Get(constants.CustomersRepoKey).(domain.CustomerCacheRepository)
		stores := c.Get(constants.StoresRepoKey).(domain.StoreCacheRepository)
		products := c.Get(constants.ProductsRepoKey).(domain.ProductCacheRepository)

		integrationEventHandlers := handlers.NewIntegrationHandlers(orders, customers, stores, products)
		integrationEventHandlers = logging.LogEventHandlerAccess(
			integrationEventHandlers,
			"IntegrationEvents",
			service.Logger(),
		) // logging wrapper

		inboxStore := c.Get(constants.InboxStoreKey).(tm.InboxStore)
		inboxHandler := tm.InboxHandler(inboxStore)

		return am.NewEventHandler(reg, integrationEventHandlers, inboxHandler), nil
	})

	// setup driver adapters
	if err := grpc.RegisterServerTx(container, service.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, service.Mux(), service.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(service.Mux()); err != nil {
		return err
	}

	if err := handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}

	return nil
}
