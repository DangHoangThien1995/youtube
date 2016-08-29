package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"youtube/handle"
	"youtube/middleware"
	"youtube/utils/logs"
)

type Config struct {
	Server struct {
		Port string `json:"API_PORT"`
		Addr string `json:"API_ADDR"`
	} `json:"server"`
}

var l = logs.New("API_server")

type setupStruct struct {
	Config
	Handler http.Handler
}

func setup(cfg Config) *setupStruct {
	s := &setupStruct{Config: cfg}
	s.setupRouters()
	return s
}

func (this *setupStruct) setupRouters() {
	commonMids := commonMiddlewares()

	normal := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			commonMids(h).ServeHTTP(w, r)
		}
	}
	router := mux.NewRouter()

	newcontroller := handle.NewController()

	//router.Handle("/user/home", normal(newcontroller.Home))
	//router.Handle("/user/{id}", normal(newcontroller.Info)).Methods("GET")
	//log.Print("djfkaj")
	router.Handle("/user", normal(newcontroller.AddVideoFromYoutube)).Methods("GET")
	// router.Handle("/getvideo", normal(newcontroller.GetVideoFromDB)).Methods("POST")
	//router.Handle("/nextvideo", normal(newcontroller.GetNextVideo)).Methods("POST")
	//router.Handle("/login", normal(newcontroller.CheckIf)).Methods("POST")
	//router.Handle("/user/{id}", normal(newcontroller.DeleteId)).Methods("OPTIONS", "DELETE")
	//router.Handle("/user/{id}", normal(newcontroller.Update)).Methods("PUT")
	//router.Handle("/user", normal(newcontroller.RetrieveAllUser)).Methods("GET")
	//router.Handle("/room", r)

	http.ListenAndServe(":3030", router)
}

func commonMiddlewares() func(http.Handler) http.Handler {
	logger := middleware.NewLogger()
	recovery := middleware.NewRecovery()
	return func(h http.Handler) http.Handler {
		return logger(recovery(h))
	}
}
