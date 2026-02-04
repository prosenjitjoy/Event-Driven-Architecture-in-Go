package payments

//go:generate buf generate

//go:generate mockgen -destination=paymentspb/mock.go -package=paymentspb mall/payments/paymentspb PaymentsServiceClient,PaymentsServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/payments/internal/application App

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/payments/internal/domain InvoiceRepository,PaymentRepository
