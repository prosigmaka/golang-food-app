package handler

import (
	"errors"
	"food-app/internal/domain/entity"
	"food-app/internal/domain/service"
	"food-app/pkg/jwttoken"
	"food-app/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	var userHandler = UserHandler{}
	userHandler.userService = userService
	return &userHandler
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var registerUser entity.ReqisterViewModel
	err := c.ShouldBindJSON(&registerUser)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	registerUserError := registerUser.Validate()
	if len(registerUserError) > 0 {
		response.ResponseCustomError(c, registerUserError, http.StatusBadRequest)
		return
	}

	result, err := h.userService.SaveUser(&registerUser)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ResponseCreated(c, result)
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	result, err := h.userService.GetListUser()
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if result == nil {
		result = &[]entity.UserViewModel{}
	}

	response.ResponseOKWithData(c, result)
}

func (h *UserHandler) GetDetailUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.userService.GetDetailUser(userId)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if result == nil {
		result = &entity.UserViewModel{}
	}

	response.ResponseOKWithData(c, result)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response.ResponseError(c, errors.New("Invalid User ID").Error(), http.StatusBadRequest)
		return
	}

	var updateUser entity.User
	err = c.ShouldBindJSON(&updateUser)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	updateUser.ID = userId

	updateUserError := updateUser.Validate()
	if len(updateUserError) > 0 {
		response.ResponseCustomError(c, updateUserError, http.StatusBadRequest)
		return
	}

	result, err := h.userService.UpdateUser(&updateUser)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if result == nil {
		result = &entity.UserViewModel{}
	}

	response.ResponseOKWithData(c, result)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response.ResponseError(c, errors.New("Invalid User ID").Error(), http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(userId)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ResponseOK(c, "Successfully User Deleted")
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginVM entity.LoginViewModel

	err := c.ShouldBindJSON(&loginVM)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	validateUser, err := h.userService.GetUserByEmailPassword(loginVM)
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if validateUser == nil {
		validateUser = &entity.User{}
	}

	// Generete JWT
	token, err := jwttoken.CreateToken(int64(validateUser.ID))
	if err != nil {
		response.ResponseError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	userData := map[string]interface{}{
		"access_token": token.AccessToken,
		"expired":      token.ExpiredToken,
		"user_id":      validateUser.ID,
	}

	response.ResponseOKWithData(c, userData)
}
