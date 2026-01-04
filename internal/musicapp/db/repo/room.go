package repo

import (
	"context"
	"fmt"
	"mlm/models"
	"strings"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type RoomRepo struct{}

func NewRoomRepo() *RoomRepo {
	return &RoomRepo{}
}

func (r *RoomRepo) Insert(
	ctx context.Context,
	exec boil.ContextExecutor,
	room *models.Room,
) (*models.Room, error) {
	err := room.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("insert Room: %w", err)
	}

	return room, nil
}

func (r *RoomRepo) BulkInsert(
	ctx context.Context,
	exec boil.ContextExecutor,
	rooms []*models.Room,
) error {
	if len(rooms) == 0 {
		return nil
	}

	placeholders := make([]string, len(rooms))
	args := make([]interface{}, 0, len(rooms)*5)

	for i, room := range rooms {
		placeholders[i] = "(?, ?, ?, ?, ?)"
		args = append(args, room.ID, room.Name, room.CreatedBy, room.IsActive, room.CreatedAt)
	}

	query := fmt.Sprintf(`INSERT INTO rooms (id, name, created_by, is_active, created_at) VALUES %s`, strings.Join(placeholders, ", "))

	_, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("bulk insert rooms: %w", err)
	}

	return nil
}

func (r *RoomRepo) Upsert(
	ctx context.Context,
	exec boil.ContextExecutor,
	room *models.Room,
) (*models.Room, error) {
	err := room.Upsert(
		ctx,
		exec,
		boil.Infer(),
		boil.Infer(),
	)
	if err != nil {
		return nil, fmt.Errorf("upsert room: %w", err)
	}

	return room, nil
}
