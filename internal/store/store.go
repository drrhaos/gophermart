package store

import "context"

type StorageInterface interface {
	UserRegister(ctx context.Context, user string, password string) bool
	Ping(ctx context.Context) bool
}

type StorageContext struct {
	storage StorageInterface
}

func (sc *StorageContext) SetStorage(storage StorageInterface) {
	sc.storage = storage
}

func (sc *StorageContext) UserRegister(ctx context.Context, user string, password string) bool {
	return sc.storage.UserRegister(ctx, user, password)
}

func (sc *StorageContext) Ping(ctx context.Context) (exists bool) {
	return sc.storage.Ping(ctx)
}
