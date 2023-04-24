package controller

import (
	"net/http"
	"strconv"
	"time"

	"realworld-go-nolambda/util"
	"realworld-go-nolambda/service"
	"realworld-go-nolambda/model"

	"github.com/gorilla/mux"
)

type AResponse struct {
	Articles      []ArticleResponse `json:"articles"`
	ArticlesCount int               `json:"articlesCount"`
}

type A1Response struct {
	Article ArticleResponse `json:"article"`
}

type ArticleResponse struct {
	Slug           string         `json:"slug"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Body           string         `json:"body"`
	TagList        []string       `json:"tagList"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedAt      string         `json:"updatedAt"`
	Favorited      bool           `json:"favorited"`
	FavoritesCount int64          `json:"favoritesCount"`
	Author         AuthorResponse `json:"author"`
}

type AuthorResponse struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type APuRequest struct {
	Article ArticlePutRequest `json:"article"`
}

type ArticlePutRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

func GetArticlesFeed(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewErrorResponse(http.StatusUnauthorized, err, w)
	}

	query := r.URL.Query()
	offset, err := strconv.Atoi(query.Get("offset"))

	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 20
	}

	articles, err := service.GetFeed(user.Username, offset, limit)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
	}

	isFavorited, authors, _, err := service.GetArticleRelatedProperties(user, articles, false)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	articleResponses := make([]ArticleResponse, 0, len(articles))

	for i, article := range articles {
		articleResponses = append(articleResponses, ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[i],
			FavoritesCount: article.FavoritesCount,
			Author: AuthorResponse{
				Username:  authors[i].Username,
				Bio:       authors[i].Bio,
				Image:     authors[i].Image,
				Following: true,
			},
		})
	}

	response := AResponse{
		Articles:      articleResponses,
		ArticlesCount: len(articleResponses),
	}

	util.NewSuccessResponse(response, w, r)
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewErrorResponse(http.StatusUnauthorized, err, w)
		return
	}
	query := r.URL.Query()

	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 20
	}

	author := query.Get("author")
	tag := query.Get("tag")
	favorited := query.Get("favorited")

	articles, err := service.GetArticles(offset, limit, author, tag, favorited)
	if err != nil {
		util.NewErrorResponse(http.StatusNotFound, err, w)
		return
	}

	isFavorited, authors, following, err := service.GetArticleRelatedProperties(user, articles, true)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	articleResponses := make([]ArticleResponse, 0, len(articles))

	for i, article := range articles {
		articleResponses = append(articleResponses, ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        article.TagList,
			CreatedAt:      time.Unix(0, article.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, article.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[i],
			FavoritesCount: article.FavoritesCount,
			Author: AuthorResponse{
				Username:  authors[i].Username,
				Bio:       authors[i].Bio,
				Image:     authors[i].Image,
				Following: following[i],
			},
		})
	}

	response := AResponse{
		Articles:      articleResponses,
		ArticlesCount: len(articleResponses),
	}

	util.NewSuccessResponse(response, w, r)
}

func GetArticleSlug(w http.ResponseWriter, r *http.Request) {
	user, _, _ := service.GetCurrentUser(r.Header.Get("Authorization"))

	vars := mux.Vars(r)
	slug := vars["slug"]
	article, err := service.GetArticleBySlug(slug)
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

func PutArticleSlug(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	request := &APuRequest{}
	err = util.ParseBody(r, request)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	vars := mux.Vars(r)
	slug := vars["slug"]
	oldArticle, err := service.GetArticleBySlug(slug)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	newArticle := createNewArticle(*request, oldArticle)

	err = service.UpdateArticle(oldArticle, &newArticle)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	isFavorited, authors, following, err := service.GetArticleRelatedProperties(user, []model.Article{newArticle}, true)
	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	response := A1Response{
		Article: ArticleResponse{
			Slug:           newArticle.Slug,
			Title:          newArticle.Title,
			Description:    newArticle.Description,
			Body:           newArticle.Body,
			TagList:        newArticle.TagList,
			CreatedAt:      time.Unix(0, newArticle.CreatedAt).Format(model.TimestampFormat),
			UpdatedAt:      time.Unix(0, newArticle.UpdatedAt).Format(model.TimestampFormat),
			Favorited:      isFavorited[0],
			FavoritesCount: newArticle.FavoritesCount,
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

func createNewArticle(request APuRequest, oldArticle model.Article) model.Article {
	newArticle := model.Article{
		ArticleId:      oldArticle.ArticleId,
		Title:          request.Article.Title,
		Description:    request.Article.Description,
		Body:           request.Article.Body,
		TagList:        request.Article.TagList,
		CreatedAt:      oldArticle.CreatedAt,
		UpdatedAt:      time.Now().UTC().UnixNano(),
		FavoritesCount: oldArticle.FavoritesCount,
		Author:         oldArticle.Author,
	}

	if newArticle.Title == "" {
		newArticle.Title = oldArticle.Title
	}

	if newArticle.Description == "" {
		newArticle.Description = oldArticle.Description
	}

	if newArticle.Body == "" {
		newArticle.Body = oldArticle.Body
	}

	if newArticle.TagList == nil {
		newArticle.TagList = oldArticle.TagList
	}

	return newArticle
}

func PostArticles(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	request := &APuRequest{}
	err = util.ParseBody(r, request)
	if err != nil {
		util.NewErrorResponse(http.StatusBadRequest, err, w)
		return
	}
	now := time.Now().UTC()
	nowUnixNano := now.UnixNano()
	nowStr := now.Format(model.TimestampFormat)

	article := model.Article{
		Title:       request.Article.Title,
		Description: request.Article.Description,
		Body:        request.Article.Body,
		TagList:     request.Article.TagList, // TODO .distinct()
		CreatedAt:   nowUnixNano,
		UpdatedAt:   nowUnixNano,
		Author:      user.Username,
	}

	err = service.PutArticle(&article)
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
			CreatedAt:      nowStr,
			UpdatedAt:      nowStr,
			Favorited:      false,
			FavoritesCount: 0,
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

func DeleteArticleSlug(w http.ResponseWriter, r *http.Request) {
	user, _, err := service.GetCurrentUser(r.Header.Get("Authorization"))
	if err != nil {
		util.NewUnauthorizedResponse(w)
	}

	vars := mux.Vars(r)
	slug := vars["slug"]
	err = service.DeleteArticle(slug, user.Username)

	if err != nil {
		util.NewErrorResponse(http.StatusInternalServerError, err, w)
	}

	util.NewSuccessResponse(nil, w, r)
}
