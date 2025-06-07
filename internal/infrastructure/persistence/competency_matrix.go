package persistence

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
	"github.com/troxanna/pr-chat-backend/internal/domain/entity"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connString string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &Postgres{Pool: pool}, nil
}


type DBCompetencyMatrix struct {
	db *pgxpool.Pool
}

func NewDBCompetencyMatrix(db *pgxpool.Pool) DBCompetencyMatrix { return DBCompetencyMatrix{db: db} }

func (p DBCompetencyMatrix) CreateCompetencyMatrix(ctx context.Context, groups []entity.GroupSkills, skills []entity.Skill) error {
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
