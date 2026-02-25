package constants

// grpc service names
const (
	StoresServiceName    = "STORES"
	CustomersServiceName = "CUSTOEMRS"
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

	ShoppingListsRepoKey = "shoppingListsRepo"
	StoresCacheRepoKey   = "storesCacheRepo"
	ProductsCacheRepoKey = "productsCacheRepo"
)

// repository table names
const (
	OutboxTableName    = "depot.outbox"
	InboxTableName     = "depot.inbox"
	EventsTableName    = "depot.events"
	SnapshotsTableName = "depot.snapshots"
	SagasTableName     = "depot.sagas"

	ShoppingListsTableName = "depot.shopping_lists"
	StoresCacheTableName   = "depot.stores_cache"
	ProductsCacheTableName = "depot.products_cache"
)
