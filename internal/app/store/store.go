package store

type Store interface {
	User() UserRepository
	Task() TaskRepository
}
