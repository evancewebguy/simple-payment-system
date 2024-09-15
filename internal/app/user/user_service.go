package user

import (
	"errors"
	"log/slog"
	"mamlaka/internal/app/common"
	"mamlaka/internal/pkg/auth"
	"mamlaka/internal/pkg/tokens"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// UserService defines the methods available in the user service.
type UserService interface {
	LoginUser(c echo.Context) error
	RegisterUser(c echo.Context) error
	RefreshToken(c echo.Context) error
}

// userService is the implementation of UserService.
type userService struct {
	logger     *slog.Logger
	repository UserRepository
}

// LoginUser handles user login requests by validating credentials, checking user status, and generating tokens.
func (u userService) LoginUser(c echo.Context) error {
	// Retrieve and parse the request body into a LoginRequest structure
	var loginRequest SignInRequest
	if err := c.Bind(&loginRequest); err != nil {
		u.logger.Error("Error parsing login request body", err)
		return u.handleError(c, err, http.StatusBadRequest)
	}

	// Validate the login request body
	if err := common.ValidateModel(loginRequest); err != nil {
		u.logger.Error("Invalid login request body", err)
		return u.handleError(c, err, http.StatusBadRequest)
	}

	// Determine the provider from the request
	provider := loginRequest.Provider

	var user *User
	var err error

	// Fetch user details based on the provider
	switch provider {
	case ProviderEmail:
		if loginRequest.Email == "" {
			return u.handleError(c, errors.New("email address is required"), http.StatusBadRequest)
		}
		if loginRequest.Password == "" {
			return u.handleError(c, errors.New("password is required"), http.StatusBadRequest)
		}
		user, err = u.repository.GetUserByEmail(loginRequest.Email)
	default:
		if loginRequest.Email == "" {
			return u.handleError(c, errors.New("email address is required"), http.StatusBadRequest)
		}
		if loginRequest.Password == "" {
			return u.handleError(c, errors.New("password is required"), http.StatusBadRequest)
		}
		user, err = u.repository.GetUserByEmail(loginRequest.Email)
	}

	if err != nil {
		u.logger.Error("Error retrieving user from repository", err)
		return u.handleError(c, err, http.StatusInternalServerError)
	}

	if user == nil {
		return u.handleError(c, errors.New("user with provided credentials not found"), http.StatusNotFound)
	}

	// Verify password
	if provider == ProviderEmail {
		if !auth.CheckPasswordHash(loginRequest.Password, user.Password) {
			return u.handleError(c, errors.New("invalid credentials"), http.StatusUnauthorized)
		}
	}

	// Check if the user account is active and verified
	if !user.IsActive {
		return u.handleError(c, errors.New("account is not active"), http.StatusUnauthorized)
	}

	if !user.IsVerified {
		return u.handleError(c, errors.New("account is not verified"), http.StatusUnauthorized)
	}

	// Generate access and refresh tokens
	accessToken, err := tokens.GenerateAccessToken(strconv.Itoa(int(user.ID)), time.Minute*60)
	if err != nil {
		u.logger.Error("Error generating access tokens", err)
		return u.handleError(c, err, http.StatusInternalServerError)
	}

	// Generate refresh tokens
	refreshToken, err := tokens.GenerateRefreshToken(strconv.Itoa(int(user.ID)), time.Minute*60)
	if err != nil {
		u.logger.Error("Error generating refresh tokens", err)
		return u.handleError(c, err, http.StatusInternalServerError)
	}

	// Return user data and tokens
	response := LoginResponse{
		User: UserDto{
			ID:          user.ID,
			FullName:    user.FullName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			IsActive:    user.IsActive,
			IsVerified:  user.IsVerified,
		},
		Token: RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	_, err = u.repository.UpdateUser(user.ID, user)
	if err != nil {
		u.logger.Error("Error updating user", err)
		return u.handleError(c, err, http.StatusInternalServerError)
	}

	u.logger.Info("User login successful", user)
	return c.JSON(http.StatusOK, common.BaseResponse{
		Status:  http.StatusOK,
		Message: "Login successful",
		Data:    response,
	})
}

// handleError is a helper function for creating error responses.
func (u userService) handleError(c echo.Context, err error, status int) error {
	return c.JSON(status, common.ErrorResponse{
		Status: status,
		Error:  err.Error(),
	})
}

// getProfileDetails handles fetching and processing user profile details.
func (u userService) getProfileDetails(provider ProviderType, token string) (*User, error) {
	var userEmail, fullName string
	switch provider {
	case ProviderEmail:
		userEmail = ""
	default:
		return nil, errors.New("unsupported provider")
	}

	return &User{
		FullName:   fullName,
		Email:      userEmail,
		IsActive:   false,
		IsVerified: false,
	}, nil
}

// RegisterUser handles user registration.
func (u userService) RegisterUser(c echo.Context) error {
	var signUpRequest SignUpRequest

	if err := c.Bind(&signUpRequest); err != nil {
		u.logger.Error("Error parsing request body", err)
		return u.handleError(c, err, http.StatusBadRequest)
	}

	if err := common.ValidateModel(signUpRequest); err != nil {
		u.logger.Error("Error validating request body", err)
		return u.handleError(c, err, http.StatusBadRequest)
	}

	var user *User
	if signUpRequest.Provider != ProviderEmail {
		if signUpRequest.Token == "" {
			return u.handleError(c, errors.New("token is required"), http.StatusBadRequest)
		}

		var err error
		user, err = u.getProfileDetails(signUpRequest.Provider, signUpRequest.Token)
		if err != nil {
			u.logger.Error("Error getting profile details", err)
			return u.handleError(c, err, http.StatusBadRequest)
		}
	} else {
		if signUpRequest.Email == "" {
			return u.handleError(c, errors.New("email is required"), http.StatusBadRequest)
		}
		if signUpRequest.FullName == "" {
			return u.handleError(c, errors.New("full name is required"), http.StatusBadRequest)
		}
		if signUpRequest.Password == "" {
			return u.handleError(c, errors.New("password is required"), http.StatusBadRequest)
		}
		hash, err := auth.HashPassword(signUpRequest.Password)
		if err != nil {
			u.logger.Error("Error hashing password", err)
			return u.handleError(c, err, http.StatusInternalServerError)
		}
		user = &User{
			FullName:   signUpRequest.FullName,
			Email:      signUpRequest.Email,
			Password:   hash,
			IsActive:   false,
			IsVerified: false,
		}
	}

	if user, _ := u.repository.GetUserByEmail(user.Email); user != nil {
		return u.handleError(c, errors.New("user with provided email already exists"), http.StatusBadRequest)
	}

	if _, err := u.repository.CreateUser(user); err != nil {
		u.logger.Error("Error creating user", err)
		return u.handleError(c, err, http.StatusInternalServerError)
	}

	u.logger.Info("User account created successfully", user)
	return c.JSON(http.StatusCreated, common.BaseResponse{
		Status:  http.StatusCreated,
		Message: "User account successfully created",
		Data:    user,
	})
}

func (u userService) RefreshToken(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.BaseResponse{
		Status:  http.StatusNotImplemented,
		Message: "Not Implemented",
	})
}

func (u userService) GetUserAccountProfile(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.BaseResponse{
		Status:  http.StatusNotImplemented,
		Message: "Not Implemented",
	})
}

func (u userService) UpdateUserAccountProfile(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.BaseResponse{
		Status:  http.StatusNotImplemented,
		Message: "Not Implemented",
	})
}

func (u userService) DeleteUserAccountProfile(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.BaseResponse{
		Status:  http.StatusNotImplemented,
		Message: "Not Implemented",
	})
}

func (u userService) UpdateUserAccountPreferences(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.BaseResponse{
		Status:  http.StatusNotImplemented,
		Message: "Not Implemented",
	})
}

func (u userService) CreateUserAccountPreferences(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, common.BaseResponse{
		Status:  http.StatusNotImplemented,
		Message: "Not Implemented",
	})
}

// NewUserService creates a new instance of userService.
func NewUserService(logger *slog.Logger, repository UserRepository) UserService {
	return userService{
		logger:     logger,
		repository: repository,
	}
}
