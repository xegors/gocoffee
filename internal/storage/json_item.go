package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xegors/gocoffee/internal/domain"
	"github.com/xegors/gocoffee/internal/service"
)

type itemDTO struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Category    string       `json:"category"`
	Price       domain.Money `json:"price"`
	WeightGrams int          `json:"weight_grams"`
}

// ItemRepository keeps the catalog in a JSON file.
//
// The zero value is not usable: build one with NewItemRepository.
type ItemRepository struct {
	path string
}

var _ service.ItemRepository = (*ItemRepository)(nil)

// NewItemRepository returns a repository that keeps the catalog in the file at
// path.
//
// The path is not touched here: a missing or empty file is a valid empty
// catalog, and the parent directory is created by the first Save.
func NewItemRepository(path string) *ItemRepository {
	return &ItemRepository{path: path}
}

func (d itemDTO) toDomain() (domain.Item, error) {
	cat, err := domain.ParseCategory(d.Category)
	if err != nil {
		return domain.Item{}, err
	}

	return domain.Item{
		ID:          d.ID,
		Name:        d.Name,
		Category:    cat,
		Price:       d.Price,
		WeightGrams: d.WeightGrams,
	}, nil
}

func fromDomain(i domain.Item) itemDTO {
	return itemDTO{
		ID:          i.ID,
		Name:        i.Name,
		Category:    string(i.Category),
		Price:       i.Price,
		WeightGrams: i.WeightGrams,
	}
}

// Load returns the catalog stored in the file.
//
// A missing file, an empty file and a file holding an empty JSON array all
// yield an empty catalog and a nil error. Any read error other than a missing
// file is returned as is, wrapped with the path; a file that cannot be parsed
// is reported as ErrCorruptedCatalog, and an unusable field in a parsed file
// fails the whole load with the index of the offending item.
func (r *ItemRepository) Load() ([]domain.Item, error) {
	data, err := os.ReadFile(r.path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("read catalog %s: %w", r.path, err)
	case len(data) == 0:
		return nil, nil
	}

	var slice []itemDTO
	if err := json.Unmarshal(data, &slice); err != nil {
		return nil, fmt.Errorf("%w: %s: %w", ErrCorruptedCatalog, r.path, err)
	}

	items := make([]domain.Item, 0, len(slice))
	for i, d := range slice {
		item, err := d.toDomain()
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		items = append(items, item)
	}

	return items, nil
}

// Save writes items to the file, creating the parent directory when needed.
//
// The file is rewritten in full, so items left out of the slice disappear from
// the catalog. An empty slice is stored as an empty JSON array rather than
// null.
func (r *ItemRepository) Save(items []domain.Item) error {
	slice := make([]itemDTO, 0, len(items))
	for _, item := range items {
		slice = append(slice, fromDomain(item))
	}

	if len(slice) == 0 {
		slice = []itemDTO{}
	}

	data, err := json.MarshalIndent(slice, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal catalog: %w", err)
	}

	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create catalog dir: %w", err)
	}

	if err := os.WriteFile(r.path, data, 0o644); err != nil {
		return fmt.Errorf("write catalog: %w", err)
	}

	return nil
}
