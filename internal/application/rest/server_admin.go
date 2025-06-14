package rest

import (
	"context"

	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
)


type competencyMatrix interface {
	CreateCompetencyMatrix(context.Context, []entity.GroupSkills, []entity.Skill, entity.Matrix) error
	GetCompetencyMatrixs(context.Context) ([]entity.GroupSkills, error)
	// UpdateCompetencyMatrix(context.Context, value.PanelID, *value.PanelID, *string, *string) error
	// DeleteCompetencyMatrix(context.Context, value.PanelID) error
}

type ServerAdmin struct {
	competencyMatrix  competencyMatrix
}

func NewServerAdmin(
	competencyMatrix competencyMatrix,
) ServerAdmin {
	return ServerAdmin{
		competencyMatrix: competencyMatrix,
	}
}

