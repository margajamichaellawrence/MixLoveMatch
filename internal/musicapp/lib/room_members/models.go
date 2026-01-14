package room_members

import (
	"time"

	"github.com/aarondl/null/v8"
)

type RoomMembers struct {
	ID       string
	RoomID   string
	UserID   string
	JoinedAt time.Time
	LeftAt   time.Time
}

type RoomMemberQueryFilter struct {
	IDs      []string
	RoomID   null.String
	UserID   null.String
	JoinedAt null.Time
	LeftAt   null.Time
}

type UpdateRoomMember struct {
	IDs      []string
	RoomID   null.String
	UserID   null.String
	JoinedAt null.Time
	LeftAt   null.Time
}
