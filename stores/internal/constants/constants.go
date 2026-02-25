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
	AggregateStoreKey           = "aggregateStore"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"

	CatalogHandlersKey = "catalogHandlers"
	MallHandlersKey    = "mallHandlers"

	StoresRepoKey   = "storesRepo"
	ProductsRepoKey = "productsRepo"
	CatalogRepoKey  = "catalogRepo"
	MallRepoKey     = "mallRepo"
)

// repository table names
const (
	OutboxTableName    = "stores.outbox"
	InboxTableName     = "stores.inbox"
	EventsTableName    = "stores.events"
	SnapshotsTableName = "stores.snapshots"
	SagasTableName     = "stores.sagas"
	CatalogTableName   = "stores.products"
	MallTableName      = "stores.stores"
)
