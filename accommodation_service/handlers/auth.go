package handlers

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	token := r.Header.Get("Auth-Token")
	//	id, err := uuid.Parse(token)
	//	if err == nil {
	//		authenticated := domain.User{
	//			Id: id,
	//		}
	//		r = r.WithContext(context.WithValue(r.Context(), "auth", &authenticated))
	//		r = r.WithContext(context.WithValue(r.Context(), "token", token))
	//	}
	//	next.ServeHTTP(w, r)
	//})
	return nil
}
