package server

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrInsufficientPermissions = errors.New("insufficient access rights")

type PostIntRep struct {
	gorm.Model
	Author  string
	Content string
}

type dbWrapper struct {
	*gorm.DB
}

func NewDBWrapper() (dbWrapper, error) {
	username := os.Getenv("CONTENT_DB_USER")
	password := os.Getenv("CONTENT_DB_PASSWORD")
	// port := os.Getenv("CONTENT_DB_PORT")
	dbname := os.Getenv("CONTENT_DB")

	dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s sslmode=disable", username, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return dbWrapper{}, nil
	}

	err = db.AutoMigrate(&PostIntRep{})
	if err != nil {
		return dbWrapper{}, nil
	}
	return dbWrapper{db}, nil
}

func (db *dbWrapper) GetPostObj(id uint, author string) (*PostIntRep, error) {
	post := PostIntRep{Model: gorm.Model{ID: id}}
	result := db.First(&post, "ID=?", id)

	if result.Error != nil {
		return nil, result.Error
	}
	if post.Author != author {
		return nil, ErrInsufficientPermissions
	}
	return &post, nil
}

func (db *dbWrapper) CreatePost(post *PostIntRep) error {
	result := db.Create(post)
	return result.Error
}

func (db *dbWrapper) UpdatePost(post *PostIntRep, id uint) error {
	// Check access rights
	_, err := db.GetPostObj(id, post.Author)
	if err != nil {
		return err
	}
	// update query
	result := db.Model(&PostIntRep{}).Where("ID = ?", id).Updates(PostIntRep{
		Author:  post.Author,
		Content: post.Content,
	})
	return result.Error
}

func (db *dbWrapper) DeletePost(postAuthor *string, id uint) error {
	// Chech access rights
	_, err := db.GetPostObj(id, *postAuthor)
	if err != nil {
		return err
	}
	return db.Delete(&PostIntRep{}, id).Error
}

func (db *dbWrapper) GetPosts(offset int, batchSize int, author string) ([]PostIntRep, error) {
	var posts []PostIntRep
	result := db.Where("author = ?", author).Limit(int(batchSize)).Find(&posts)
	return posts, result.Error
}
