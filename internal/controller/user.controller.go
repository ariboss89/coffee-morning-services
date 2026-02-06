package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/ariboss89/coffee-morning-services/internal/response"
	"github.com/ariboss89/coffee-morning-services/internal/service"
	"github.com/ariboss89/coffee-morning-services/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Get User Profiles
// @Summary      Get User Profiles
// @Tags         user
// @Produce      json
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 404 {object} dto.ResponseError
// @Failure 		 403 {object} dto.ResponseError
// @Failure 		 500 {object} dto.ResponseError
// @Router       /user/profile [get]
// @security			BearerAuth
func (u UserController) GetUserProfileById(c *gin.Context) {
	token, isExist := c.Get("token")
	if !isExist {
		response.Error(c, http.StatusForbidden, "Forbidden Access")
		return
	}
	accessToken, _ := token.(jwt.JWTClaims)

	data, err := u.userService.GetUserProfileById(c.Request.Context(), accessToken.Id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			response.Error(c, http.StatusNotFound, "Data Not Found")
			return
		}
		response.Error(c, http.StatusForbidden, "Internal Server Error")
		return
	}
	response.Success(c, http.StatusOK, "Profile Retrieved Successfully !!", data)
}

// Update Profile
// @Summary      Update Profile
// @Tags         user
// @accept			 multipart/form-data
// @Produce      json
// @Param        fullname	formData string false  "Full Name"
// @Param        bio	formData string false  "Bio"
// @Param        avatar_file formData file false  "Update Avatar"
// @Success      200  {object}  dto.ResponseSuccess
// @Failure 		 400 {object} dto.ResponseError
// @Failure 		 403 {object} dto.ResponseError
// @Failure 		 404 {object} dto.ResponseError
// @Failure 		 500 {object} dto.ResponseError
// @Router       /user/ [patch]
// @Security			BearerAuth
func (u UserController) UpdateProfile(c *gin.Context) {
	const maxSize = 2 * 800 * 600
	var updateUser dto.UserRequest
	token, isExist := c.Get("token")
	if !isExist {
		c.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
			Message: "Forbidden Access",
			Status:  "Error",
		})
		return
	}
	accessToken, _ := token.(jwt.JWTClaims)

	if err := c.ShouldBindWith(&updateUser, binding.FormMultipart); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if updateUser.AvatarFile != nil {
		ext := path.Ext(updateUser.AvatarFile.Filename)
		re := regexp.MustCompile("^[.](jpg|png)$")
		if !re.Match([]byte(ext)) {
			response.Error(c, http.StatusBadRequest, "File have to be jpg or png")
			return
		}
		// validasi ukuran
		if updateUser.AvatarFile.Size > maxSize {
			response.Error(c, http.StatusBadRequest, "File maximum 1 MB")
			return
		}

		data, errPhoto := u.userService.GetUserProfileById(c.Request.Context(), accessToken.Id)

		if errPhoto != nil {
			// Unexpected Error
			log.Println("photo is not uploaded yet")
		}

		if data.Avatar != "" {
			prevPhoto := data.Avatar
			filePath := "public" + prevPhoto
			absPath, err := filepath.Abs(filePath)
			errPath := os.Remove(absPath)

			if errPath != nil {
				log.Println(err)
			}
		}

		filename := fmt.Sprintf("%d_profile_%d%s", time.Now().UnixNano(), accessToken.Id, ext)
		updateUser.Avatar = filename

		if e := c.SaveUploadedFile(updateUser.AvatarFile, filepath.Join("public", "profile", filename)); e != nil {
			log.Println(e.Error())
			response.Error(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	id := accessToken.Id

	if id != 0 {
		if err := u.userService.UpdateProfile(c.Request.Context(), updateUser, id); err != nil {
			response.Error(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		photoUrl := fmt.Sprintf("http://localhost:8002/static/img/profile/%s", updateUser.Avatar)
		if updateUser.Avatar != "" {
			response.Success(c, http.StatusOK, "Update successfully", photoUrl)
			return
		}

		response.Success(c, http.StatusOK, "Update successfully", nil)
		return

	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
			Message: "Forbidden Access",
			Status:  "Error",
		})
		return
	}
}
