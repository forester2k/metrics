package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func DBInit(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("DataBase Initialize: dataSourceName is empty. " +
			"Need flag -d orr DATABASE_DSN environment variable")
	}
	fmt.Println("DSN string ", dsn) //============== del this string
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		str := fmt.Sprintf("DataBase Initialize: can't open DataBase with DSN-string %s", dsn)
		return fmt.Errorf(str+"%w", err)
	}
	DB = db
	return nil
}
