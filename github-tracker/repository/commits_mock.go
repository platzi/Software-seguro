package repository

import (
	"context"
	"github-tracker/github-tracker/repository/entity"

	"github.com/stretchr/testify/mock"
)

type MockCommit struct {
	*mock.Mock
	Commit
}

func (m MockCommit) Insert(ctx context.Context, commit *entity.Commit) (err error) {
	results := m.Called(ctx, commit)
	return results.Error(0)
}

func (m MockCommit) GetCommitByAuthorEmail(ctx context.Context, email string) (commits []entity.Commit, err error) {
	results := m.Called(ctx, email)
	return results.Get(0).([]entity.Commit), results.Error(1)
}
