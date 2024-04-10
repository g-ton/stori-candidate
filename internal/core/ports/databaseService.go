package ports

import "github.com/g-ton/stori-candidate/internal/core/entities"

// As services

type DatabaseService interface {
	CreateTransaction(*entities.Transaction) error
}
