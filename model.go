package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./example.sqlite3")
	if err != nil {
		panic(err)
	}
}

// 記事の全件検索
func getPosts(limit int) (posts []Post, err error) {
	stmt := "SELECT id, title, body, author FROM posts LIMIT ?"
	rows, err := Db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return posts, nil
}

// 記事の1件検索
func retrieve(id int) (post Post, err error) {
	post = Post{}
	stmt := "SELECT * FROM posts WHERE id = ?"
	err = Db.QueryRow(stmt, id).Scan(&post.Id, &post.Title, &post.Body, &post.Author)
	return post, nil
}

func (post *Post) create() error {
	stmt := "INSERT INTO posts (title, body, author) VALUES (?, ?, ?) RETURNING id"
	err := Db.QueryRow(stmt, post.Title, post.Body, post.Author).Scan(&post.Id)
	return err
}

//更新処理
func (post *Post) update() error {
	stmt := "UPDATE posts set title=?, body=?, author=? WHERE id=?"
	_, err := Db.Exec(stmt, post.Title, post.Body, post.Author, post.Id)
	return err
}

//削除処理
func (post *Post) delete() error {
	stmt := "DELETE FROM posts WHERE id=?"
	_, err := Db.Exec(stmt, post.Id)
	return err
}
