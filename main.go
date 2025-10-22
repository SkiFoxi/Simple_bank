package main

import (
	"database/sql"
	"log"

	"github.com/SkiFoxi/Simple_bank/api"
	"github.com/SkiFoxi/Simple_bank/db/sqlc"
	_ "github.com/lib/pq"
)
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5444/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8081"
)
func main() {
	conn, err := sql.Open(dbDriver, dbSource) //создаем новое соединение с БД
	if err != nil {
		log.Fatal("connot connect to db")
	}	

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Сервер не может быть запущен", err)
	}
}
