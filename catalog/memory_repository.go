package catalog

import (
	"context"
	"sort"
	"strings"
	"sync"
)

type memoryRepository struct {
	mu       sync.RWMutex
	products map[string]Product
}

// NewInMemoryRepository returns a Repository backed by an in-memory map.
// If seed is nil, a small default catalog is pre-populated.
func NewInMemoryRepository(seed []Product) Repository {
	repo := &memoryRepository{
		products: make(map[string]Product),
	}

	initial := seed
	if len(initial) == 0 {
		initial = defaultCatalog()
	}

	for _, p := range initial {
		repo.products[p.ID] = p
	}

	return repo
}

func (r *memoryRepository) Close() {}

func (r *memoryRepository) PutProduct(ctx context.Context, p Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[p.ID] = p
	return nil
}

func (r *memoryRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.products[id]; ok {
		cp := p
		return &cp, nil
	}
	return nil, ErrorNotFound
}

func (r *memoryRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if take == 0 {
		return []Product{}, nil
	}

	items := r.sortedProducts()
	start := min(int(skip), len(items))
	end := min(start+int(take), len(items))
	copied := make([]Product, end-start)
	copy(copied, items[start:end])
	return copied, nil
}

func (r *memoryRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]Product, 0, len(ids))
	for _, id := range ids {
		if p, ok := r.products[id]; ok {
			results = append(results, p)
		}
	}
	return results, nil
}

func (r *memoryRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if take == 0 {
		return []Product{}, nil
	}

	q := strings.ToLower(query)
	filtered := make([]Product, 0, len(r.products))
	for _, p := range r.products {
		if strings.Contains(strings.ToLower(p.Name), q) || strings.Contains(strings.ToLower(p.Description), q) {
			filtered = append(filtered, p)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Name == filtered[j].Name {
			return filtered[i].ID < filtered[j].ID
		}
		return filtered[i].Name < filtered[j].Name
	})

	start := min(int(skip), len(filtered))
	end := min(start+int(take), len(filtered))
	copied := make([]Product, end-start)
	copy(copied, filtered[start:end])
	return copied, nil
}

func (r *memoryRepository) sortedProducts() []Product {
	items := make([]Product, 0, len(r.products))
	for _, p := range r.products {
		items = append(items, p)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Name == items[j].Name {
			return items[i].ID < items[j].ID
		}
		return items[i].Name < items[j].Name
	})
	return items
}

func defaultCatalog() []Product {
	return []Product{
		{
			ID:          "prod-espresso",
			Name:        "Espresso Beans",
			Description: "Single-origin espresso beans with chocolate and hazelnut notes.",
			Price:       12.5,
		},
		{
			ID:          "prod-aeropress",
			Name:        "Aeropress Brewer",
			Description: "Lightweight coffee maker perfect for brewing on the go.",
			Price:       34.0,
		},
		{
			ID:          "prod-mug",
			Name:        "Stoneware Mug",
			Description: "12oz mug with heat-retaining stoneware and ergonomic handle.",
			Price:       18.0,
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
