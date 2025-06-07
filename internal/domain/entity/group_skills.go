package entity

type GroupSkills struct {
	ID string `db:"id"`
	Name string `db:"name"`
	Description string `db:"description"`
	Type string  `db:"type"`
}