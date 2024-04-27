package server

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound                = errors.New("Object not found in database")
	ErrCreatingPost            = errors.New("Error creating post")
	ErrInsufficientPermissions = errors.New("insufficient access rights")
)

type PostIntRep struct {
	gorm.Model
	Author  int32
	Content string
}

type dbWrapper struct {
	*gorm.DB
}

func NewDBWrapper() (dbWrapper, error) {
	username := os.Getenv("CONTENT_DB_USER")
	password := os.Getenv("CONTENT_DB_PASSWORD")
	dbname := os.Getenv("CONTENT_DB")

	dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s sslmode=disable", username, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return dbWrapper{}, nil
	}

	err = db.AutoMigrate(&PostIntRep{})
	if err != nil {
		return dbWrapper{}, err
	}
	return dbWrapper{db}, nil
}

func (db *dbWrapper) GetPostObj(id uint, author int32) (*PostIntRep, error) {
	log.Println("GetPostObj query")

	post := PostIntRep{Model: gorm.Model{ID: id}}
	result := db.First(&post, "ID=?", id)
	if result == nil {
		return nil, ErrNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}
	if post.Author != author {
		return nil, ErrInsufficientPermissions
	}
	return &post, nil
}

func (db *dbWrapper) CreatePost(post *PostIntRep) error {
	log.Println("CreatePost query")

	result := db.Create(post)
	if result == nil {
		return ErrCreatingPost
	}
	return result.Error
}

func (db *dbWrapper) UpdatePost(post *PostIntRep, id uint) error {
	log.Println("UpdatePost query")

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

func (db *dbWrapper) DeletePost(postAuthor int32, id uint) error {
	log.Println("DeletePost query")

	// Chech access rights
	_, err := db.GetPostObj(id, postAuthor)
	if err != nil {
		return err
	}
	return db.Delete(&PostIntRep{}, id).Error
}

func (db *dbWrapper) GetPosts(offset int, batchSize int, author int32) ([]PostIntRep, error) {
	log.Println("GetPosts query")

	var posts []PostIntRep
	result := db.Where("author = ?", author).Limit(int(batchSize)).Find(&posts)
	return posts, result.Error
}
