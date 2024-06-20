package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type eventType int32

const (
	likes eventType = iota
	views
)

const (
	topPostLim int = 5
	topUserLim int = 3
)

type post struct {
	id          uint64
	likes       int64
	views       int64
	post_author int32
}

type userInfo struct {
	postAuthor int32
	likes      int64
}

type dbWrapper struct {
	*sql.DB
}

func newDBWrapper() (dbWrapper, error) {
	clickhousePort := os.Getenv("STAT_DB_TCP_PORT")
	dsn := fmt.Sprintf("clickhouse://default:@clickhouse:%s/default", clickhousePort)
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return dbWrapper{}, fmt.Errorf("error connecting to db: %v", err.Error())
	}

	// Ping the database to ensure the connection is established
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		return dbWrapper{}, fmt.Errorf("error pinging db: %v", err.Error())
	}

	fmt.Println("Connected to ClickHouse successfully!")

	return dbWrapper{db}, nil
}

func (db dbWrapper) countLikesViews(id uint64) (likes_, views_ int64, err error) {
	log.Printf("counting likes and views for %v \n", id)

	query := fmt.Sprintf(`
        SELECT
			post_id,
            sumIf(1, event_type = %d) AS likes,
            sumIf(1, event_type = %d) AS views
        FROM
            event_table
        WHERE
            post_id = ?
        GROUP BY
            post_id
    `, likes, views)

	var post_id int64
	err = db.QueryRow(query, id).Scan(&post_id, &views_, &likes_)
	return
}

func (db dbWrapper) getTopPosts(e eventType) ([]post, error) {
	log.Println("counting top user by post likes/views")

	// var order string = "likes"
	// if e == views {
	// 	order = "views"
	// }

	// query := fmt.Sprintf(`
	// SELECT
	//     post_id,
	//     sumIf(1, event_type = %d) AS likes,
	//     sumIf(1, event_type = %d) AS views
	// 	post_author
	// FROM
	//     event_table
	// GROUP BY
	//     post_id
	// ORDER BY
	//     %s DESC
	// LIMIT
	//     %d;
	// `, likes, views, order, topNum)
	query := fmt.Sprintf(`
	SELECT
		post_id,
		post_author,
		COUNT() AS count
	FROM
		event_table
	WHERE
		event_type = %d
	GROUP BY
		post_id, post_author
	ORDER BY
		count DESC
	LIMIT %d
	`, e, topPostLim)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]post, 0, topPostLim)
	for rows.Next() {
		var p post
		var err error

		if e == likes {
			err = rows.Scan(&p.id, &p.likes, &p.post_author)
		} else {
			err = rows.Scan(&p.id, &p.views, &p.post_author)
		}

		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil

}

func (db dbWrapper) getTopUsers() ([]userInfo, error) {
	log.Println("counting top user by post likes/views")

	query := fmt.Sprintf(`
    SELECT 
		event_author AS user_id, 
		sumIf(1, event_type = %d) AS total_likes
	FROM 
		event_table
	GROUP BY 
		event_author
    ORDER BY 
        total_likes DESC
    LIMIT 
        %d;
    `, likes, topUserLim)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]userInfo, 0, topUserLim)
	for rows.Next() {
		var info userInfo
		err := rows.Scan(&info.postAuthor, &info.likes)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}
