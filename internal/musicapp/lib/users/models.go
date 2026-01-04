package users

import (
	"time"

	"github.com/aarondl/null/v8"
)

// User - Clean domain model (no DB tags)
type User struct {
	ID          string
	Username    string
	Email       string
	DisplayName string
	Gender      Gender
	CreatedAt   time.Time
}

// Gender enum
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// UserQueryFilter - uses null types for optional filters
type UserQueryFilter struct {
	IDs      []string
	Username null.String
	Email    null.String
	Gender   null.String

	// Sorting
	OrderBy null.String // "created_at", "username"
	Sort    null.String // "asc", "desc", "ASC", "DESC"
	Limit   null.Int
	Offset  null.Int
}

// UpdateUser - nullable fields for partial updates
type UpdateUser struct {
	IDs         []string // Which users to update
	Username    null.String
	DisplayName null.String
	Gender      null.String
}
