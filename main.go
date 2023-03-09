package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/posts/", postsHandler)
	server.ListenAndServe()
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	path := path.Base(r.URL.Path)

	// pathがallの場合全件取得
	if path == "all" {
		posts, err := getPosts(100)
		if err != nil {
			return err
		}
		output, err := json.Marshal(&posts)
		if err != nil {
			return err
		}
		w.Write(output)
		return nil
	}

	// pathがidの場合id検索
	id, err := strconv.Atoi(path)
	if err != nil {
		return err
	}

	post, err := retrieve(id)
	if err != nil {
		return err
	}

	output, err := json.Marshal(&post)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(output)
	return nil
}

func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	contentLength := r.ContentLength
	contentBody := make([]byte, contentLength)
	r.Body.Read(contentBody)

	var post Post
	err = json.Unmarshal(contentBody, &post)
	if err != nil {
		return err
	}
	err = post.create()
	if err != nil {
		return err
	}

	output, err := json.Marshal(&post)
	if err != nil {
		return
	}

	w.WriteHeader(200)
	w.Write(output)
	return nil
}

func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return err
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	contentLength := r.ContentLength
	contentBody := make([]byte, contentLength)
	r.Body.Read(contentBody)

	err = json.Unmarshal(contentBody, &post)
	if err != nil {
		return err
	}

	err = post.update()
	if err != nil {
		return err
	}

	output, err := json.Marshal(post)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	w.Write(output)
	return nil
}

func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return err
	}

	post, err := retrieve(id)
	if err != nil {
		return
	}

	err = post.delete()
	if err != nil {
		return err
	}

	output, err := json.Marshal(&post)
	w.WriteHeader(200)
	w.Write(output)
	return nil
}
