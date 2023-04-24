package routes

import (
	"encoding/json"
	"net/http"

	"realworld-go-nolambda/controller"
	
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(map[string]string{"data": "Hello from Mux & mongoDB"})
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/articles/feed", controller.GetArticlesFeed).Methods("GET")
	router.HandleFunc("/articles", controller.GetArticles).Methods("GET")
	router.HandleFunc("/articles/{slug}", controller.GetArticleSlug).Methods("GET")
	router.HandleFunc("/articles/{slug}", controller.PutArticleSlug).Methods("PUT")
	router.HandleFunc("/articles", controller.PostArticles).Methods("POST")
	router.HandleFunc("/articles/{slug}", controller.DeleteArticleSlug).Methods("DELETE")

	router.HandleFunc("/articles/{slug}/comments/{id}", controller.DeleteComment).Methods("DELETE")
	router.HandleFunc("/articles/{slug}/comments", controller.GetComments).Methods("GET")
	router.HandleFunc("/articles/{slug}/comments", controller.PostComment).Methods("POST")

	router.HandleFunc("/articles/{slug}/favorite", controller.DeleteFavorite).Methods("DELETE")
	router.HandleFunc("/articles/{slug}/favorite", controller.PostFavorite).Methods("POST")

	router.HandleFunc("/profiles/{username}/follow", controller.DeleteProfileFollow).Methods("DELETE")
	router.HandleFunc("/profiles/{username}/follow", controller.PostProfileFollow).Methods("POST")
	router.HandleFunc("/profiles/{username}", controller.GetProfiles).Methods("GET")


	router.HandleFunc("/tags", controller.GetTags).Methods("GET")

	router.HandleFunc("/user", controller.GetUser).Methods("GET")
	router.HandleFunc("/users/login", controller.UserLogin).Methods("POST")
	router.HandleFunc("/users", controller.PostUser).Methods("POST")
	router.HandleFunc("/user", controller.PutUser).Methods("PUT")
}