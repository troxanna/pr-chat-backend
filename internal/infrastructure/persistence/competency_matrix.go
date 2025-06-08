package persistence

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
)

type DBCompetencyMatrix struct {
	db *pgxpool.Pool
}

func NewDBCompetencyMatrix(db *pgxpool.Pool) DBCompetencyMatrix { return DBCompetencyMatrix{db: db} }

func (p DBCompetencyMatrix) CreateCompetencyMatrix(ctx context.Context, groups []entity.GroupSkills, skills []entity.Skill) error {
	log.Println(groups)
	log.Println(skills)
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
