package notifications

//go:generate buf generate

//go:generate mockgen -destination=notificationspb/mock.go -package=notificationspb mall/notifications/notificationspb NotificationsServiceClient,NotificationsServiceServer

//go:generate mockgen -destination=internal/application/mock.go -package=application mall/notifications/internal/application App

//go:generate mockgen -destination=internal/domain/mock.go -package=domain mall/notifications/internal/domain CustomerCacheRepository,CustomerRepository
