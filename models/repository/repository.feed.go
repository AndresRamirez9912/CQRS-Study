package repository

import (
	"database/sql"
	"events/models"
	"fmt"

	_ "github.com/lib/pq"
)

type IRepositoryFeed interface {
	Close()
	InsertFeed(feed *models.Feed) error
	ListFeeds() ([]*models.Feed, error)
}

type FeedRepository struct {
	db *sql.DB
}

func NewFeed() (*FeedRepository, error) {
	url := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable", "postgres", "45665482", "Feeds", "postgres")
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &FeedRepository{
		db: db,
	}, nil
}

func (repo *FeedRepository) Close() {
	repo.db.Close()
}

func (repo *FeedRepository) InsertFeed(feed *models.Feed) error {
	_, err := repo.db.Exec("INSERT INTO feeds (id, title, description) VALUES ($1,$2,$3)", feed.Id, feed.Title, feed.Description)
	return err
}

func (repo *FeedRepository) ListFeeds() ([]*models.Feed, error) {
	rows, err := repo.db.Query("SELECT * FROM  feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []*models.Feed{}
	for rows.Next() {
		feed := &models.Feed{}
		err := rows.Scan(&feed.Id, &feed.Title, &feed.Description, &feed.CreatedAt)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}
