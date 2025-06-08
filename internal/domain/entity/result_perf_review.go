package entity

type ResultPerfRewiew struct {
	UserName string
	GroupsSkills []GroupSkillsResult
}

type GroupSkillsResult struct {
	Name string
	Average  float64
	Skills []SkillResult
}

type SkillResult struct {
	Name string 
	Score  float64 
}
