package threads

import (
	"context"
	repo "gofiber-api/repository"
)

type RepositoryThread interface {
	GetThreads(ctx context.Context) []repo.Thread
	AddThread(ctx context.Context, author string, content string) (string, error)
	EditThread(ctx context.Context, id string, content string) error
	DeleteThread(ctx context.Context, id string) error
}

// test this with mock tomorrow

type ThreadService struct {
	RepositoryThread
}

func NewThread(r RepositoryThread) *ThreadService {
	return &ThreadService{
		RepositoryThread: r,
	}
}

func (t *ThreadService) GetAll(ctx context.Context) []repo.Thread {
	return t.GetThreads(ctx)
}

func (t *ThreadService) Add(ctx context.Context, author string, content string) error {
	_, err := t.AddThread(ctx, author, content)
	return err
}

func (t *ThreadService) Edit(ctx context.Context, id string, content string) error {
	return t.EditThread(ctx, id, content)
}

func (t *ThreadService) Delete(ctx context.Context, id string) error {
	return t.DeleteThread(ctx, id)
}
