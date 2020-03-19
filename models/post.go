package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gitlab.com/golang-commonmark/markdown"
)

type IPost interface {
	ToHtml() *string
	ToJson() *string
}

type Post struct {
	PostId      uuid.UUID `bson:"_id,omitempty" json:"postId"`
	PublishDate time.Time `bson:"publishDate" json:"publishDate"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Summary     string    `json:"summary"`
}

func NewPost(title string, body string) (*Post, error) {
	retVal := new(Post)
	id, err := uuid.NewUUID()

	if err == nil {
		retVal.PostId = id
		retVal.Title = title
		retVal.Body = body
	} else {
		retVal = nil
	}

	return retVal, err
}

func FromJson(serialized *string) (*Post, error) {
	var result Post
	err := json.Unmarshal([]byte(*serialized), &result)

	return &result, err
}

func (p *Post) ToHtml() *string {
	// convert markdown to html
	md := markdown.New(markdown.XHTMLOutput(true))
	result := md.RenderToString([]byte(p.Body))

	return &result
}

func ToJson(posts []*Post) *string {
	result, err := json.Marshal(posts)

	if err != nil {
		panic(err)
	}

	sr := string(result)

	return &sr
}

func (p *Post) ToJson() *string {
	result, err := json.Marshal(p)

	if err != nil {
		panic(err)
	}

	sr := string(result)

	return &sr
}
