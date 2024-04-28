package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

type dbWrapper struct {
	*sql.DB
}

type EngagementEvent struct {
	EventID uint32 `json:"id"`
	Type    int    `json:"type"`
	PostID  int64  `json:"post"`
	Author  int    `json:"author"`
}

func newDBWrapper() (dbWrapper, error) {
	clickhousePort := os.Getenv("STAT_DB_PORT")
	dsn := fmt.Sprintf("clickhouse://default:@clickhouse:%s/default", clickhousePort)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return dbWrapper{}, err
	}

	// Ping the database to ensure the connection is established
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return dbWrapper{}, err
	}

	fmt.Println("Connected to ClickHouse successfully!")

	return dbWrapper{db}, nil
}

// func (db *dbWrapper) Insert(msg []byte) error {
// 	var event EngagementEvent
// 	err := json.Unmarshal(msg, &event)
// 	if err != nil {
// 		fmt.Printf("Error decoding message: %v\n", err)
// 		return err
// 	}

// 	_, err = db.ExecContext(context.Background(), `
// 	INSERT INTO events
// 	(id, type, post, author)
// 	VALUES (?, ?, ?, ?)`,
// 		event.EventID,
// 		event.Type,
// 		event.PostID,
// 		event.Author,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
