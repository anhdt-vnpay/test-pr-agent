package sql

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func InitOracleClient(dbConfig map[string]string) (*sql.DB, error) {
	host := dbConfig["host"]
	username := dbConfig["username"]
	password := dbConfig["password"]
	port := dbConfig["port"]
	dbname := dbConfig["database"]

	var dsn = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`, username, password, host, port, dbname)

	if val, ok := dbConfig["connection-timeout"]; ok {
		var err error
		connectionTimeout, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		dsn = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s?connect_timeout=%d"`, username, password, host, port, dbname, connectionTimeout)
		// dsn = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s" poolSessionTimeout=%s`, username, password, host, port, dbname, val)
	}

	db, err := sql.Open("godror", dsn)
	if err != nil {
		panic(err)
		// return nil, err
	}

	db.SetConnMaxIdleTime(time.Hour)

	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(1000)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Printf("Initialized connection to Oracle DB : %+v\n", db.Stats())

	return db, nil
}
