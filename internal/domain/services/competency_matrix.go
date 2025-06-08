package service

import (
	"context"
	"fmt"

	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
)

type dbCompetencyMatrix interface {
	CreateCompetencyMatrix(context.Context, []entity.GroupSkills, []entity.Skill, entity.Matrix) error
	GetCompetencyMatrixs(context.Context) ([]entity.GroupSkills, error)
}

type CompetencyMatrix struct {
	db dbCompetencyMatrix
}

func NewCompetencyMatrix(
	db dbCompetencyMatrix,
) *CompetencyMatrix {
	return &CompetencyMatrix{
		db: db,
	}
}

func (p CompetencyMatrix) CreateCompetencyMatrix(ctx context.Context, groups []entity.GroupSkills, skills []entity.Skill, matrix entity.Matrix) error {
	if err := p.db.CreateCompetencyMatrix(ctx, groups, skills, matrix); err != nil {
		return fmt.Errorf("db.CreateCompetencyMatrix: %w", err)
	}
	return nil
}

func (p CompetencyMatrix) GetCompetencyMatrixs(ctx context.Context) ([]entity.GroupSkills, error) {
	return []entity.GroupSkills{}, nil
	// if err := p.db.CreateCompetencyMatrix(ctx, groups, skills); err != nil {
	// 	return fmt.Errorf("db.CreateCompetencyMatrix: %w", err)
	// }
	// return nil
}
