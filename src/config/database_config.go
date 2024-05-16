package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	GrpcDB *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func NewGrpcDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		GrpcDB: GrpcDB(envConfig),
	}
	return databaseConfig
}

func GrpcDB(envConfig *EnvConfig) *PostgresDatabase {
	var url string
	if envConfig.GrpcDB.Password == "" {
		url = fmt.Sprintf(
			"postgresql://%s@%s:%s/%s",
			envConfig.GrpcDB.User,
			envConfig.GrpcDB.Host,
			envConfig.GrpcDB.Port,
			envConfig.GrpcDB.Database,
		)
	} else {
		url = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			envConfig.GrpcDB.User,
			envConfig.GrpcDB.Password,
			envConfig.GrpcDB.Host,
			envConfig.GrpcDB.Port,
			envConfig.GrpcDB.Database,
		)
	}

	connection, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	connection.SetConnMaxLifetime(300 * time.Second)
	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(10)

	GrpcDB := &PostgresDatabase{
		Connection: connection,
	}
	return GrpcDB
}
