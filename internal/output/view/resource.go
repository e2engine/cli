package view

import (
	"time"

	coremodel "github.com/e2engine/core/model"
	idpkg "github.com/e2engine/core/pkg/id"
)

type resource struct {
	ID        string
	Name      string
	Version   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableHeaders returns headers for table-format resource list rendering.
func (r resource) TableHeaders() []string {
	return []string{"ID", "NAME", "VERSION", "CREATED", "UPDATED"}
}

// TableRows constructs rows for table-format resource list rendering.
func (r resource) TableRows() [][]string {
	return [][]string{{
		idpkg.TruncateID(r.ID),
		r.Name,
		r.Version,
		r.CreatedAt.Format(time.RFC3339Nano),
		r.UpdatedAt.Format(time.RFC3339Nano),
	}}
}

func newResourceView[S coremodel.ResourceSpec](r coremodel.Resource[S]) resource {
	return resource{
		ID:        r.ID,
		Name:      r.Name,
		Version:   r.Version,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
