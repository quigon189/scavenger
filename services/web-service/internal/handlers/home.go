package handlers

import (
	"net/http"
	"web-service/internal/models"
	"web-service/internal/views/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Username: "test",
		Name: "Test",
		Role: "test role",
		Email: "test@test",
		Group: "test123",
	}

	Render(w,r,pages.Home(user))
}
