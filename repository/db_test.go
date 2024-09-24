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

func (s *DbTestSuite) TestGetThread() {
	// add thread to fb
	s.db.AddThread(context.Background(), "the-author", "the contents")

	threads := s.db.GetThreads(context.Background())
	s.NotEmpty(threads)
	s.Equal(1, len(threads))
	s.Equal("the-author", threads[0].Author)
	s.Equal("the contents", threads[0].Content)
	s.NotEmpty(threads[0].Created)
	s.False(threads[0].IsEdited)
	s.NotEmpty(threads[0].LastUpdate)

	s.db.AddThread(context.Background(), "the-author-1", "the contents 2")
	threads = s.db.GetThreads(context.Background())
	s.NotEmpty(threads)
	s.Equal(2, len(threads))
	s.Equal("the-author", threads[0].Author)
	s.Equal("the contents", threads[0].Content)
	s.NotEmpty(threads[0].Created)
	s.False(threads[0].IsEdited)
	s.NotEmpty(threads[0].LastUpdate)
}

func (s *DbTestSuite) TestEditThread() {
	// add thread to fb
	s.db.AddThread(context.Background(), "the-author", "the contents for edit")

	threads := s.db.GetThreads(context.Background())
	// edit thread
	err := s.db.EditThread(context.Background(), threads[0].ID, "the edited content")
	s.NoError(err)

	threads = s.db.GetThreads(context.Background())
	s.NotEmpty(threads)
	s.Equal(1, len(threads))
	s.Equal("the-author", threads[0].Author)
	s.Equal("the edited content", threads[0].Content)
	s.NotEmpty(threads[0].Created)
	s.True(threads[0].IsEdited)
	s.NotEmpty(threads[0].LastUpdate)
}

func (s *DbTestSuite) TestDeleteThread() {
	// add thread to db
	s.db.AddThread(context.Background(), "the-author", "the contents")
	s.db.AddThread(context.Background(), "the-author-2", "the content 2")
	s.db.AddThread(context.Background(), "the-author-3", "the content 3")

	// delete thread
	s.db.DeleteThread(context.Background(), "1")

	threads := s.db.GetThreads(context.Background())
	s.NotEmpty(threads)
	s.Equal(2, len(threads))

	thread, err := s.db.GetThreadByID("1")
	s.Empty(thread)
	s.Error(err)
}
