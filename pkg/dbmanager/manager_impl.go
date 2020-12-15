package dbmanager

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type managerImpl struct {
	mongo   MongoManager
	mariadb MariaDBManager
}

func (m *managerImpl) Close(ctx context.Context) error {
	var eg errgroup.Group

	eg.Go(func() error {
		return m.mongo.Close(ctx)
	})
	eg.Go(func() error {
		return m.mariadb.Close(ctx)
	})

	return eg.Wait()
}
