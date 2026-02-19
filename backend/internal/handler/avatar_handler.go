package handler

import (
	"net/http"

	"github.com/genvid/backend/internal/model"
	"github.com/go-chi/chi/v5"
)

type AvatarHandler struct{}

func NewAvatarHandler() *AvatarHandler {
	return &AvatarHandler{}
}

func (h *AvatarHandler) List(w http.ResponseWriter, r *http.Request) {
	avatars := getMockAvatars()
	respondJSON(w, http.StatusOK, model.SuccessResponse(avatars))
}

func (h *AvatarHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	avatarID := chi.URLParam(r, "id")

	avatars := getMockAvatars()
	for _, avatar := range avatars {
		if avatar.ID == avatarID {
			respondJSON(w, http.StatusOK, model.SuccessResponse(avatar))
			return
		}
	}

	respondError(w, http.StatusNotFound, "NOT_FOUND", "Avatar not found", nil)
}

func getMockAvatars() []model.Avatar {
	return []model.Avatar{
		{
			ID:          "avatar-001",
			Name:        "emma",
			DisplayName: strPtr("Emma"),
			Gender:      strPtr("female"),
			AgeRange:    strPtr("20s"),
			Style:       "casual",
			Languages:   []string{"en", "es"},
			IsPremium:   false,
			UsageCount:  150,
		},
		{
			ID:          "avatar-002",
			Name:        "james",
			DisplayName: strPtr("James"),
			Gender:      strPtr("male"),
			AgeRange:    strPtr("30s"),
			Style:       "professional",
			Languages:   []string{"en"},
			IsPremium:   false,
			UsageCount:  120,
		},
		{
			ID:          "avatar-003",
			Name:        "sofia",
			DisplayName: strPtr("Sofia"),
			Gender:      strPtr("female"),
			AgeRange:    strPtr("20s"),
			Style:       "energetic",
			Languages:   []string{"en", "pt"},
			IsPremium:   true,
			UsageCount:  200,
		},
		{
			ID:          "avatar-004",
			Name:        "li",
			DisplayName: strPtr("Li"),
			Gender:      strPtr("male"),
			AgeRange:    strPtr("30s"),
			Style:       "friendly",
			Languages:   []string{"en", "zh"},
			IsPremium:   false,
			UsageCount:  95,
		},
		{
			ID:          "avatar-005",
			Name:        "maria",
			DisplayName: strPtr("Maria"),
			Gender:      strPtr("female"),
			AgeRange:    strPtr("40s"),
			Style:       "elegant",
			Languages:   []string{"en", "es", "pt"},
			IsPremium:   true,
			UsageCount:  180,
		},
		{
			ID:          "avatar-006",
			Name:        "alex",
			DisplayName: strPtr("Alex"),
			Gender:      strPtr("male"),
			AgeRange:    strPtr("20s"),
			Style:       "trendy",
			Languages:   []string{"en"},
			IsPremium:   false,
			UsageCount:  85,
		},
		{
			ID:          "avatar-007",
			Name:        "yuki",
			DisplayName: strPtr("Yuki"),
			Gender:      strPtr("female"),
			AgeRange:    strPtr("20s"),
			Style:       "casual",
			Languages:   []string{"en", "ja"},
			IsPremium:   true,
			UsageCount:  160,
		},
		{
			ID:          "avatar-008",
			Name:        "david",
			DisplayName: strPtr("David"),
			Gender:      strPtr("male"),
			AgeRange:    strPtr("40s"),
			Style:       "professional",
			Languages:   []string{"en"},
			IsPremium:   false,
			UsageCount:  70,
		},
	}
}

func strPtr(s string) *string {
	return &s
}
