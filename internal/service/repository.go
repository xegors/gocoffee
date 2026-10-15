package service

import "github.com/xegors/gocoffee/internal/domain"

// ItemRepository reads and writes the catalog.
type ItemRepository interface {
	Load() ([]domain.Item, error)
	Save(items []domain.Item) error
}
