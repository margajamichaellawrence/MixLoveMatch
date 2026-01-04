package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"

	"mlm/internal/musicapp/lib/rooms"
	"mlm/models" // SQLBoiler generated models
)

// Store handles room queries
type Store struct {
	// No dependencies - store is pure query logic
}

// New creates a new room store
func New() *Store {
	return &Store{}
}

// Rooms returns 0 or more rooms matching the filter
func (s *Store) Rooms(
	ctx context.Context,
	exec boil.ContextExecutor,
	filter rooms.RoomQueryFilter,
) ([]*rooms.Room, error) {
	mods := []qm.QueryMod{}

	// IDs filter
	if len(filter.IDs) > 0 {
		// Convert string IDs to uint64 for current DB schema
		ids := make([]interface{}, len(filter.IDs))
		for i, id := range filter.IDs {
			idNum, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid room ID %s: %w", id, err)
			}
			ids[i] = idNum
		}
		mods = append(mods, qm.WhereIn("id IN ?", ids...))
	}

	// Name filter
	if filter.Name.Valid {
		mods = append(mods, qm.Where("name = ?", filter.Name.String))
	}

	// CreatedBy filter
	if filter.CreatedBy.Valid {
		createdByNum, err := strconv.ParseUint(filter.CreatedBy.String, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid created_by ID %s: %w", filter.CreatedBy.String, err)
		}
		mods = append(mods, qm.Where("created_by = ?", createdByNum))
	}

	// IsActive filter
	if filter.IsActive.Valid {
		mods = append(mods, qm.Where("is_active = ?", filter.IsActive.Bool))
	}

	// CreatedAt filter
	if filter.CreatedAt.Valid {
		mods = append(mods, qm.Where("created_at = ?", filter.CreatedAt.Time))
	}

	// Sorting
	if filter.OrderBy.Valid {
		sortDir := "ASC"
		if filter.Sort.Valid {
			sortDir = filter.Sort.String
		}
		mods = append(mods, qm.OrderBy(filter.OrderBy.String+" "+sortDir))
	}

	// Pagination
	if filter.Limit.Valid {
		mods = append(mods, qm.Limit(filter.Limit.Int))
	}
	if filter.Offset.Valid {
		mods = append(mods, qm.Offset(filter.Offset.Int))
	}

	// Execute query
	dbRooms, err := models.Rooms(mods...).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("query rooms: %w", err)
	}

	return dbRoomsToRooms(dbRooms), nil
}

// Room returns exactly 1 room, errors if 0 or >1 found
func (s *Store) Room(
	ctx context.Context,
	exec boil.ContextExecutor,
	filter rooms.RoomQueryFilter,
) (*rooms.Room, error) {
	results, err := s.Rooms(ctx, exec, filter)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no room found")
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("expected 1 room, got %d", len(results))
	}

	return results[0], nil
}

// Update performs generic update with nullable fields
func (s *Store) Update(
	ctx context.Context,
	exec boil.ContextExecutor,
	update rooms.UpdateRoom,
) error {
	if len(update.IDs) == 0 {
		return fmt.Errorf("no room IDs provided")
	}

	// Build update columns
	cols := make(map[string]interface{})

	if update.Name.Valid {
		cols["name"] = update.Name.String
	}
	if update.IsActive.Valid {
		cols["is_active"] = update.IsActive.Bool
	}
	if update.CreatedAt.Valid {
		cols["created_at"] = update.CreatedAt.Time
	}
	if update.CreatedBy.Valid {
		createdByNum, err := strconv.ParseUint(update.CreatedBy.String, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid created_by ID %s: %w", update.CreatedBy.String, err)
		}
		cols["created_by"] = createdByNum
	}

	if len(cols) == 0 {
		return nil // Nothing to update
	}

	// Convert string IDs to uint64
	ids := make([]interface{}, len(update.IDs))
	for i, id := range update.IDs {
		idNum, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid room ID %s: %w", id, err)
		}
		ids[i] = idNum
	}

	// Execute update
	_, err := models.Rooms(
		qm.WhereIn("id IN ?", ids...),
	).UpdateAll(ctx, exec, cols)

	return err
}

// dbRoomsToRooms converts DB models to domain models
func dbRoomsToRooms(dbRooms []*models.Room) []*rooms.Room {
	result := make([]*rooms.Room, len(dbRooms))
	for i, db := range dbRooms {
		var isActive bool
		if db.IsActive.Valid {
			isActive = db.IsActive.Bool
		}

		var createdAt time.Time
		if db.CreatedAt.Valid {
			createdAt = db.CreatedAt.Time
		}

		result[i] = &rooms.Room{
			ID:        fmt.Sprintf("%d", db.ID),
			Name:      db.Name,
			CreatedBy: fmt.Sprintf("%d", db.CreatedBy),
			IsActive:  isActive,
			CreatedAt: createdAt,
		}
	}
	return result
}

// roomToDBRoom converts domain model to DB model
func roomToDBRoom(room *rooms.Room) (*models.Room, error) {
	id, err := strconv.ParseUint(room.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid room ID: %w", err)
	}

	createdBy, err := strconv.ParseUint(room.CreatedBy, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created_by ID: %w", err)
	}

	return &models.Room{
		ID:        id,
		Name:      room.Name,
		CreatedBy: createdBy,
		IsActive: null.Bool{
			Bool:  room.IsActive,
			Valid: true,
		},
		CreatedAt: null.Time{
			Time:  room.CreatedAt,
			Valid: !room.CreatedAt.IsZero(),
		},
	}, nil
}
