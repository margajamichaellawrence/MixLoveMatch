package rooms

import (
	"time"

	"github.com/aarondl/null/v8"
)

type Room struct {
	ID        string
	Name      string
	CreatedBy string
	IsActive  bool
	CreatedAt time.Time
}

type RoomQueryFilter struct {
	IDs       []string
	Name      null.String
	CreatedBy null.String
	IsActive  null.Bool
	CreatedAt null.Time

	// Sorting
	OrderBy null.String // "created_at", "name"
	Sort    null.String // "ASC", "DESC"
	Limit   null.Int
	Offset  null.Int
}

type UpdateRoom struct {
	IDs       []string
	Name      null.String
	IsActive  null.Bool
	CreatedAt null.Time
	CreatedBy null.String
}
