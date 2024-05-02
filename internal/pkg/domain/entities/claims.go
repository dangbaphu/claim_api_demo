package entities

import "time"

type Claim struct {
	ID        string        `bson:"id,omitempty"`
	UserID    string        `bson:"user_id,omitempty"`
	Ammount   int           `bson:"ammount,omitempty"`
	Documents []FileStorage `bson:"documents,omitempty"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
}
