package ddd

//go:generate mockgen -destination=mock.go -package=ddd mall/internal/ddd Aggregate,CommandHandler,Entity,EventHandler,EventPublisher,EventSubscriber,ReplyHandler
