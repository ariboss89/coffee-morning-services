package controller

import (
	"fmt"
	"log"
	"net/http"
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

type InteractionController struct {
	interactionService *service.InteractionService
}

func NewInteractionController(interactionService *service.InteractionService) *InteractionController {
	return &InteractionController{interactionService: interactionService}
}

// Posts godoc
//
//	@Summary		Post new content
//	@Description	Create a new user account
//	@Tags			Interaction
//
// @accept			 multipart/form-data
// @Produce      json
// @Param        content_file formData file true  "Upload Content Image"
// @Param        caption	formData string true  "Caption"
//
//	@Success		200		{object}	dto.RegisterResponse
//	@Success		201		{object}	dto.RegisterResponse
//	@Failure		401		{object}	dto.ResponseError
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/interaction/content [post]
//
// @Security			BearerAuth
func (ic *InteractionController) PostContent(c *gin.Context) {
	const maxSize = 2 * 800 * 600
	var postContent dto.InteractionRequest
	token, isExist := c.Get("token")
	if !isExist {
		c.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
			Message: "Forbidden Access",
			Status:  "Error",
		})
		return
	}
	accessToken, _ := token.(jwt.JWTClaims)

	if err := c.ShouldBindWith(&postContent, binding.FormMultipart); err != nil {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if postContent.ContentFile != nil {
		ext := path.Ext(postContent.ContentFile.Filename)
		re := regexp.MustCompile("^[.](jpg|png)$")
		if !re.Match([]byte(ext)) {
			response.Error(c, http.StatusBadRequest, "File have to be jpg or png")
			return
		}
		// validasi ukuran
		if postContent.ContentFile.Size > maxSize {
			response.Error(c, http.StatusBadRequest, "File maximum 2 MB")
			return
		}

		filename := fmt.Sprintf("%d_profile_%d%s", time.Now().UnixNano(), accessToken.Id, ext)
		postContent.ContentName = filename

		if e := c.SaveUploadedFile(postContent.ContentFile, filepath.Join("public", "content", filename)); e != nil {
			log.Println(e.Error())
			response.Error(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	id := accessToken.Id

	if id != 0 {
		if err := ic.interactionService.PostContent(c.Request.Context(), postContent, id); err != nil {
			response.Error(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		photoUrl := fmt.Sprintf("http://localhost:8002/static/img/content/%s", postContent.ContentName)
		if postContent.ContentName != "" {
			response.Success(c, http.StatusOK, "Content inserted !!", photoUrl)
			return
		}

		response.Success(c, http.StatusOK, "Content inserted !!", nil)
		return

	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, dto.ResponseError{
			Message: "Forbidden Access",
			Status:  "Error",
		})
		return
	}
}
