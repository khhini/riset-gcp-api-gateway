package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khhini/riset-gcp-api-gateway-auth/config"
	"github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/entity"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	cfg := config.DefaultConfig("test", "0.1")
	cfg.Listen.LoadFromEnv()
	cfg.DBConfig.LoadFromEnv()

	pool, err := pgxpool.New(context.Background(), cfg.DBConfig.ConnStr())
	if err != nil {
		t.Fatalf("%v", err)
		return nil, nil
	}

	return pool, func() {
		if pool != nil {
			pool.Close()
		}
	}
}

func setupTestData(t *testing.T, repo entity.Repository) (string, string) {
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
	userID, err := repo.InsertUser(context.Background(), &userToInsert)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}

	return userID.String(), userToInsert.Username
}

func cleanUpTestData(t *testing.T, repo entity.Repository, userID string) {
	err := repo.DeleteUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("failed to create test data: %v", err)
	}
}

func TestPostgresRepositoryInsertUser(t *testing.T) {
	//Setup Test
	pool, _ := setupTestDB(t)
	repo := NewPostgreRepository(pool)

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
	userID, err := repo.InsertUser(context.Background(), &userToInsert)
	assert.Nil(t, err)
	assert.NotEmpty(t, userID)

	cleanUpTestData(t, repo, userID.String())
}

func TestPostgresRepositoryGetUserByID(t *testing.T) {
	// Setup Test
	pool, _ := setupTestDB(t)
	repo := NewPostgreRepository(pool)
	testUserID, _ := setupTestData(t, repo)
	var (
		user *entity.User
		err  error
	)

	// Test GetUserByID with Valid ID
	user, err = repo.GetUserByID(context.Background(), testUserID)
	assert.Nil(t, err)
	assert.NotNil(t, user)

	// Test GetUserByID with unknown ID
	user, err = repo.GetUserByID(context.Background(), "cf253ef2-83c0-11ee-b962-0242ac120002")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")

	// Test GetUserByID with invalid ID
	user, err = repo.GetUserByID(context.Background(), "invalid_user_id")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "error fetching user")

	cleanUpTestData(t, repo, testUserID)
}

func TestPostgresRepositoryGetUserByUsername(t *testing.T) {
	// Setup Test
	pool, _ := setupTestDB(t)
	repo := NewPostgreRepository(pool)
	testUserID, testUsername := setupTestData(t, repo)
	var (
		user *entity.User
		err  error
	)

	// Test GetUserByID with registered unsername
	user, err = repo.GetUserByUsername(context.Background(), testUsername)
	assert.Nil(t, err)
	assert.NotNil(t, user)

	// Test GetUserByID with unregistered username
	user, err = repo.GetUserByUsername(context.Background(), "unregistered_user")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")

	cleanUpTestData(t, repo, testUserID)
}

func TestPostgresRepositoryUpdateUser(t *testing.T) {
	// Setup Test
	pool, _ := setupTestDB(t)
	repo := NewPostgreRepository(pool)
	testUserID, _ := setupTestData(t, repo)
	var (
		updatedUser entity.User
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
	err = repo.UpdateUser(context.Background(), testUserID, &updatedUser)
	assert.Nil(t, err)

	// Verify updated user data
	updatedUserFromBD, err := repo.GetUserByID(context.Background(), testUserID)
	assert.Nil(t, err)
	assert.NotNil(t, updatedUserFromBD)
	assert.Equal(t, updatedUser.Username, updatedUserFromBD.Username)
	assert.Equal(t, updatedUser.Email, updatedUserFromBD.Email)
	assert.Equal(t, updatedUser.LastLogin.Time.Unix(), updatedUserFromBD.LastLogin.Time.Unix())

	cleanUpTestData(t, repo, testUserID)
}

func TestPostgresRepositoryDeleteUser(t *testing.T) {
	// Setup Test
	pool, _ := setupTestDB(t)
	repo := NewPostgreRepository(pool)
	testUserID, _ := setupTestData(t, repo)
	var err error

	// Test DeleteUser
	err = repo.DeleteUser(context.Background(), testUserID)
	assert.Nil(t, err)

	// Verify user data deleted
	deletedUser, err := repo.GetUserByID(context.Background(), testUserID)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, deletedUser)
}
