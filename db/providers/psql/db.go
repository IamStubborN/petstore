package psql

import (
	"fmt"
	"strconv"
	"time"

	"github.com/IamStubborN/petstore/config"
	"github.com/IamStubborN/petstore/templates"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	"go.uber.org/zap"
)

type Database struct {
	pool *sqlx.DB
}

func (d Database) InitDatabase(cfg config.DB) *Database {
	return &Database{
		pool: initialSQLConn(cfg),
	}
}

func (d *Database) Close() error {
	if err := d.pool.Close(); err != nil {
		return err
	}

	return nil
}

func initialSQLConn(cfg config.DB) *sqlx.DB {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%s",
		cfg.Host, strconv.Itoa(cfg.Port), cfg.User, cfg.Password, cfg.DB, cfg.SSL, cfg.Timeout)
	db, err := sqlx.Open("postgres", dbInfo)
	if err != nil {
		zap.L().Fatal("can't open connection to database", zap.Error(err))
	}

	retryConnect(db, cfg.Retry)

	migrationLogic(db, cfg.RandomDataCount)

	return db
}

func checkError(f func() error) {
	if err := f(); err != nil {
		zap.L().Error("error in defer", zap.Error(err))
	}
}

func retryConnect(db *sqlx.DB, fatalRetry int) {
	var retryCount int
	for range time.NewTicker(time.Second).C {
		retryCount++
		zap.L().Info("try connect to db",
			zap.String("№", strconv.Itoa(retryCount)))

		if err := db.Ping(); err == nil {
			zap.L().Info("database connected")
			return
		}

		if fatalRetry == retryCount {
			zap.L().Fatal("can't connect to database")
		}
	}
}

func migrationLogic(db *sqlx.DB, randomDataCount int) {
	var isExist bool
	err := db.QueryRow(qm[categoryIsExist]).Scan(&isExist)
	if err != nil {
		zap.L().Fatal("can't exec db migrations", zap.Error(err))
	}
	if !isExist {
		if randomDataCount > 0 {
			if err = templates.GenerateRandomSQLData(randomDataCount); err != nil {
				zap.L().Fatal("can't generate random data for db", zap.Error(err))
			}
		}

		migrations := &migrate.FileMigrationSource{
			Dir: "db/migrations",
		}

		_, err = migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
		if err != nil {
			zap.L().Fatal("can't exec db migrations", zap.Error(err))
		}

		zap.L().Info("migrations complete")
	}
}
