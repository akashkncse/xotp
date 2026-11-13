package api

import (
	"net/http"
)

func (s *Server) aliveHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, aliveResponse{Alive: true})
}
