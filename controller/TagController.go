package controller

import (
	"net/http"
	"realworld-go-nolambda/util"
	"realworld-go-nolambda/service"
)

type TResponse struct {
	Tags []string `json:"tags"`
}

func GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := service.GetTags()
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	response := TResponse{
		Tags: tags,
	}

	util.NewSuccessResponse(response, w, r)
}
