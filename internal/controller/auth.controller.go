package controller

import (
	"net/http"
	"strings"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/response"
	"github.com/ariboss89/coffee-morning-services/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user with email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	dto.LoginResponse
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		401		{object}	dto.ResponseError
//	@Router			/auth/login [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Email field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			response.Error(ctx, http.StatusBadRequest, "Email must be a valid email address")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	data, err := ac.authService.Login(ctx, req)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "email must be a valid email address" {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	response.Success(ctx, http.StatusOK, "Login successful", dto.JWT{Token: data.Data.Token})
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	Create a new user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"User registration data"
//	@Success		201		{object}	dto.RegisterResponse
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		errStr := err.Error()

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Email field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Email") && strings.Contains(errStr, "email") {
			response.Error(ctx, http.StatusBadRequest, "Email must be a valid email address")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "required") {
			response.Error(ctx, http.StatusBadRequest, "Password field cannot be empty")
			return
		}

		if strings.Contains(errStr, "Password") && strings.Contains(errStr, "min") {
			response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	data, err := ac.authService.Register(ctx, req)

	if err != nil {
		if strings.Contains(err.Error(), "user exist") {
			response.Error(ctx, http.StatusBadRequest, "Email already in use")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response.Success(ctx, http.StatusCreated, "Registered successfully", data)
}

// Logout User godoc
// @Summary      Logout user
// @Tags         Auth
// @Produce      json
// @Success      201  {object}  dto.ResponseSuccess
// @Failure 		 400 {object} dto.ResponseError
// @Failure 		 500 {object} dto.ResponseError
// @Failure 		 403 {object} dto.ResponseError
// @Failure 		 401 {object} dto.ResponseError
// @Router       /auth/logout [delete]
// @security		 BearerAuth
func (a AuthController) Logout(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	tokenJWT := strings.Split(auth, "Bearer ")

	token := tokenJWT[1]

	// if err := a.authService.LogoutUser(c.Request.Context(), token); err != nil {
	_, err := a.authService.LogoutUser(c.Request.Context(), token)

	if err != nil {
		str := err.Error()
		if strings.Contains(str, "token obsoleted") {
			response.Error(c, http.StatusBadRequest, "Token Obsoleted")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response.Success(c, http.StatusCreated, "Logout successfully", nil)
}
