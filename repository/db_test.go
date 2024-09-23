package repository_test

import (
	"context"
	"gofiber-api/repository"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DbTestSuite struct {
	suite.Suite
	db repository.Db
}

////////////////////////////
// trial each suite state //
////////////////////////////

// called only in beginning
func (suite *DbTestSuite) SetupSuite() {
	suite.db.Init()
}

// called in every test function
func (suite *DbTestSuite) SetupTest() {
	suite.db.Clear()
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

////////////////////////////

func (s *DbTestSuite) TestCreateThread() {
	id, err := s.db.AddThread(context.Background(), "the-author", "the-content")
	s.Nil(err)

	thread, err := s.db.GetThreadByID(id)
	s.Nil(err)
	s.Equal("the-author", thread.Author)
	s.Equal("the-content", thread.Content)
	s.NotEmpty(thread.Created)
	s.NotEmpty(thread.LastUpdate)
	s.False(thread.IsEdited)
}
