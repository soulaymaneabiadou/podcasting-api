package repositories

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"podcast/database"
	"podcast/hasher"
	"podcast/models"
	"podcast/tests"
	"podcast/types"
)

var (
	ur           = NewUsersRepository(database.DB)
	err          error
	user         models.User
	name         string     = "Tester"
	email        string     = "test12@testing.com"
	password     string     = "123456789"
	creatorRole  types.Role = types.CREATOR_ROLE
	listenerRole types.Role = types.LISTENER_ROLE
)

func setupReposTests(t *testing.T) {
	user, err = ur.Create(types.CreateUserInput{Name: name, Email: email, Password: password, Role: listenerRole})
	if err != nil {
		t.Fatal(err.Error())
	}

	if user.ID == 0 {
		t.Fatalf("A user document should have been created")
	}

	if user.Name != name {
		t.Fatalf("The user's name mismatch, want %s, got %s", name, user.Name)
	}

	if user.Email != email {
		t.Fatalf("The user's email mismatch, want %s, got %s", email, user.Email)
	}

	if user.Role != string(listenerRole) {
		t.Fatalf("The user's role mismatch, want %s, got %s", string(listenerRole), user.Role)
	}
}

func TestGetAll(t *testing.T) {
	setupReposTests(t)
	defer tests.Teardown()

	users, err := ur.GetAll(database.Paginator{Limit: 10, Page: 1})

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(users) != 1 {
		t.Fatalf(`should find 1 user in the database, found %d`, len(users))
	}

	if users[0].ID != user.ID {
		t.Fatalf("user should include a valid id, want %d found %d", user.ID, users[0].ID)
	}

	if users[0].Email != email {
		t.Fatalf("user should include the proper email, want %s, found %s", email, user.Email)
	}
}

func TestGetById(t *testing.T) {
	setupReposTests(t)
	defer tests.Teardown()

	u, err := ur.GetById(fmt.Sprint(user.ID))
	if err != nil {
		t.Fatalf(err.Error())
	}

	if u.Email != user.Email {
		t.Fatalf("should find the user by a given id, want %s found %s", user.Email, u.Email)
	}
}

func TestGetByEmail(t *testing.T) {
	setupReposTests(t)
	defer tests.Teardown()

	u, err := ur.GetByEmail(email)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if u.ID != user.ID {
		t.Fatalf("should find the user by their email, want %d found %d", user.ID, u.ID)
	}
}

func TestGetByResetPasswordToken(t *testing.T) {
	setupReposTests(t)
	defer tests.Teardown()

	token, err := hasher.GenerateSecureToken(20)
	if err != nil {
		t.Fatalf(err.Error())
	}

	h := sha256.Sum256([]byte(token))
	hash := hex.EncodeToString(h[:])

	u, err := ur.Update(user, types.UpdateUserInput{
		ResetPasswordToken:  hash,
		ResetPasswordExpire: time.Now().Add(time.Minute * 10),
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	u, err = ur.GetByResetPasswordToken(hash)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if u.ID != user.ID {
		t.Fatalf("should find the right user given their reset token hash, want %d found %d", user.ID, u.ID)
	}
}

func TestCreateUser(t *testing.T) {
	defer tests.Teardown()

	user, err := ur.Create(types.CreateUserInput{Email: "test2@testing.com", Password: "12345678"})
	if err != nil {
		t.Fatal(err.Error())
	}

	if user.ID == 0 {
		t.Fatalf("should create the user without issues")
	}
}

func TestUpdateUser(t *testing.T) {
	setupReposTests(t)
	defer tests.Teardown()

	u, err := ur.Update(user, types.UpdateUserInput{Email: "updated@testing.com"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	if u.Email != "updated@testing.com" {
		t.Fatalf("A user document should have the updated email, want %s found %s", "updated@testing.com", u.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	setupReposTests(t)
	defer tests.Teardown()

	del, err := ur.Destroy(user)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !del {
		t.Fatalf("The user should be deleted")
	}
}
