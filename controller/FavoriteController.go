package controller

import (
	"net/http"
	"time"

	"realworld-go-nolambda/util"
	"realworld-go-nolambda/service"
	"realworld-go-nolambda/model"
	"github.com/gorilla/mux"
)

type PResponse struct {
	Profile ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	Username  string `json:"username"`
	Image     string `json:"image"`
	Bio       string `json:"bio"`
	Following bool   `json:"following"`
}

func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	vars := mux.Vars(r)
	slug := vars["slug"]
	articleId, err := model.SlugToArticleId(slug)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	favoriteArticleKey := model.FavoriteArticleKey{
		Username:  user.Username,
		ArticleId: articleId,
	}

	err = service.UnfavoriteArticle(favoriteArticleKey)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	article, err := service.GetArticleByArticleId(articleId)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	isFavorited, authors, following, err := service.GetArticleRelatedProperties(user, []model.Article{article}, true)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	response := A1Response{
		Article: ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[0],
			FavoritesCount: article.FavoritesCount,
			Author: AuthorResponse{
				Username:  authors[0].Username,
				Bio:       authors[0].Bio,
				Image:     authors[0].Image,
				Following: following[0],
			},
		},
	}

	util.NewSuccessResponse(response, w, r)
}

func PostFavorite(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	vars := mux.Vars(r)
	slug := vars["slug"]
	articleId, err := model.SlugToArticleId(slug)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	favoriteArticle := model.FavoriteArticle{
		FavoriteArticleKey: model.FavoriteArticleKey{
			Username:  user.Username,
			ArticleId: articleId,
		},
		FavoritedAt: time.Now().UTC().UnixNano(),
	}

	err = service.SetFavoriteArticle(favoriteArticle)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	article, err := service.GetArticleByArticleId(articleId)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	isFavorited, authors, following, err := service.GetArticleRelatedProperties(user, []model.Article{article}, true)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	response := A1Response{
		Article: ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[0],
			FavoritesCount: article.FavoritesCount,
			Author: AuthorResponse{
				Username:  authors[0].Username,
				Bio:       authors[0].Bio,
				Image:     authors[0].Image,
				Following: following[0],
			},
		},
	}

	util.NewSuccessResponse(response, w, r)
}

func DeleteProfileFollow(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	vars := mux.Vars(r)
	username := vars["username"]
	publisher, err := service.GetUserByUsername(username)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
	}

	err = service.Unfollow(user.Username, publisher.Username)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	response := PResponse{
		Profile: ProfileResponse{
			Username:  publisher.Username,
			Image:     publisher.Image,
			Bio:       publisher.Bio,
			Following: false,
		},
	}

	util.NewSuccessResponse(response, w, r)
}

func PostProfileFollow(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	vars := mux.Vars(r)
	username := vars["username"]
	publisher, err := service.GetUserByUsername(username)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	err = service.Follow(user.Username, publisher.Username)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	response := PResponse{
		Profile: ProfileResponse{
			Username:  publisher.Username,
			Image:     publisher.Image,
			Bio:       publisher.Bio,
			Following: true,
		},
	}

	util.NewSuccessResponse(response, w, r)
}

func GetProfiles(w http.ResponseWriter, r *http.Request) {
	user, _, _ := service.GetCurrentUser(r.Header.Get("Authorization"))

	vars := mux.Vars(r)
	username := vars["username"]
	publisher, err := service.GetUserByUsername(username)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	following, err := service.IsFollowing(user, []string{publisher.Username})
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	response := PResponse{
		Profile: ProfileResponse{
			Username:  publisher.Username,
			Image:     publisher.Image,
			Bio:       publisher.Bio,
			Following: following[0],
		},
	}

	util.NewSuccessResponse(response, w, r)
}

// func main() {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/profiles-get/{username}", Handle).Methods("GET")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }
