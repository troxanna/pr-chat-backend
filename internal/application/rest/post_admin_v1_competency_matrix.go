package rest

import (
	"fmt"
	"net/http"

	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
	"github.com/google/uuid"
)

type GroupSkillsModel struct {
	Name string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Skills []SkillModel `json:"skills" validate:"required"`
	Type string `json:"type"`
}

type SkillModel struct {
	Name string `json:"name"`
}

// fix for one item
type PostAdminCompetencyMatrixRequest struct {
	Name string `json:"name" validate:"required"`
	GroupsSkills []GroupSkillsModel `json:"groups" validate:"required"`
}

// type PostAdminV1ProjectsJSONRequestBody = PostAdminCompetencyMatrixRequest


func (s ServerAdmin) PostAdminV1CompetencyMatrix(w http.ResponseWriter, r *http.Request) error {
	var request PostAdminCompetencyMatrixRequest

	if err := readRequest(r, &request); err != nil {
		return fmt.Errorf("readRequest: %w", err)
	}
	var groups []entity.GroupSkills
	var skills []entity.Skill

	for _, group := range request.GroupsSkills {
		groupID := uuid.New().String()
		groups = append(groups, entity.GroupSkills{
			ID:          groupID,
			Name:        group.Name,
			Description: group.Description,
			Type:        group.Type,
		})

		for _, s := range group.Skills {
			skills = append(skills, entity.Skill{
				ID:          uuid.New().String(),
				Name:        s.Name,
				Description: "",
				GroupID:     groupID,
			})
		}
	}

	matrix := entity.Matrix{
		Name: request.Name,
		GroupsSkills: groups,
	}

	if err := s.competencyMatrix.CreateCompetencyMatrix(r.Context(), groups, skills, matrix); err != nil {
		return fmt.Errorf("CreateCompetencyMatrix: %w", err)
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

