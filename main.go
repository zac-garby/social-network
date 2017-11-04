package main

import (
	"database/sql"
	"log"

	"github.com/Zac-Garby/social-network/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@unix(/Applications/MAMP/tmp/mysql/mysql.sock)/social-network")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	s := server.Server{
		Database: db,
	}

	s.Start(":8080")
}
