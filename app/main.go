package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go-chat-app/my_handlers"
	"log"
	"net/http"
)

func main() {
	forumsPath := "/api/forums"
	router := mux.NewRouter()
	router.HandleFunc("/api/register", my_handlers.Register).Methods("POST")
	router.HandleFunc("/api/login", my_handlers.Login).Methods("POST")
	router.HandleFunc("/api/messages", my_handlers.GetMessages).Methods("GET")
	router.HandleFunc("/api/messages", my_handlers.CreateMessage).Methods("POST")
	router.HandleFunc(forumsPath, my_handlers.GetForums).Methods("GET")
	router.HandleFunc(forumsPath, my_handlers.CreateForum).Methods("POST")
	router.HandleFunc(forumsPath, my_handlers.RemoveForum).Methods("DELETE")
	router.HandleFunc("/api/user", my_handlers.GetUserInfo).Methods("GET")
	log.Println("Server started on :8081")
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://nginx:80"})
	err := http.ListenAndServe(":8081", handlers.CORS(headers, methods, origins)(router))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
