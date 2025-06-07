package service

import (
	"context"
	"fmt"

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
	fmt.Println("Hello")
	return nil

	// if err := p.db.CreateProject(ctx, entity.Project{ //nolint:exhaustruct
	// 	ProjectID:   value.ProjectID(xid.New().String()),
	// 	PanelID:     panelID,
	// 	Name:        name,
	// 	Description: description,
	// 	CreatedAt:   time.Now().UTC(),
	// }); err != nil {
	// 	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
	// 		return failure.NewInvalidArgumentErrorFromError(
	// 			fmt.Errorf("db.CreateProject: %w", err),
	// 			failure.WithDescription("Already exist"))
	// 	}

	// 	return fmt.Errorf("db.CreateProject: %w", err)
	// }

	// return nil
}
