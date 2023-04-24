package controller

import (
	"bytes"
	"net/http"

	"realworld-go-nolambda/util"
	"realworld-go-nolambda/service"
	"realworld-go-nolambda/model"
)

type UResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
	Token    string `json:"token"`
}
type ULRequest struct {
	User UserLoginRequest `json:"user"`
}
type UPuRequest struct {
	User UserPutRequest `json:"user"`
}
type UPoRequest struct {
	User UserPostRequest `json:"user"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPostRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPutRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, token, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	response := UResponse{
		User: UserResponse{
			Username: user.Username,
			Email:    user.Email,
			Image:    user.Image,
			Bio:      user.Bio,
			Token:    token,
		},
	}

	util.NewSuccessResponse(response, w, r)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	request := &ULRequest{}
	err := util.ParseBody(r, request)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	user, err := service.GetUserByEmail(request.User.Email)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	passwordHash, err := model.Scrypt(request.User.Password)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	if !bytes.Equal(passwordHash, user.PasswordHash) {
		util.NewErrorResponse(http.StatusBadRequest, model.NewInputError("password", "wrong password"), w)
		return
	}

	token, err := model.GenerateToken(user.Username)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	response := UResponse{
		User: UserResponse{
			Username: user.Username,
			Email:    user.Email,
			Image:    user.Image,
			Bio:      user.Bio,
			Token:    token,
		},
	}

	util.NewSuccessResponse(response, w, r)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	request := &UPoRequest{}
	err := util.ParseBody(r, request)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	err = model.ValidatePassword(request.User.Password)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	passwordHash, err := model.Scrypt(request.User.Password)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	user := model.User{
		Username:     request.User.Username,
		Email:        request.User.Email,
		PasswordHash: passwordHash,
	}

	err = service.PutUser(user)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	token, err := model.GenerateToken(user.Username)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	response := UResponse{
		User: UserResponse{
			Username: user.Username,
			Email:    user.Email,
			Image:    user.Image,
			Bio:      user.Bio,
			Token:    token,
		},
	}

	util.NewSuccessResponse(response, w, r)
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	oldUser, token, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
		return
	}

	request := &UPuRequest{}
	err = util.ParseBody(r, request)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	err = model.ValidatePassword(request.User.Password)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	passwordHash, err := model.Scrypt(request.User.Password)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	newUser := model.User{
		Username:     oldUser.Username,
		Email:        request.User.Email,
		PasswordHash: passwordHash,
		Image:        request.User.Image,
		Bio:          request.User.Bio,
	}

	err = service.UpdateUser(*oldUser, newUser)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	response := UResponse{
		User: UserResponse{
			Username: newUser.Username,
			Email:    newUser.Email,
			Image:    newUser.Image,
			Bio:      newUser.Bio,
			Token:    token,
		},
	}

	util.NewSuccessResponse(response, w, r)
}
