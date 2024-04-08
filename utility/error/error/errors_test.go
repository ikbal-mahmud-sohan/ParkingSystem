package utility

import (
	"errors"
	"net/http"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestHandleServiceErr(t *testing.T) {
	user := User{
		FirstName: "",
		LastName:  "",
	}

	valErr := user.Validate()

	testCases := []struct {
		name         string
		inputError   error
		expectedCode codes.Code
		expectedMsg  string
		expectedErrs map[string]string
	}{
		{
			name:         "nil error input",
			inputError:   nil,
			expectedCode: 0,
			expectedMsg:  "",
			expectedErrs: nil,
		},
		{
			name:         "validation error",
			inputError:   valErr,
			expectedCode: http.StatusUnprocessableEntity,
			expectedMsg:  "FirstName: cannot be blank; LastName: cannot be blank.",
			expectedErrs: map[string]string{"FirstName": "cannot be blank", "LastName": "cannot be blank"},
		},
		{
			name:         "postgres error",
			inputError:   errors.New("pq: error message"),
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Internal server error",
			expectedErrs: nil,
		},
		{
			name:         "json error",
			inputError:   errors.New("json: error message"),
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Internal server error",
			expectedErrs: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := HandleServiceErr(tc.inputError)

			if err == nil {
				assert.Nil(t, err)
			}
		})
	}
}

type User struct {
	FirstName string
	LastName  string
}

func (a User) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required),
		validation.Field(&a.LastName, validation.Required),
	)
}
