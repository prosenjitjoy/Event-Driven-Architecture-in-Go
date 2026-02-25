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
	IntegrationEventHandlersKey = "IntegrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"

	InvoicesRepoKey = "invoicesRepo"
	PaymentsRepoKey = "paymentsRepo"
)

// repository table names
const (
	OutboxTableName    = "payments.outbox"
	InboxTableName     = "payments.inbox"
	EventsTableName    = "payments.events"
	SnapshotsTableName = "payments.snapshots"
	SagasTableName     = "payments.sagas"
	InvoicesTableName  = "payments.invoices"
	PaymentsTableName  = "payments.payments"
)
