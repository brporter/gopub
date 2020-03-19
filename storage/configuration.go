package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/brporter/gopub/models"
)

type StorageConfiguration struct {
	connectionString string
	database         string
	collection       string
}

func NewStorageConfiguration(path string) models.IConfiguration {
	var s StorageConfiguration

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	lines := make([]string, 0, 3) // we expect three lines

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) > 3 {
		panic("unexpected number of storage configuration elements")
	}

	s.connectionString = lines[0]
	s.database = lines[1]
	s.collection = lines[2]

	return &s
}

func (s *StorageConfiguration) GetSecret(name string) (*string, error) {
	lname := strings.ToLower(name)

	switch lname {
	case "storage":
		return &s.connectionString, nil
	case "database":
		return &s.database, nil
	case "collection":
		return &s.collection, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown secret named %v", name))
	}

	return nil, errors.New("unable to determine secret")
}
