package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brporter/gopub/models"
	"github.com/brporter/gopub/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type IPostService interface {
	GetPostHandler(repo storage.IPostRepo) http.HandlerFunc
	PutPostHandler(repo storage.IPostRepo) http.HandlerFunc
	DeletePostHandler(repo storage.IPostRepo) http.HandlerFunc
}

type PostService struct{}

func (s *PostService) GetPostHandler(repo storage.IPostRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if vars["postId"] != "" {
			postId, err := uuid.Parse(vars["postId"])

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// a specific post was requested
			p, err := repo.FetchOne(&postId)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			j := p.ToJson()

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, *j)
		} else {
			pageSize := 10
			publishDate := time.Now().UTC()

			publishDateParam := r.URL.Query().Get("publishDate")

			if publishDateParam != "" {
				fmt.Printf("Fetching documents older than %v\n", publishDateParam)
				parsedDateTime, err := time.Parse(time.RFC3339, publishDateParam)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				publishDate = parsedDateTime
			}

			p, err := repo.FetchMany(publishDate, pageSize)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			j := models.ToJson(p)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, *j)
		}
	}
}

func (s *PostService) PutPostHandler(repo storage.IPostRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// body should be a json of a post object
		bodyString := string(body)
		result, err := models.FromJson(&bodyString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = repo.Save(result)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *PostService) DeletePostHandler(repo storage.IPostRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		postId, err := uuid.Parse(vars["postId"])

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = repo.Remove(&postId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
