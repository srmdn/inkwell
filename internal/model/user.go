package model

import "time"

type User struct {
	ID         int64     `db:"id"`
	Email      string    `db:"email"`
	PasswdHash string    `db:"passwd_hash"`
	CreatedAt  time.Time `db:"created_at"`
}
