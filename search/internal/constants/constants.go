package constants

// grpc service names
const (
	StoresServiceName    = "STORES"
	CustomersServiceName = "CUSTOMERS"
)

// dependency injection keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTransactionKey      = "tx"
	MessagePublisherKey         = "messagePublisher"
	MessageSubscriberKey        = "messageSubscriber"
	EventPublisherKey           = "eventPublisher"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"

	OrdersRepoKey    = "ordersRepo"
	CustomersRepoKey = "customersRepo"
	StoresRepoKey    = "storesRepo"
	ProductsRepoKey  = "productsRepo"
)

// repository table names
const (
	OutboxTableName         = "search.outbox"
	InboxTableName          = "search.inbox"
	EventsTableName         = "search.events"
	SnapshotsTableName      = "search.snapshots"
	SagasTableName          = "search.sagas"
	OrdersTableName         = "search.orders"
	CustomersCacheTableName = "search.customers_cache"
	StoresCacheTableName    = "search.stores_cache"
	ProductsCacheTableName  = "search.products_cache"
)
