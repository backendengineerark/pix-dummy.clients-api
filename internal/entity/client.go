package entity

import (
	"errors"

	"github.com/google/uuid"
)

type Client struct {
	Id        string
	Name      string
	Document  string
	BirthDate string
}

func NewClient(name, document, bithDate string) (*Client, []error) {
	client := &Client{
		Id:        uuid.New().String(),
		Name:      name,
		Document:  document,
		BirthDate: bithDate,
	}

	errors := client.IsValid()
	if len(errors) > 0 {
		return nil, errors
	}

	return client, nil
}

func (c *Client) IsValid() []error {
	errorsList := []error{}

	if c.Id == "" {
		errorsList = append(errorsList, errors.New("id is required"))
	}

	if c.Name == "" {
		errorsList = append(errorsList, errors.New("name is required"))
	} else if len(c.Name) < 3 {
		errorsList = append(errorsList, errors.New("name should have 3 or more characters"))
	}

	if c.Document == "" {
		errorsList = append(errorsList, errors.New("document is required"))
	} else if len(c.Document) != 11 {
		errorsList = append(errorsList, errors.New("document should have 11 numbers"))
	}

	// validate birth date

	return errorsList
}
