package user

import (
	"log/slog"
	_ "mamlaka/internal/app/common"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Login(c echo.Context) error
	Register(c echo.Context) error
	RefreshToken(c echo.Context) error
}

type userHandler struct {
	logger      *slog.Logger
	userService UserService
}

// RefreshToken godoc
// @Summary Refresh the user's authentication token
// @Description Generates a new authentication token using the refresh token
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   user body RefreshToken true "Update Profile Request"
// @Success 200 {object} common.BaseResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router  /users/refresh-token [post]
func (u userHandler) RefreshToken(c echo.Context) error {
	return u.userService.RefreshToken(c)
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   SignInRequest body SignInRequest true "Sign In Request"
// @Success 200 {object} common.BaseResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router  /auth/login [post]
func (u userHandler) Login(c echo.Context) error {
	return u.userService.LoginUser(c)
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user account and sends a verification email
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   SignUpRequest body SignUpRequest true "Sign Up Request"
// @Success 201 {object} common.BaseResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router  /auth/register [post]
func (u userHandler) Register(c echo.Context) error {
	return u.userService.RegisterUser(c)
}

func NewUserHandler(logger *slog.Logger, service UserService) UserHandler {
	return userHandler{logger: logger, userService: service}
}
