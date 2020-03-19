module github.com/brporter/gopub/storage

go 1.13

require (
	github.com/brporter/gopub/models v0.0.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/google/uuid v1.1.1
	gitlab.com/golang-commonmark/markdown v0.0.0-20191127184510-91b5b3c99c19
	go.mongodb.org/mongo-driver v1.3.0
)

replace github.com/brporter/gopub/models v0.0.0 => /Users/brporter/projects/gopub/models

replace github.com/brporter/gopub/storage v0.0.0 => /Users/brporter/projects/gopub/storage
