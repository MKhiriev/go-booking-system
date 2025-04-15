package repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"humoBooking/internal/database/repositories"
	"humoBooking/internal/models"
	"regexp"
	"testing"
	"time"
)

func TestUserRepository_Create(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	userToCreate := models.User{
		Name:      "Ahmad",
		Email:     "ahmad@example.com",
		Telephone: "+992989991745",
	}

	rows := sqlmock.NewRows([]string{"user_id", "name", "email", "telephone", "role_id", "active", "created_at"}).
		AddRow(1, userToCreate.Name, userToCreate.Email, userToCreate.Telephone, userToCreate.RoleId, true, time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		// column order is defined by struct's fields order
		`INSERT INTO "users" ("name","email","telephone","role_id","active","created_at") VALUES ($1,$2,$3,$4,$5,$6)`,
	)).
		WithArgs(userToCreate.Name, userToCreate.Email, userToCreate.Telephone, userToCreate.RoleId, userToCreate.Active, NotNullTimeArg()).
		WillReturnRows(rows)
	mock.ExpectCommit()

	// 2. Act
	updatedUser, err := repo.Create(userToCreate)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, updatedUser.UserId)
	assert.Equal(t, 1, updatedUser.UserId)
	assert.Equal(t, userToCreate.Name, updatedUser.Name)
	assert.Equal(t, userToCreate.Email, updatedUser.Email)
	assert.Equal(t, userToCreate.Telephone, updatedUser.Telephone)
	assert.Equal(t, userToCreate.RoleId, updatedUser.RoleId)
	assert.Equal(t, true, updatedUser.Active)
	assert.Equal(t, false, updatedUser.CreatedAt.IsZero())
}

func TestUserRepository_GetAll(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	expectedUser1 := models.User{
		UserId:    1,
		RoleId:    1,
		Name:      "Ahmad",
		Email:     "ahmad@example.com",
		Telephone: "+992989991745",
		Active:    true,
	}
	expectedUser2 := models.User{
		UserId:    2,
		RoleId:    2,
		Name:      "Jamshed",
		Email:     "JamshedS2@example.com",
		Telephone: "+1388",
		Active:    true,
	}

	rows := sqlmock.NewRows([]string{"user_id", "name", "role_id", "email", "telephone", "active"}).
		AddRow(expectedUser1.UserId, expectedUser1.Name, expectedUser1.RoleId, expectedUser1.Email, expectedUser1.Telephone, expectedUser1.Active).
		AddRow(expectedUser2.UserId, expectedUser2.Name, expectedUser2.RoleId, expectedUser2.Email, expectedUser2.Telephone, expectedUser2.Active)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(rows)

	allUsers := repo.GetAll()

	assert.Equal(t, expectedUser1.UserId, allUsers[0].UserId)
	assert.Equal(t, expectedUser1.Name, allUsers[0].Name)
	assert.Equal(t, expectedUser1.RoleId, allUsers[0].RoleId)
	assert.Equal(t, expectedUser1.Email, allUsers[0].Email)
	assert.Equal(t, expectedUser1.Telephone, allUsers[0].Telephone)
	assert.Equal(t, expectedUser1.Active, allUsers[0].Active)

	assert.Equal(t, expectedUser2.UserId, allUsers[1].UserId)
	assert.Equal(t, expectedUser2.Name, allUsers[1].Name)
	assert.Equal(t, expectedUser2.RoleId, allUsers[1].RoleId)
	assert.Equal(t, expectedUser2.Email, allUsers[1].Email)
	assert.Equal(t, expectedUser2.Telephone, allUsers[1].Telephone)
	assert.Equal(t, expectedUser2.Active, allUsers[1].Active)
}

