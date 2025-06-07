package rest

import (
	"fmt"
	"net/http"

	// "github.com/troxanna/pr-chat-backend/internal/domain/entity"
	// "github.com/google/uuid"
)

type GetAdminCompetencyMatrixResponse struct {
	GroupSkills []GroupSkillsModel `json:"groups" validate:"required"`
}

func (s ServerAdmin) GetAdminV1CompetencyMatrix(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("test")
	return nil
}