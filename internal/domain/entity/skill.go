package entity

type Skill struct {
	ID string `db:"id"`
	Name string `db:"name"`
	Description string `db:"description"`
	GroupID string `db:"group_id"`
}