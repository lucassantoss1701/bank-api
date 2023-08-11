package webserver

import (
	"context"
	"log"
	"lucassantoss1701/bank/internal/entity"
	"lucassantoss1701/bank/internal/infra/web/responses"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	customMiddleware "lucassantoss1701/bank/internal/infra/web/webserver/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	http.HandlerFunc
	method string
	path   string
}

type WebServer struct {
	Router        chi.Router
	Handlers      []Handler
	srv           *http.Server
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      []Handler{},
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, method string, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, Handler{
		path:        path,
		method:      method,
		HandlerFunc: handler,
	})
}

func (s *WebServer) Start() {
	s.startCHI()
	server := &http.Server{Addr: s.WebServerPort, Handler: s.Router}
	s.srv = server

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func (s *WebServer) Stop() {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	err := s.srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Println("Error stopping the server:", err)
		return
	}

	log.Println("Server stopped gracefully")
}

func (s *WebServer) startCHI() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(customMiddleware.Auth)

	s.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var errorHandler entity.ErrorHandler
		errorHandler.TypeError = entity.NOT_FOUND_ERROR
		errorHandler.Add("route does not exist")
		err = &errorHandler
		responses.Err(w, err)
	})

	s.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var errorHandler entity.ErrorHandler
		errorHandler.TypeError = entity.NOT_ALLOWED_ERROR
		errorHandler.Add("method is not valid")
		err = &errorHandler
		responses.Err(w, err)
	})

	for _, handler := range s.Handlers {
		s.Router.Method(handler.method, handler.path, handler.HandlerFunc)
	}
}
