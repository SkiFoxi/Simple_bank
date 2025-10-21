package db

import (
	"database/sql"
	"log"
	"testing"
	"os"
	_ "github.com/lib/pq"
)
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5444/simple_bank?sslmode=disable"
)
var testQueries *Queries //Объявляем глобально т.к. будем использовать во всех наших модулях
var testDB *sql.DB
func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource) //создаем новое соединение с БД
	if err != nil {
		log.Fatal("connot connect to db")
	}	
	testQueries = New(testDB) //New для нас создал sqlc в файле db.go чтобы удобно подключаться

	os.Exit(m.Run()) // m.Run это выполнение модульного теста с вывводом результата 
	//os.Exit завершает программу и выводит результат из m.Run()
}