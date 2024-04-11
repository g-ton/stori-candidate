package api

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/g-ton/stori-candidate/api/helper"
	db "github.com/g-ton/stori-candidate/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransactionRequest struct {
	AccountID   int64   `json:"account_id" binding:"required,min=1"`
	Date        string  `json:"date" binding:"required"`
	Transaction float64 `json:"transaction" binding:"required"`
}

// -- CREATE
func (server *Server) CreateTransaction(ctx *gin.Context) {
	var req createTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We validate that the account exists to proceed
	if !server.validAccount(ctx, req.AccountID) {
		return
	}

	arg := db.CreateTransactionParams{
		AccountID:   req.AccountID,
		Date:        req.Date,
		Transaction: req.Transaction,
	}

	transaction, err := server.Store.CreateTransaction(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

type getTransactionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// -- GET
func (server *Server) GetTransaction(ctx *gin.Context) {
	var req getTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transaction, err := server.Store.GetTransaction(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

type listTransactionRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// -- LIST
func (server *Server) ListTransactions(ctx *gin.Context) {
	var req listTransactionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListTransactionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	transactions, err := server.Store.ListTransactions(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64) bool {
	_, err := server.Store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	return true
}

type listTransactionsByAccountRequest struct {
	AccountID int64  `json:"account_id" binding:"required,min=1"`
	Mails     string `json:"mails" binding:"required"`
}

func (server *Server) GetSummaryInfoByDB(ctx *gin.Context) {
	var req listTransactionsByAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transactions, err := server.Store.ListTransactionsByAccount(ctx, req.AccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	r := helper.GetSummaryInfo(transactions)

	// Receiver email addresses.
	// The list of mails is split by a comma to separate one mail from other
	mailsTo := strings.Split(req.Mails, ",")
	err = helper.ProcessTemplateEmailForTransaction(r, mailsTo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

type listTransactionsByFileRequest struct {
	FilePath string `json:"file_path" binding:"required"`
	Mails    string `json:"mails" binding:"required"`
}

func (server *Server) GetSummaryInfoByFile(ctx *gin.Context) {
	var req listTransactionsByFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transactions, err := helper.ProcessFile(req.FilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	r := helper.GetSummaryInfo(transactions)

	// Receiver email addresses.
	// The list of mails is split by a comma to separate one mail from other
	mailsTo := strings.Split(req.Mails, ",")
	err = helper.ProcessTemplateEmailForTransaction(r, mailsTo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}
