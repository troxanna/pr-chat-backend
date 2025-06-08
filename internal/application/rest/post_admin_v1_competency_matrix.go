package rest

import (
	"fmt"
	"net/http"

	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
	"github.com/google/uuid"
)

type MatrixModel struct {
	Name string `json:"name" validate:"required"`
	GroupsSkills []GroupSkillsModel `json:"groups" validate:"required"`
}

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
	Matrixs []MatrixModel `json:"matrixs" validate:"required"`
}

// type PostAdminV1ProjectsJSONRequestBody = PostAdminCompetencyMatrixRequest


func (s ServerAdmin) PostAdminV1CompetencyMatrix(w http.ResponseWriter, r *http.Request) error {
	var request PostAdminCompetencyMatrixRequest

	if err := readRequest(r, &request); err != nil {
		return fmt.Errorf("readRequest: %w", err)
	}
	var groups []entity.GroupSkills
	var skills []entity.Skill
	var matrixs []entity.Matrix

	for _, matrix := range request.Matrixs {
		matrixID := uuid.New().String()

		matrixs = append(matrixs, entity.Matrix{
			ID:   matrixID,
			Name: matrix.Name,
		})

		for _, group := range matrix.GroupsSkills {
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
	}

	if err := s.competencyMatrix.CreateCompetencyMatrix(r.Context(), groups, skills, matrixs); err != nil {
		return fmt.Errorf("CreateCompetencyMatrix: %w", err)
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

