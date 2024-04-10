// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: transaction.sql

package db

import (
	"context"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    account_id,
    date,
    transaction
) VALUES (
  $1, $2, $3
) RETURNING id, account_id, date, transaction
`

type CreateTransactionParams struct {
	AccountID   int64   `json:"account_id"`
	Date        string  `json:"date"`
	Transaction float64 `json:"transaction"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction, arg.AccountID, arg.Date, arg.Transaction)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Date,
		&i.Transaction,
	)
	return i, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, account_id, date, transaction FROM transactions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransaction(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Date,
		&i.Transaction,
	)
	return i, err
}

const listTransactions = `-- name: ListTransactions :many
SELECT id, account_id, date, transaction FROM transactions
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListTransactionsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, listTransactions, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transaction{}
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Date,
			&i.Transaction,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}