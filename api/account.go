package api

import (
	"database/sql"
	"net/http"

	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required" example:"Paty Jiménez"`
	// We wait for a card number with a len of 16 digits
	CardNumber string `json:"card_number" binding:"required,len=16" example:"2011232967685539"`
}

// CreateTags		godoc
// @Summary			CreateAccount creates an stori account
// @Description		An account for a customer is created including the card number.
// @Accept			json
// @Produce			application/json
// @Param   account         body     createAccountRequest        true  "Account request"
// @Tags			Account
// @Router			/accounts [post]
// @Success 200	{object} db.Account "Account structure"
func (server *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:      req.Owner,
		CardNumber: req.CardNumber,
	}

	account, err := server.Store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// CreateTags		godoc
// @Summary			GetAccount gets an account by ID
// @Description		The info of an account is returned.
// @Produce			application/json
// @Tags			Account
// @Param   id         path     int        true  "The id of the account to be searched"          minimum(1)
// @Router			/accounts/{id} [get]
// @Success 200	{object} db.Account "Account structure"
func (server *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.Store.GetAccount(ctx, req.ID)
	if err != nil {
		// No rows found with the specific ID
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// CreateTags		godoc
// @Summary			ListAccounts gets a list of accounts
// @Description		A list of accounts and their info is returned by pagination.
// @Produce			application/json
// @Tags			Account
// @Param   page_id         query     int        true  "The id of the page where to start"          minimum(1)
// @Param   page_size         query     int        true  "The number of registers to show per page"          minimum(1) maximum(10)
// @Router			/accounts [get]
// @Success 200	{object} []db.Account "Account list structure"
func (server *Server) ListAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.Store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
