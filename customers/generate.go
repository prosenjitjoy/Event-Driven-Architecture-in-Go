package customers

//go:generate buf generate

//go:generate mockgen -destination=customerspb/mock.go -package=customerspb mall/customers/customerspb CustomersServiceClient,CustomersServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/customers/internal/application App

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/customers/internal/domain CustomerRepository
