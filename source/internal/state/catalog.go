package state

import (
	"fmt"
	"sort"
	"sync"

	"example.edu/6cc545/cw1/internal/model"
)

// Catalog stores file metadata, not the actual file contents.
type Catalog struct {
	mu    sync.RWMutex
	files []model.FileMetadata
}

func NewCatalog() *Catalog {
	return &Catalog{files: []model.FileMetadata{}}
}

// Replace validates the complete candidate before changing the stored data.
func (c *Catalog) Replace(input []model.FileMetadata) error {
	candidate := make([]model.FileMetadata, len(input))
	copy(candidate, input)

	seen := make(map[string]bool, len(candidate))

	for _, file := range candidate {
		if err := file.Validate(); err != nil {
			return fmt.Errorf("invalid metadata for %q: %w", file.Name, err)
		}

		if seen[file.Name] {
			return fmt.Errorf("duplicate file name: %q", file.Name)
		}
		seen[file.Name] = true
	}

	sort.Slice(candidate, func(i, j int) bool {
		return candidate[i].Name < candidate[j].Name
	})

	c.mu.Lock()
	defer c.mu.Unlock()

	c.files = candidate
	return nil
}

// Snapshot returns an independent copy sorted by filename.
func (c *Catalog) Snapshot() []model.FileMetadata {
	c.mu.RLock()
	files := make([]model.FileMetadata, len(c.files))
	copy(files, c.files)
	c.mu.RUnlock()

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	return files
}
