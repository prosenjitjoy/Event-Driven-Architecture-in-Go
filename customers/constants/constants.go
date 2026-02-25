package constants

// grpc service names
const (
	StoreServiceName     = "STORES"
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
	ReplyHandlerKey             = "replyHandlers"

	CustomersRepoKey = "customersRepo"
)

// repository table names
const (
	OutboxTableName    = "customers.outbox"
	InboxTableName     = "customers.inbox"
	EventsTableName    = "customers.events"
	SnapshotsTableName = "customers.snapshots"
	SagasTableName     = "customers.sagas"
	CustomersTableName = "customers.customers"
)
