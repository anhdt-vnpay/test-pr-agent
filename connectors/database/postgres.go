package sql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgresClient(dbConfig map[string]string) *gorm.DB {
 /*
 InitPostgresClient initializes a new PostgreSQL client using the provided database configuration.
 It constructs a Data Source Name (DSN) from the configuration map and attempts to establish a connection.
 If the connection fails or the host is not provided, it returns nil.

 Parameters:
 dbConfig - A map containing the database configuration with keys: "host-v3", "username-v3", "password-v3", "port-v3", "database-v3".

 Returns:
 *gorm.DB - A pointer to the GORM DB instance if the connection is successful, nil otherwise.
 */
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
