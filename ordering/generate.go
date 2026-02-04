package ordering

//go:generate buf generate

//go:generate mockgen -destination=orderingpb/mock.go -package=orderingpb mall/ordering/orderingpb OrderingServiceClient,OrderingServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/ordering/internal/application App

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/ordering/internal/domain CustomerRepository,InvoiceRepository,NotificationRepository,OrderRepository,PaymentRepository,ShoppingRepository
