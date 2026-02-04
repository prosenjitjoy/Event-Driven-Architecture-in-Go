package es

//go:generate mockgen -destination=mock.go -package=es mall/internal/es AggregateRepository,AggregateStore,EventSourcedAggregate,Aggregate
