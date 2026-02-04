package stores

//go:generate buf generate

//go:generate mockgen -destination=storespb/mock.go -package=storespb mall/stores/storespb StoresServiceClient,StoresServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/stores/internal/application App

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/stores/internal/domain CatalogRepository,MallRepository,ProductRepository,StoreRepository
