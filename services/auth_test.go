package services

import (
	"fmt"
	"testing"

	"podcast/database"
	"podcast/gateway"
	"podcast/models"
	"podcast/repositories"
	"podcast/tests"
	"podcast/types"
)

var (
	as       = NewAuthService(repositories.NewUsersRepository(database.DB), NewEmailService(gateway.NewSMTPGateway()))
	user     models.User
	name     string = "John Doe"
	email    string = "john@testing.com"
	password string = "123456789"
)

func TestSignup(t *testing.T) {
	defer tests.Teardown()

	u, err := as.Signup(types.SignupInput{Email: email, Password: password})
	if u.Email != email && err != nil {
		t.Fatalf("should signup and return the proper email, want %s, found %s", email, user.Email)
	}
}

func TestSignin(t *testing.T) {
	defer tests.Teardown()

	as.Signup(types.SignupInput{Email: email, Password: password})

	u, err := as.Signin(types.SigninInput{Email: email, Password: password})
	if u.Email != email && err != nil {
		t.Fatalf("should signin and return the proper user email, want %s, found %s", email, user.Email)
	}
}

func TestUpdateDetails(t *testing.T) {
	defer tests.Teardown()

	u, _ := as.Signup(types.SignupInput{Email: email, Password: password})

	u, err := as.UpdateDetails(fmt.Sprint(u.ID), types.UpdateDetailsInput{
		Email: "updated@test.com",
	})
	if u.Email != "updated@test.com" || err != nil {
		t.Fatalf("should update the user email, want %s, found %s", "updated@test.com", u.Email)
	}

	u, err = as.UpdateDetails(fmt.Sprint(u.ID), types.UpdateDetailsInput{
		Email: "updatetest.com",
	})
	if u.Email == "updatetest.com" || err == nil {
		t.Fatalf("should not update the user email, want %s, found %s", "updated@test.com", u.Email)
	}
}

func TestUpdatePassword(t *testing.T) {
	defer tests.Teardown()

	u, _ := as.Signup(types.SignupInput{Email: email, Password: password})

	upd, err := as.UpdatePassword(fmt.Sprint(u.ID), types.UpdatePasswordInput{
		NewPassword:     "12345678",
		CurrentPassword: password,
	})
	if !upd || err != nil {
		t.Fatalf("should update the user's password")
	}

	upd, err = as.UpdatePassword(fmt.Sprint(u.ID), types.UpdatePasswordInput{
		NewPassword:     "12345678",
		CurrentPassword: "123",
	})
	if upd || err == nil {
		t.Fatalf("should not update the user's password")
	}
}

func TestForgotPassword(t *testing.T) {
	defer tests.Teardown()

	as.Signup(types.SignupInput{Email: email, Password: password})

	_, err := as.ForgotPassword(types.ForgotPasswordInput{
		Email: email,
	})
	if err != nil {
		t.Fatalf("should create and return forgot password token")
	}

	_, err = as.ForgotPassword(types.ForgotPasswordInput{
		Email: "random@testing.com",
	})
	if err == nil {
		t.Fatalf("should not create and return forgot password token")
	}
}

func TestResetPassword(t *testing.T) {
	defer tests.Teardown()

	as.Signup(types.SignupInput{Email: email, Password: password})
	token, _ := as.ForgotPassword(types.ForgotPasswordInput{
		Email: email,
	})

	done, err := as.ResetPassword(token, "1234567")
	if !done || err != nil {
		t.Fatalf("should reset the user's password given the right token")
	}

	done, err = as.ResetPassword("123", "12345678")
	if done || err == nil {
		t.Fatalf("should not reset the user's password because of bad token")
	}
}
