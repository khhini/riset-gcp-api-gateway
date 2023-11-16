package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khhini/riset-gcp-api-gateway-auth/config"
	"github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/entity"
	"github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/repository"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) entity.Repository {
	cfg := config.DefaultConfig("test", "0.1")
	cfg.Listen.LoadFromEnv()
	cfg.DBConfig.LoadFromEnv()

	pool, err := pgxpool.New(context.Background(), cfg.DBConfig.ConnStr())
	if err != nil {
		t.Fatalf("%v", err)
		return nil
	}

	return repository.NewPostgreRepository(pool)
}

func setupTestData(t *testing.T, svc entity.Service) string {
	// Create a test user
	var userToInsert entity.User
	userJson := []byte(`{
		"username": "testuser",
		"email": "testuser@example.com",
		"password": "superdupersecretpassword" 
	}`)

	if err := json.Unmarshal(userJson, &userToInsert); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	userToInsert.CreatedOn = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// Test InsertUser
	userID, err := svc.InsertUser(context.Background(), &userToInsert)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	return userID.String()
}

func cleanUpTestData(t *testing.T, svc entity.Service, userID string) {
	err := svc.DeleteUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}
}

func TestServiceInsertUser(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewService(repo)

	// Create a test user
	var userToInsert entity.User
	userJson := []byte(`{
		"username": "testuser",
		"email": "testuser@example.com",
		"password": "superdupersecretpassword" 
	}`)

	if err := json.Unmarshal(userJson, &userToInsert); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	unencryptedPassword := userToInsert.Password

	userToInsert.CreatedOn = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// Test InsertUser
	userID, err := svc.InsertUser(context.Background(), &userToInsert)
	assert.Nil(t, err)
	assert.NotEmpty(t, userID)
	assert.NotEqual(t, unencryptedPassword, userToInsert.Password)

	cleanUpTestData(t, svc, userID.String())
}

func TestServiceGetUserByID(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewService(repo)
	testUserID := setupTestData(t, svc)

	var (
		user *entity.User
		err  error
	)

	// Test GetUserByID with Valid ID
	user, err = svc.GetUserByID(context.Background(), testUserID)
	assert.Nil(t, err)
	assert.NotNil(t, user)

	// Test GetUserByID with unknown ID
	user, err = svc.GetUserByID(context.Background(), "cf253ef2-83c0-11ee-b962-0242ac120002")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")

	// Test GetUserByID with invalid ID
	user, err = svc.GetUserByID(context.Background(), "invalid_user_id")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "error fetching user")

	cleanUpTestData(t, svc, testUserID)
}

func TestServiceUpdateUser(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewService(repo)
	testUserID := setupTestData(t, svc)
	var (
		updatedUser *entity.User
		err         error
	)

	// Update user data
	updatedJson := []byte(`{
		"username": "updasteduser",
		"email": "updateduser@example.com"
	}`)

	if err := json.Unmarshal(updatedJson, &updatedUser); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	updatedUser.LastLogin = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// Test UpdateUser
	err = svc.UpdateUser(context.Background(), testUserID, updatedUser)
	assert.Nil(t, err)

	// Verify updated user data
	updatedUserFromBD, err := svc.GetUserByID(context.Background(), testUserID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedUserFromBD)
	assert.Equal(t, updatedUser.Username, updatedUserFromBD.Username)
	assert.Equal(t, updatedUser.Email, updatedUserFromBD.Email)
	assert.Equal(t, updatedUser.LastLogin.Time.Unix(), updatedUserFromBD.LastLogin.Time.Unix())

	cleanUpTestData(t, svc, testUserID)
}

func TestServiceDeleteUser(t *testing.T) {
	repo := setupTestDB(t)
	svc := NewService(repo)
	testUserID := setupTestData(t, svc)
	var err error

	// Test DeleteUser
	err = svc.DeleteUser(context.Background(), testUserID)
	assert.Nil(t, err)

	// Verify user data deleted
	deletedUser, err := svc.GetUserByID(context.Background(), testUserID)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, deletedUser)
}
