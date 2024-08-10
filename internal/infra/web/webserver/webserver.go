package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      []handlerConfig
	WebServerPort string
}

type handlerConfig struct {
	path        string
	method      string
	handlerFunc http.HandlerFunc
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      []handlerConfig{},
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handlerFunc http.HandlerFunc) {
	s.Handlers = append(s.Handlers, handlerConfig{
		path,
		method,
		handlerFunc,
	})
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, conf := range s.Handlers {
		s.Router.Method(conf.method, conf.path, conf.handlerFunc)
	}
	err := http.ListenAndServe(s.WebServerPort, s.Router)
	if err != nil {
		panic(err)
	}
}
