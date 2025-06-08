package entity

type Matrix struct {
	ID string `db:"id"`
	Name string `db:"name"`
	GroupsSkills []GroupSkills 
}