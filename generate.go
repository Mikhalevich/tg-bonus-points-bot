package generate

//go:generate go tool mockgen -source=internal/domain/port/customer_repository.go -destination=internal/domain/port/customer_repository_mock.go -package=port

//go:generate go tool mockgen -source=internal/domain/port/message_sender.go -destination=internal/domain/port/message_sender_mock.go -package=port
