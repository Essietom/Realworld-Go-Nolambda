package controller

import (
	"net/http"
	"realworld-go-nolambda/util"
	"realworld-go-nolambda/service"
	"realworld-go-nolambda/model"


	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CResponse struct {
	Comments []CommentResponse `json:"comments"`
}

type C1Response struct {
	Comment CommentResponse `json:"comment"`
}

type CommentResponse struct {
	Id        int64          `json:"id"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	Body      string         `json:"body"`
	Author    AuthorResponse `json:"author"`
}

type CRequest struct {
	Comment CommentRequest `json:"comment"`
}

type CommentRequest struct {
	Body string `json:"body"`
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	commentId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, model.NewInputError("id", "invalid"), w)
	}

	slug := vars["slug"]
	err = service.DeleteComment(slug, commentId, user.Username)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	util.NewSuccessResponse(nil, w, r)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	user, _, _ := service.GetCurrentUser(r.Header.Get("Authorization"))

	vars := mux.Vars(r)
	slug := vars["slug"]
	comments, err := service.GetComments(slug)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	authors, following, err := service.GetCommentRelatedProperties(user, comments)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	commentResponses := make([]CommentResponse, 0, len(comments))

	for i, comment := range comments {
		commentResponses = append(commentResponses, CommentResponse{
			Id:        comment.CommentId,
			Body:      comment.Body,
			CreatedAt: time.Unix(0, comment.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt: time.Unix(0, comment.UpdatedAt).Format(model.TimestampFormat),
			Author: AuthorResponse{
				Username:  authors[i].Username,
				Bio:       authors[i].Bio,
				Image:     authors[i].Image,
				Following: following[i],
			},
		})
	}

	response := CResponse{
		Comments: commentResponses,
	}

	util.NewSuccessResponse(response, w, r)
}

func PostComment(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	request := &CRequest{}
	err = util.ParseBody(r, request)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	vars := mux.Vars(r)
	slug := vars["slug"]
	// Make sure article exists, at least at this point
	article, err := service.GetArticleBySlug(slug)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	now := time.Now().UTC()
	nowUnixNano := now.UnixNano()
	nowStr := now.Format(model.TimestampFormat)

	comment := model.Comment{
		CommentKey: model.CommentKey{
			ArticleId: article.ArticleId,
		},
		CreatedAt: nowUnixNano,
		UpdatedAt: nowUnixNano,
		Body:      request.Comment.Body,
		Author:    user.Username,
	}

	err = service.PutComment(&comment)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	response := C1Response{
		Comment: CommentResponse{
			Id:        comment.CommentId,
			Body:      comment.Body,
			CreatedAt: nowStr,
			UpdatedAt: nowStr,
			Author: AuthorResponse{
				Username:  user.Username,
				Bio:       user.Bio,
				Image:     user.Image,
				Following: false,
			},
		},
	}

	util.NewSuccessResponse(response, w, r)
}
