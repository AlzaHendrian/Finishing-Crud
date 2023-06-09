package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func ConnectDB() {
	DBurl := "postgres://postgres:Profpakzul26@localhost:5432/Personal-web"

	var err error
	Conn, err = pgx.Connect(context.Background(), DBurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Success to connect database")
}
