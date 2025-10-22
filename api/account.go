package api

import (
	"database/sql"
	"net/http"

	db "github.com/SkiFoxi/Simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"` //oneof в джин обозначает что поле должно иметь одно из этих значений
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { //ShouldBindJSON в джин проверяет выполняются ли условия `json:"currency" binding: "required,oneof=USD EUR"` и  `json:"owner" binding: "required"`
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 400 ошибка браузера
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) //400 ошибка браузера
		return
	}

	ctx.JSON(http.StatusOK, account) //200 ошибка т.е. все прошло хорошо
}

//Получение данных одного аккаунта
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`       //для GET запроса (router.GET("/accounts:id", server.CreateAccount))
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil { //ShouldBindJSON в джин проверяет выполняются ли условия `json:"currency" binding: "required,oneof=USD EUR"` и  `json:"owner" binding: "required"`
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 400 ошибка браузера
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)  
}

//Получение данных нескольких аккаунтов
type listAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`       //для GET запроса (router.GET("/accounts:id", server.CreateAccount))
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil { 
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // 400 ошибка браузера
		return
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}


	ctx.JSON(http.StatusOK, accounts)  
}