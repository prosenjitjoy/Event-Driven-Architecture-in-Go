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

	BasketsRepoKey  = "basketsRepo"
	StoresRepoKey   = "storesRepo"
	ProductsRepoKey = "productsRepo"
)

// repository table names
const (
	OutboxTableName    = "baskets.outbox"
	InboxTableName     = "baskets.inbox"
	EventsTableName    = "baskets.events"
	SnapshotsTableName = "baskets.snapshots"
	SagasTableName     = "baskets.sagas"

	StoresCacheTableName   = "baskets.stores_cache"
	ProductsCacheTableName = "baskets.products_cache"
)

const (
	BasketsStartedCount    = "baskets_started_count"
	BasketsCheckedOutCount = "baskets_checked_out_count"
	BasketsCanceledCount   = "baskets_canceled_count"
)
