package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brporter/gopub/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FakeRepo struct {
	ErrorOnFetch  bool
	ErrorOnSave   bool
	ErrorOnDelete bool
}

func (f *FakeRepo) Open() error {
	return nil
}

func (f *FakeRepo) Close() error {
	return nil
}

func (f *FakeRepo) Save(p *models.Post) error {

	if f.ErrorOnSave {
		return errors.New("an error occurred")
	}

	return nil
}

func (f *FakeRepo) FetchOne(id *uuid.UUID) (*models.Post, error) {
	if f.ErrorOnFetch {
		return nil, errors.New("an error occurred")
	}

	p, err := models.NewPost("A Post", "A Title")

	if err != nil {
		panic(err)
	}

	p.PostId = *id

	return p, nil
}

func (f *FakeRepo) Remove(id *uuid.UUID) error {
	if f.ErrorOnDelete {
		return errors.New("an error occurred")
	}

	return nil
}

func TestGetPostHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/cdc0be7a-60dc-11ea-a0d8-acde48001122", nil)

	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	rr := httptest.NewRecorder()

	router.Handle("/posts/{postId}", service.GetPostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("fetching post failed with %v, expected %v, msg: %v", rr.Code, http.StatusOK, rr.Body)
	}
}

func TestGetPostHandler_BadUuid(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/cdc0be7a-60dc-11ea-a0d8", nil)

	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	rr := httptest.NewRecorder()

	router.Handle("/posts/{postId}", service.GetPostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("fetching post failed with %v, expected %v, msg: %v", rr.Code, http.StatusInternalServerError, rr.Body)
	}
}

func TestGetPostHandler_ErrorOnFetch(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/cdc0be7a-60dc-11ea-a0d8-acde48001122", nil)

	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	repo.ErrorOnFetch = true

	rr := httptest.NewRecorder()

	router.Handle("/posts/{postId}", service.GetPostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("fetching post failed with %v, expected %v, msg: %v", rr.Code, http.StatusInternalServerError, rr.Body)
	}
}

func TestPutPostHandler(t *testing.T) {
	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	rr := httptest.NewRecorder()

	p, _ := models.NewPost("A Title", "A Body")
	ps := p.ToJson()

	req, err := http.NewRequest("PUT", "/posts", strings.NewReader(*ps))

	if err != nil {
		t.Fatal(err)
	}

	router.Handle("/posts", service.PutPostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("fetching post failed with %v, expected %v, msg: %v", rr.Code, http.StatusOK, rr.Body)
	}
}
func TestPutPostHandler_BadJson(t *testing.T) {
	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	rr := httptest.NewRecorder()

	p, _ := models.NewPost("A Title", "A Body")
	ps := p.ToJson()
	*ps += " INVALID "

	req, err := http.NewRequest("PUT", "/posts", strings.NewReader(*ps))

	if err != nil {
		t.Fatal(err)
	}

	router.Handle("/posts", service.PutPostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("saving post returned status code %v, expected %v, msg: %v", rr.Code, http.StatusOK, rr.Body)
	}
}

func TestPutPostHandler_ErrorOnSave(t *testing.T) {
	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	repo.ErrorOnSave = true
	rr := httptest.NewRecorder()

	p, _ := models.NewPost("A Title", "A Body")
	ps := p.ToJson()

	req, err := http.NewRequest("PUT", "/posts", strings.NewReader(*ps))

	if err != nil {
		t.Fatal(err)
	}

	router.Handle("/posts", service.PutPostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("saving post returned status code %v, expected %v, msg: %v", rr.Code, http.StatusOK, rr.Body)
	}
}

func TestDeletePostHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/posts/cdc0be7a-60dc-11ea-a0d8-acde48001122", nil)

	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	rr := httptest.NewRecorder()

	router.Handle("/posts/{postId}", service.DeletePostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("deleting post failed with %v, expected %v, msg: %v", rr.Code, http.StatusOK, rr.Body)
	}
}

func TestDeletePostHandler_BadUUID(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/posts/cdc0be7a-60dc-11ea-a0d8-acde", nil)

	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	rr := httptest.NewRecorder()

	router.Handle("/posts/{postId}", service.DeletePostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("deleting post failed with %v, expected %v, msg: %v", rr.Code, http.StatusInternalServerError, rr.Body)
	}
}

func TestDeletePostHandler_ErrorOnRemove(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/posts/cdc0be7a-60dc-11ea-a0d8-acde48001122", nil)

	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()

	var service PostService
	repo := new(FakeRepo)
	repo.ErrorOnDelete = true

	rr := httptest.NewRecorder()

	router.Handle("/posts/{postId}", service.DeletePostHandler(repo))

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("deleting post failed with %v, expected %v, msg: %v", rr.Code, http.StatusInternalServerError, rr.Body)
	}
}
