package store

import (
	"context"
	"fmt"
	"mlm/internal/musicapp/lib/room_members"
	"mlm/models"
	"strconv"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type Store struct{}

func New() *Store {
	return &Store{}
}

func (s *Store) RoomMembers(
	ctx context.Context,
	exec boil.ContextExecutor,
	filter room_members.RoomMemberQueryFilter,
) ([]*room_members.RoomMembers, error) {
	mods := []qm.QueryMod{}

	if len(filter.IDs) > 0 {
		ids := make([]interface{}, len(filter.IDs))
		for i, id := range filter.IDs {
			idNum, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid room member ID %s: %w", id, err)
			}
			ids[i] = idNum
		}
		mods = append(mods, qm.WhereIn("id IN (?)", ids))
	}

	if filter.RoomID.Valid {
		mods = append(mods, qm.Where("room_id = ?", filter.RoomID.String))
	}

	if filter.UserID.Valid {
		mods = append(mods, qm.Where("user_id = ?", filter.UserID.String))
	}

	if filter.JoinedAt.Valid {
		mods = append(mods, qm.WhereIn("joined_at = ?", filter.JoinedAt.Time))
	}

	if filter.LeftAt.Valid {
		mods = append(mods, qm.WhereIn("left_at = ?", filter.LeftAt.Time))
	}

	dbRoomMembers, err := models.RoomMembers(mods...).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("query room members: %w", err)
	}

	return dbRoomMembersToRoomMembers(dbRoomMembers), nil
}

func dbRoomMembersToRoomMembers(dbRoomMembers []*models.RoomMember) []*room_members.RoomMembers {
	result := make([]*room_members.RoomMembers, len(dbRoomMembers))
	for i, db := range dbRoomMembers {
		var joinedAt null.Time
		if db.JoinedAt.Valid {
			joinedAt = null.TimeFrom(db.JoinedAt.Time)
		}

		result[i] = &room_members.RoomMembers{
			ID:       fmt.Sprintf("%d", db.ID),
			RoomID:   fmt.Sprintf("%d", db.RoomID),
			UserID:   fmt.Sprintf("%d", db.UserID),
			JoinedAt: joinedAt.Time,
			LeftAt:   db.LeftAt.Time,
		}
	}
	return result
}
