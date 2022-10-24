package controllers

import "github.com/tifarin/fullstack/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Photos routes
	s.Router.HandleFunc("/photos", middlewares.SetMiddlewareJSON(s.CreatePhoto)).Methods("POST")
	s.Router.HandleFunc("/photos", middlewares.SetMiddlewareJSON(s.GetPhotos)).Methods("GET")
	s.Router.HandleFunc("/photos/{id}", middlewares.SetMiddlewareJSON(s.GetPhoto)).Methods("GET")
	s.Router.HandleFunc("/photos/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePhoto))).Methods("PUT")
	s.Router.HandleFunc("/photos/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePhoto)).Methods("DELETE")

	//Coment routes
	s.Router.HandleFunc("/comments", middlewares.SetMiddlewareJSON(s.CreateComment)).Methods("POST")
	s.Router.HandleFunc("/comments", middlewares.SetMiddlewareJSON(s.GetComments)).Methods("GET")
	s.Router.HandleFunc("/comments/{id}", middlewares.SetMiddlewareJSON(s.GetComment)).Methods("GET")
	s.Router.HandleFunc("/comments/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateComment))).Methods("PUT")
	s.Router.HandleFunc("/comments/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteAComment)).Methods("DELETE")

	//Coment media sosial
	s.Router.HandleFunc("/media_sosial", middlewares.SetMiddlewareJSON(s.CreateMediaSosial)).Methods("POST")
	s.Router.HandleFunc("/media_sosial", middlewares.SetMiddlewareJSON(s.GetMediaSosials)).Methods("GET")
	s.Router.HandleFunc("/media_sosial/{id}", middlewares.SetMiddlewareJSON(s.GetMediaSosial)).Methods("GET")
	s.Router.HandleFunc("/media_sosial/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateAMediaSosial))).Methods("PUT")
	s.Router.HandleFunc("/media_sosial/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteAMediaSosial)).Methods("DELETE")
}
