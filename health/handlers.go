package health

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/util/response"
)

// Handles health check requests (GET /v1/health)
func (s *Service) healthcheck(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Raw("SELECT 1=1").Rows()
	defer rows.Close()

	var healthy bool
	if err == nil {
		healthy = true
	}

	response.WriteJSON(w, map[string]interface{}{
		"healthy": healthy,
	}, 200)
}
