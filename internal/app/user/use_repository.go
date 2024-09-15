package user

import (
	"errors"
	"gorm.io/gorm"
	"log/slog"
)

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	UpdateUser(userID uint, user *User) (*User, error)
	CreateUser(user *User) (*User, error)
	DeactivateUser(userID uint) (*User, error)
	DeleteUserByEmail(email string) error
}

type userRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func (u userRepository) CreateUser(user *User) (*User, error) {
	if err := u.DB.Create(user).Error; err != nil {
		u.logger.Error("Error creating user", err)
		return nil, err
	}
	u.logger.Info("User created successfully", "userID", user.ID)
	return user, nil
}

func (u userRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Info("User not found", "email", email)
			return nil, nil
		}
		u.logger.Error("Error fetching user by email", err)
		return nil, err
	}
	u.logger.Info("User fetched successfully", "userID", user.ID)
	return &user, nil
}

func (u userRepository) UpdateUser(userID uint, user *User) (*User, error) {
	if err := u.DB.Model(&User{}).Where("id = ?", userID).Updates(user).Error; err != nil {
		u.logger.Error("Error updating user", err)
		return nil, err
	}
	u.logger.Info("User updated successfully", "userID", userID)
	return user, nil
}

func (u userRepository) DeactivateUser(userID uint) (*User, error) {
	user, err := u.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.IsActive = false
	return u.UpdateUser(userID, user)
}

func (u userRepository) DeleteUserByEmail(email string) error {
	if err := u.DB.Where("email = ?", email).Delete(&User{}).Error; err != nil {
		u.logger.Error("Error deleting user by email", err)
		return err
	}
	u.logger.Info("User deleted successfully", "email", email)
	return nil
}

func (u userRepository) GetUserByID(userID uint) (*User, error) {
	var user User
	if err := u.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Info("User not found", "userID", userID)
			return nil, nil
		}
		u.logger.Error("Error fetching user by ID", err)
		return nil, err
	}
	u.logger.Info("User fetched successfully", "userID", userID)
	return &user, nil
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return userRepository{
		DB:     db,
		logger: logger,
	}
}
