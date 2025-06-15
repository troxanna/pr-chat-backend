package persistence

import (
	"context"
	"log"

	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

type DBCompetencyMatrix struct {
	db *sqlx.DB
}

func NewDBCompetencyMatrix(db *sqlx.DB) DBCompetencyMatrix { return DBCompetencyMatrix{db: db} }

func (p DBCompetencyMatrix) CreateCompetencyMatrix(ctx context.Context, groups []entity.GroupSkills, skills []entity.Skill, matrix entity.Matrix) error {
	log.Println(groups)
	log.Println(skills)
	log.Println(matrix)
	// const query = `
	// 	INSERT INTO projects(
	// 		id,
	// 		panel_id,
	// 		name,
	// 		description,
	// 		created_at
	// 	) VALUES (
	// 		:id,
	// 		:panel_id,
	// 		:name,
	// 		:description,
	// 		:created_at
	// 	)
	// `

	// if _, err := p.db.NamedExecContext(ctx, query, project); err != nil {
	// 	return fmt.Errorf("db.NamedExecContext: %w", err)
	// }

	return nil
}

func (p DBCompetencyMatrix) GetCompetencyMatrixs(ctx context.Context) ([]entity.GroupSkills, error) {
	return []entity.GroupSkills{}, nil
}