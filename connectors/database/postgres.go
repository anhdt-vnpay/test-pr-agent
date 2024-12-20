package sql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgresClient(dbConfig map[string]string) *gorm.DB {
	host := dbConfig["host-v3"]
	username := dbConfig["username-v3"]
	password := dbConfig["password-v3"]
	port := dbConfig["port-v3"]
	dbname := dbConfig["database-v3"]

	if host == "" {
		return nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, username, password, dbname, port)
	// fmt.Printf("PgService.NewPgService: dsn = %s\n", dsn)
	postgresDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage

	}), &gorm.Config{
		// PrepareStmt: true,
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Printf("Failed to connect to Postgres server : %v\n", err.Error())
		return nil
	}

	return postgresDB

}
