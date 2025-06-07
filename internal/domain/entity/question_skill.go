package entity

type QuestionSkill struct {
	ID string `db:"id"`
	Content string `db:"content"`
	SkillID string `db:"skill_id"`
	Level int `db:"level"`
}