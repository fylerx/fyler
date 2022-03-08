package middleware

import (
	"context"
	"net/http"

	"github.com/fylerx/fyler/internal/projects"
	u "github.com/fylerx/fyler/pkg/utils"
)

type ProjectMiddleware struct {
	Projects projects.Repository
}

func (pmw *ProjectMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-API-KEY")

		if token == "" {
			u.RespondWithError(w, http.StatusUnauthorized, "missing auth token")
			return
		}

		project, err := pmw.Projects.GetByAPIKey(token)
		if err != nil {
			u.RespondWithError(w, http.StatusUnauthorized, "missing project")
			return
		}

		ctx := context.WithValue(r.Context(), projects.ProjectKey, project)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
