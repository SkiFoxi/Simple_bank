package api

import (
	db "github.com/SkiFoxi/Simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Сервер который обрабатывает все http запросы для банка
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Создает новый сервер и настраивает все маршруты api
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//Добавляем маршрут
	router.POST("/accounts", server.CreateAccount) //Запись нового аккаунта
	router.GET("/accounts/:id", server.getAccount) //Получение одного аккаунта
	router.GET("/accounts", server.listAccount) //Получение нескольких аккаунтов
	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Запуск сервера
func (server *Server) Start(address string) error {
	return server.router.Run(address) //Функция джина
}
