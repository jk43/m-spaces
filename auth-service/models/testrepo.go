package models

import (
	"errors"
)

type TestDBRepo struct{}

func (repo *TestDBRepo) InsertCredentials(cred *Credentials) (int, error) {
	if cred.Email == "error@test.com" {
		return 0, errors.New("Error in InsertCredentials")
	}
	return 1, nil
}

func (repo *TestDBRepo) GetCredentialsWithEmail(email string) (*Credentials, error) {
	if email == "error@test.com" {
		return nil, errors.New("GetCredentialsWithEmail")
	}
	if email == "test@test.com" {
		// user := service.User{
		// 	Email:     "test@test.com",
		// 	ID:  "63b4b225132f9fad1a325f0f",
		// 	FirstName: "First Name",
		// 	LastName:  "Last Name",
		// }
		return &Credentials{
			// User:     &user,
			// Password: "$2a$10$BBsvVBQ8UdyLNkL3./7rTeRbKI/dfq0lEvFcRO7vWithANRVa/8ZVHO",
		}, nil
	}
	return nil, nil
}

func (repo *TestDBRepo) GetCredentialsWithUserID(oid string) (*Credentials, error) {
	// user := service.User{
	// 	FirstName: "First Name",
	// 	LastName:  "Last Name",
	// 	Email:     "test@test.com",
	// }
	return &Credentials{
		// User: &user,
	}, nil
}
