package storage_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/xegors/gocoffee/internal/domain"
	"github.com/xegors/gocoffee/internal/storage"
)

func writeFixture(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}

	return path
}

func basicWant() []domain.Item {
	return []domain.Item{
		{
			ID: "ethiopia-chelchele-250", Name: "Ethiopia Chelchele",
			Price: 799, WeightGrams: 250,
			Category: domain.CategoryCoffee, Brew: domain.BrewFilter},
		{
			ID: "basic-temper-57", Name: "Basic Temper 57mm",
			Price:    1500,
			Category: domain.CategoryAccessories},
		{
			ID: "random-filter-250", Name: "Random Filter Coffee",
			Price: 850, WeightGrams: 250,
			Category: domain.CategorySurprise, Brew: domain.BrewFilter},
	}
}

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "items.json")
	repo := storage.NewItemRepository(path)

	if err := repo.Save(basicWant()); err != nil {
		t.Fatalf("cannot save to repo %s: %v", path, err)
	}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("cannot load repo %s: %v", path, err)
	}

	if !reflect.DeepEqual(got, basicWant()) {
		t.Errorf("got %+v, want %+v", got, basicWant())
	}
}

func TestLoadValidFile(t *testing.T) {
	path := writeFixture(t, "items.json", `[
  {"id":"ethiopia-chelchele-250","name":"Ethiopia Chelchele","category":"coffee","price":799,"weight_grams":250,"brew":"filter"},
  {"id":"basic-temper-57","name":"Basic Temper 57mm","category":"accessories","price":1500,"brew":""},
  {"id":"random-filter-250","name":"Random Filter Coffee","category":"surprise","price":850,"weight_grams":250,"brew":"filter"}
]`)
	repo := storage.NewItemRepository(path)

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("cannot load repo %s: %v", path, err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d items, want 3", len(got))
	}

	// Первая позиция: все шесть полей присутствуют в файле.
	if got[0].ID != "ethiopia-chelchele-250" {
		t.Errorf("got[0].ID = %q, want %q", got[0].ID, "ethiopia-chelchele-250")
	}
	if got[0].Name != "Ethiopia Chelchele" {
		t.Errorf("got[0].Name = %q, want %q", got[0].Name, "Ethiopia Chelchele")
	}
	if got[0].Price != 799 {
		t.Errorf("got[0].Price = %d, want %d", got[0].Price, 799)
	}
	if got[0].Category != domain.CategoryCoffee {
		t.Errorf("got[0].Category = %q, want %q", got[0].Category, domain.CategoryCoffee)
	}
	if got[0].WeightGrams != 250 {
		t.Errorf("got[0].WeightGrams = %d, want %d", got[0].WeightGrams, 250)
	}
	if got[0].Brew != domain.BrewFilter {
		t.Errorf("got[0].Brew = %q, want %q", got[0].Brew, domain.BrewFilter)
	}

	// Вторая позиция: weight_grams и brew в файле отсутствуют, ждём нули.
	if got[1].ID != "basic-temper-57" {
		t.Errorf("got[1].ID = %q, want %q", got[1].ID, "basic-temper-57")
	}
	if got[1].Name != "Basic Temper 57mm" {
		t.Errorf("got[1].Name = %q, want %q", got[1].Name, "Basic Temper 57mm")
	}
	if got[1].Price != 1500 {
		t.Errorf("got[1].Price = %d, want %d", got[1].Price, 1500)
	}
	if got[1].Category != domain.CategoryAccessories {
		t.Errorf("got[1].Category = %q, want %q", got[1].Category, domain.CategoryAccessories)
	}
	if got[1].WeightGrams != 0 {
		t.Errorf("got[1].WeightGrams = %d, want 0", got[1].WeightGrams)
	}
	if got[1].Brew != domain.BrewNone {
		t.Errorf("got[1].Brew = %q, want %q", got[1].Brew, domain.BrewNone)
	}

	// Третья позиция: форма та же, отличаются ID и категория.
	if got[2].ID != "random-filter-250" {
		t.Errorf("got[2].ID = %q, want %q", got[2].ID, "random-filter-250")
	}
	if got[2].Category != domain.CategorySurprise {
		t.Errorf("got[2].Category = %q, want %q", got[2].Category, domain.CategorySurprise)
	}
	if got[2].WeightGrams != 250 {
		t.Errorf("got[2].WeightGrams = %d, want %d", got[2].WeightGrams, 250)
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "items.json")
	repo := storage.NewItemRepository(path)

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("got error %v, want %v", err, nil)
	}
	if len(got) != 0 {
		t.Errorf("got len %d, want %d", len(got), 0)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Load created the file: %s", path)
	}
}

