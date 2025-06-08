package service

import (
	"context"
	"fmt"
	"log"

	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
)

type dbCompetencyMatrix interface {
	CreateCompetencyMatrix(context.Context, []entity.GroupSkills, []entity.Skill) error
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

func (p CompetencyMatrix) CreateCompetencyMatrix(ctx context.Context, groups []entity.GroupSkills, skills []entity.Skill) error {
	if err := p.db.CreateCompetencyMatrix(ctx, groups, skills); err != nil {
		return fmt.Errorf("db.CreateCompetencyMatrix: %w", err)
	}
	log.Println("test1")
	return nil
}
