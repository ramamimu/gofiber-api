package repository

import (
	"context"
	"errors"
	"strconv"
	"time"
)

type Thread struct {
	ID         string    `json:"id"`
	Created    time.Time `json:"created"`
	LastUpdate time.Time `json:"last_update"`
	Author     string    `json:"author"`
	Content    string    `json:"content"`
	IsEdited   bool      `json:"is_edited"`
}

type Db struct {
	threads   map[string]Thread
	increment int
}

func (db *Db) Init() {
	db.increment = 0
	db.threads = make(map[string]Thread)
}

func (db *Db) Clear() {
	for t := range db.threads {
		delete(db.threads, t)
	}
}

func (db *Db) GetThreadByID(id string) (Thread, error) {
	val, ok := db.threads[id]
	if !ok {
		return Thread{}, errors.New("thread not found")
	}
	return val, nil
}

func (db *Db) GetThreadsEntity() map[string]Thread {
	return db.threads
}

func (db *Db) GetThreads(ctx context.Context) []Thread {
	t := []Thread{}
	for _, thread := range db.threads {
		t = append(t, thread)
	}
	return t
}

func (db *Db) AddThread(ctx context.Context, author string, content string) (string, error) {
	thread := Thread{
		ID:         strconv.Itoa(db.increment),
		Created:    time.Now(),
		LastUpdate: time.Now(),
		Author:     author,
		Content:    content,
		IsEdited:   false,
	}
	db.threads[strconv.Itoa(db.increment)] = thread
	db.increment++
	return thread.ID, nil
}

func (db *Db) EditThread(ctx context.Context, id string, content string) error {
	val, ok := db.threads[id]
	if !ok {
		return errors.New("thread is not available")
	}

	val.Content = content
	val.LastUpdate = time.Now()
	val.IsEdited = true

	db.threads[id] = val

	return nil
}

func (db *Db) DeleteThread(ctx context.Context, id string) error {
	_, ok := db.threads[id]
	if !ok {
		return errors.New("thread is not available")
	}

	delete(db.threads, id)
	return nil
}
