package main

import (
	"flag"
	"net/http"

	"github.com/brporter/gopub/services"
	"github.com/brporter/gopub/storage"
	"github.com/gorilla/mux"
)

func main() {

	storageCfg := flag.String("storage", "storage.cfg", "The path to the storage.cfg file containing storage repository connection details.")
	flag.Parse()

	r := mux.NewRouter()
	var s services.PostService

	storageConfig := storage.NewStorageConfiguration(*storageCfg)
	repo, _ := storage.NewMongoRepo(true, storageConfig)
	defer repo.Close()

	r.Handle("/posts/{postId}", s.DeletePostHandler(repo)).Methods("DELETE")
	r.Handle("/posts/{postId}", s.GetPostHandler(repo)).Methods("GET")
	r.Handle("/posts", s.GetPostHandler(repo)).Methods("GET")
	r.Handle("/posts", s.PutPostHandler(repo)).Methods("PUT")

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid", http.StatusBadRequest)
	}))

	r.Use(services.CORSMiddleware)
	r.Use(services.AuthMiddleware)

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
