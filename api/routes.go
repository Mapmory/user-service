package main

import (
    "github.com/gorilla/mux"
)

func setupRoutes() *mux.Router {
    router := mux.NewRouter()

    // Public routes
    router.HandleFunc("/api/users", RegisterHandler).Methods("POST")
    router.HandleFunc("/api/auth/login", LoginHandler).Methods("POST")

    // Protected routes
    router.HandleFunc("/api/users/me", TokenVerifyMiddleware(UserInfoHandler)).Methods("GET")

    return router
}