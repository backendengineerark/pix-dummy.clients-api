package entity

import (
	"time"

	"github.com/backendengineerark/clients-api/pkg/customerrors"
	"github.com/backendengineerark/clients-api/pkg/dates"
	"github.com/google/uuid"
)

type Client struct {
	Id        string
	Name      string
	Document  string
	BirthDate time.Time
	CreatedAt time.Time
}

func NewClient(name string, document string, bithDate string) (*Client, []customerrors.Error) {
	errors := []customerrors.Error{}

	birthDateFormatted, err := dates.GenerateBirthDate(bithDate)
	if err != nil {
		errors = append(errors, *customerrors.NewError(customerrors.INVALID_PARAM, err.Error()))
	}

	client := &Client{
		Id:        uuid.New().String(),
		Name:      name,
		Document:  document,
		BirthDate: birthDateFormatted,
		CreatedAt: time.Now(),
	}

	validationErrors := client.IsValid()
	if len(validationErrors) > 0 {
		errors = append(errors, validationErrors...)
		return nil, errors
	}

	return client, nil
}

func (c *Client) IsValid() []customerrors.Error {
	errorsList := []customerrors.Error{}

	if c.Id == "" {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Id is required"))
	}

	if c.Name == "" {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Name is required"))
	} else if len(c.Name) < 3 {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Name should have 3 or more characters"))
	}

	if c.Document == "" {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Document is required"))
	} else if len(c.Document) != 11 {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Document should have 11 numbers"))
	}

	if c.BirthDate.IsZero() {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Birth date is required"))
	}

	if c.CreatedAt.IsZero() {
		errorsList = append(errorsList, *customerrors.NewError(customerrors.INVALID_PARAM, "Created at is required"))
	}

	return errorsList
}
