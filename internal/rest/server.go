package rest

import (
	"net/http"
	"yaly-1/internal/rest/handlers"

	"log"
)

// Сервер REST
type Server struct {
	httpServer *http.Server
}

// Конструктор для сервера
func NewServer(port string, service handlers.CalcService) *Server {
	mux := http.NewServeMux()
	server := &Server{httpServer: &http.Server{Addr: ":" + port, Handler: mux}}

	calcHandlers := handlers.NewCalcHandlers(service)

	mux.HandleFunc("/api/v1/calculate", calcHandlers.Calculate)
	return server
}

// Запуск сервера
func (s *Server) Run() error {
	log.Println("Strting server on port " + s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Завершение работы сервера
func (s *Server) Shutdown() error {
	log.Println("Shutting down server")
	return s.httpServer.Shutdown(nil)
}
