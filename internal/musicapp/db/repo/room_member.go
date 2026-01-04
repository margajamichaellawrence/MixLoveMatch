package repo

import (
	"context"
	"fmt"
	"mlm/models"
	"strings"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type RoomMemberRepo struct{}

func NewRoomMember() *RoomMemberRepo {
	return &RoomMemberRepo{}
}

func (r *RoomMemberRepo) Insert(
	ctx context.Context,
	exec boil.ContextExecutor,
	roomMember *models.RoomMember) (*models.RoomMember, error) {
	err := roomMember.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("insert room member: %w", err)
	}

	return roomMember, nil
}

func (r *RoomMemberRepo) BulkInsert(
	ctx context.Context,
	exec boil.ContextExecutor,
	roomMembers []*models.RoomMember) error {
	if len(roomMembers) == 0 {
		return nil
	}

	placeholders := make([]string, len(roomMembers))
	args := make([]interface{}, 0, len(roomMembers)*5)

	for i, roomMember := range roomMembers {
		placeholders[i] = "(?, ?, ?, ?, ?, ?)"
		args = append(args, roomMember.ID, roomMember.RoomID, roomMember.UserID, roomMember.JoinedAt, roomMember.LeftAt)
	}

	query := fmt.Sprintf(`INSERT INTO room_member (id, room_id, user_id, joined_at, left_at) VALUES %s`, strings.Join(placeholders, ","))

	_, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert room member: %w", err)
	}

	return nil
}

func (r *RoomMemberRepo) Upsert(
	ctx context.Context,
	exec boil.ContextExecutor,
	roomMember *models.RoomMember,
) (*models.RoomMember, error) {
	err := roomMember.Upsert(
		ctx,
		exec,
		boil.Infer(),
		boil.Infer(),
	)
	if err != nil {
		return nil, fmt.Errorf("upsert room member: %w", err)
	}

	return roomMember, nil
}