func TestLoadEmptyFile(t *testing.T) {
	path := writeFixture(t, "items.json", "")
	repo := storage.NewItemRepository(path)

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("got error %v, want %v", err, nil)
	}
	if len(got) != 0 {
		t.Errorf("got len %d, want %d", len(got), 0)
	}
}

func TestLoadBrokenJSON(t *testing.T) {
	path := writeFixture(t, "items.json", `{"id":`)
	repo := storage.NewItemRepository(path)

	_, err := repo.Load()
	if !errors.Is(err, storage.ErrCorruptedCatalog) {
		t.Fatalf("got error %v, want %v", err, storage.ErrCorruptedCatalog)
	}
	// Мусор в файле и плохая категория — разные дефекты, разные ошибки.
	if errors.Is(err, domain.ErrUnknownCategory) {
		t.Errorf("got %v, want a parsing error instead of %v", err, domain.ErrUnknownCategory)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error %q does not mention the path %q", err, path)
	}
}

func TestLoadUnknownCategory(t *testing.T) {
	path := writeFixture(t, "items.json", `[
  {"id":"ok","category":"coffee","brew":"filter"},
  {"id":"bad","category":"serv","brew":"filter"}
]`)
	repo := storage.NewItemRepository(path)

	got, err := repo.Load()
	if !errors.Is(err, domain.ErrUnknownCategory) {
		t.Fatalf("got error %v, want %v", err, domain.ErrUnknownCategory)
	}
	if got != nil {
		t.Errorf("got %+v, want nil: a failed load returns nothing", got)
	}
	// Индекс нужен, чтобы понять, какую позицию в файле чинить.
	if !strings.Contains(err.Error(), "item 1") {
		t.Errorf("error %q does not point at item 1", err)
	}
}

func TestLoadUnknownBrew(t *testing.T) {
	path := writeFixture(t, "items.json", `[
  {"id":"ok","category":"coffee","brew":"filter"},
  {"id":"bad","category":"coffee","brew":"turkish"}
]`)
	repo := storage.NewItemRepository(path)

	_, err := repo.Load()
	if !errors.Is(err, domain.ErrUnknownBrew) {
		t.Fatalf("got error %v, want %v", err, domain.ErrUnknownBrew)
	}
}

func TestLoadDirectoryPath(t *testing.T) {
	// Путь указывает на каталог: прочитать нельзя, но и JSON он не содержит.
	dir := t.TempDir()
	repo := storage.NewItemRepository(dir)

	_, err := repo.Load()
	if err == nil {
		t.Fatalf("got nil error, want a read error")
	}
	if errors.Is(err, storage.ErrCorruptedCatalog) {
		t.Errorf("got %v, want a read error instead of %v", err, storage.ErrCorruptedCatalog)
	}
	if !strings.Contains(err.Error(), dir) {
		t.Errorf("error %q does not mention the path %q", err, dir)
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data", "items.json")
	repo := storage.NewItemRepository(path)

	if _, err := os.Stat(filepath.Dir(path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("setup: %s must not exist yet", filepath.Dir(path))
	}

	if err := repo.Save(basicWant()); err != nil {
		t.Fatalf("save: %v", err)
	}

	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("stat %s: %v", filepath.Dir(path), err)
	}
	if !info.IsDir() {
		t.Errorf("%s is not a directory", filepath.Dir(path))
	}
}

func TestSaveWritesExpectedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "items.json")
	repo := storage.NewItemRepository(path)

	if err := repo.Save([]domain.Item{{
		ID:          "ethiopia",
		Name:        "Ethiopia",
		Category:    domain.CategoryCoffee,
		Price:       799,
		WeightGrams: 250,
		Brew:        domain.BrewFilter,
	}}); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}

	// Порядок ключей повторяет порядок полей в itemDTO.
	wantJSON := `[
  {
    "id": "ethiopia",
    "name": "Ethiopia",
    "category": "coffee",
    "price": 799,
    "weight_grams": 250,
    "brew": "filter"
  }
]`
	if string(got) != wantJSON {
		t.Errorf("got:\n%s\nwant:\n%s", got, wantJSON)
	}
}

func TestSaveEmptyCatalogWritesEmptyArray(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "items.json")
	repo := storage.NewItemRepository(path)

	for _, items := range [][]domain.Item{nil, {}} {
		if err := repo.Save(items); err != nil {
			t.Fatalf("save %v: %v", items, err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read back: %v", err)
		}
		// null читается обратно без ошибки, но файл каталога должен читаться.
		if string(got) != "[]" {
			t.Errorf("Save(%v) wrote %q, want %q", items, got, "[]")
		}
	}
}
