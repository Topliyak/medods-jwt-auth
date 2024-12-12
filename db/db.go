package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/medods-jwt-auth/config"
)

var (
	pool *pgxpool.Pool
)

func Init() {
	var err error
	var ctx = context.Background()

	if pool, err = CreatePool(ctx); err != nil {
		logrus.Fatal(err.Error())
		panic(err.Error())
	}

	if err := pool.Ping(ctx); err != nil {
		logrus.Fatal(err.Error())
		panic(err.Error())
	}

	logrus.Info("Database pinged successfully")
}

func CreatePool(ctx context.Context) (*pgxpool.Pool, error) {
	dbUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)
	
	return pgxpool.New(ctx, dbUrl)
}

func GetReadOnlyTransaction(ctx context.Context) pgx.Tx {
	return GetTransaction(ctx, pgx.ReadOnly)
}

func GetReadWriteTransaction(ctx context.Context) pgx.Tx {
	return GetTransaction(ctx, pgx.ReadWrite)
}

func GetTransaction(ctx context.Context, accessMode pgx.TxAccessMode) pgx.Tx {
	tx, _ := pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
		AccessMode: accessMode,
		DeferrableMode: pgx.Deferrable,
	})

	return tx
}
