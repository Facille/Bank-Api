package main

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/therealadik/bank-api/internal/config"
	"github.com/therealadik/bank-api/internal/db"
)

func runMigrations(dsn string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logrus.Fatalf("Ошибка миграций : %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatalf("Ошибка миграций: %v", err)
	}
	logrus.Info("Миграции применены")
}

func main() {
	ctx := context.Background()

	dbCfg := config.LoadDB()

	dsn := db.BuildDSN(dbCfg)
	runMigrations(dsn)

	pool, err := db.New(ctx, dbCfg)

	if err != nil {
		logrus.Fatalf("Ошибка подключения к БД: %v", err)
	}

	defer pool.Close()
}