func TestUserRepository_GetUserByID(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	expectedUser := models.User{
		UserId:    1,
		RoleId:    1,
		Name:      "Ahmad",
		Email:     "ahmad@example.com",
		Telephone: "+992989991745",
		Active:    true,
	}

	rows := sqlmock.NewRows([]string{"user_id", "name", "role_id", "email", "telephone", "active"}).
		AddRow(expectedUser.UserId, expectedUser.Name, expectedUser.RoleId, expectedUser.Email, expectedUser.Telephone, expectedUser.Active)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "user_id" = $1`)).
		WithArgs(expectedUser.UserId).
		WillReturnRows(rows)

	// 2. Act
	user, err := repo.GetUserById(1)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.UserId, user.UserId)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.RoleId, user.RoleId)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Telephone, user.Telephone)
	assert.Equal(t, expectedUser.Active, user.Active)
}

func TestUserRepository_Update(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	user := models.User{
		UserId:    1,
		RoleId:    1,
		Name:      "Ahmad",
		Email:     "ahmad@example.com",
		Telephone: "+992989991745",
		Active:    true,
		CreatedAt: time.Now(),
	}
	userId := user.UserId
	newRoleId := 2
	newName := "Ahmad Jr Developer"
	newEmail := "ahmad@alif.tj"
	newTelephone := "+992989992222"
	updateData := models.User{
		UserId:    user.UserId,
		RoleId:    newRoleId,
		Name:      newName,
		Email:     newEmail,
		Telephone: newTelephone,
	}

	// check if transaction started
	mock.ExpectBegin()
	// check if GORM sql-query matches our pattern
	mock.ExpectExec(regexp.QuoteMeta(
		// column order is defined by struct's fields order
		`UPDATE "users" SET "name"=$1,"email"=$2,"telephone"=$3,"role_id"=$4,"updated_at"=$5  WHERE "user_id" = $6`,
	)).
		// WithArgs(updateData.Name, updateData.Email, updateData.Telephone, updateData.RoleId, sqlmock.AnyArg(), userId).
		WithArgs(updateData.Name, updateData.Email, updateData.Telephone, updateData.RoleId, NotNullTimeArg(), userId).
		// will return 1 row that is going to BE UPDATED
		WillReturnResult(sqlmock.NewResult(0, 1))
	// check if COMMIT
	mock.ExpectCommit()

	// 2. Act
	updatedUser, err := repo.Update(updateData)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, time.Time{}, updatedUser.UpdatedAt)
	assert.Equal(t, userId, updatedUser.UserId)
	assert.Equal(t, newRoleId, updatedUser.RoleId)
	assert.Equal(t, newName, updatedUser.Name)
	assert.Equal(t, newEmail, updatedUser.Email)
	assert.Equal(t, newTelephone, updatedUser.Telephone)
}

func TestUserRepository_Delete(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	userToDelete := models.User{
		UserId:    1,
		RoleId:    1,
		Name:      "Ahmad",
		Email:     "ahmad@example.com",
		Telephone: "+992989991745",
		Active:    true,
		CreatedAt: time.Now(),
	}
	userId := userToDelete.UserId

	// check if transaction started
	mock.ExpectBegin()
	// check if GORM sql-query matches our pattern
	mock.ExpectExec(regexp.QuoteMeta(
		// column order is defined by struct's fields order
		`UPDATE "users" SET "active"=$1,"deleted_at"=$2 WHERE "user_id" = $3`,
	)).
		// where args are ["active": false, "deleted_at": generated current_timestamp, "user_id": user_id]
		WithArgs(false, NotNullTimeArg(), userId).
		// will return count: 1 row that is going to BE UPDATED
		WillReturnResult(sqlmock.NewResult(0, 1))
	// check if COMMIT
	mock.ExpectCommit()

	// 2. Act
	isDeleted, err := repo.Delete(userId)

	// 3. Assert
	assert.NoError(t, err)
	assert.Equal(t, true, isDeleted)
}

func TestBookingRepository_UpdateUserPassword(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	userId := 3
	newPassword := "password"
	updateData := models.User{
		UserId:   userId,
		Password: newPassword,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "password"=$1,"updated_at"=$2 WHERE "user_id" = $3`,
	)).
		WithArgs(updateData.Password, NotNullTimeArg(), updateData.UserId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 2. Act
	updatedUser, err := repo.UpdatePassword(updateData)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, updatedUser.UserId)
	assert.NotEqual(t, true, updatedUser.UpdatedAt.IsZero())
	assert.Equal(t, userId, updatedUser.UserId)
	assert.Equal(t, newPassword, updatedUser.Password)
}

func TestBookingRepository_UpdateUserName(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	userId := 3
	newUsername := "killfish"
	updateData := models.User{
		UserId:   userId,
		UserName: newUsername,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "username"=$1,"updated_at"=$2 WHERE "user_id" = $3`,
	)).
		WithArgs(updateData.UserName, NotNullTimeArg(), updateData.UserId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 2. Act
	updatedUser, err := repo.UpdateUsername(updateData)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, updatedUser.UserId)
	assert.NotEqual(t, true, updatedUser.UpdatedAt.IsZero())
	assert.Equal(t, userId, updatedUser.UserId)
	assert.Equal(t, newUsername, updatedUser.UserName)
}

func TestBookingRepository_UpdateUserRole(t *testing.T) {
	// 1. Assess
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepositoryPostgres(db)

	userId := 3
	newUserRole := 3
	updateData := models.User{
		UserId: userId,
		RoleId: newUserRole,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "users" SET "role_id"=$1,"updated_at"=$2 WHERE "user_id" = $3`,
	)).
		WithArgs(updateData.RoleId, NotNullTimeArg(), updateData.UserId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 2. Act
	updatedUser, err := repo.UpdateUserRole(updateData)

	// 3. Assert
	assert.NoError(t, err)
	assert.NotEqual(t, 0, updatedUser.UserId)
	assert.NotEqual(t, true, updatedUser.UpdatedAt.IsZero())
	assert.Equal(t, userId, updatedUser.UserId)
	assert.Equal(t, newUserRole, updatedUser.RoleId)
}
