package entity

import "time"

type Commit struct {
	ID             int       `db:"id"`
	RepoName       string    `db:"repo_name"`
	CommitID       string    `db:"commit_id"`
	CommitMessage  string    `db:"commit_message"`
	AuthorUsername string    `db:"author_username"`
	AuthorEmail    string    `db:"author_email"`
	Payload        string    `db:"payload"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
