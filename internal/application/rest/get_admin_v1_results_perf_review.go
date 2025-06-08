package rest

import (
	"encoding/json"
	"net/http"

	// "github.com/troxanna/pr-chat-backend/internal/domain/entity"
	// "github.com/google/uuid"
)

type ResultPerfRewiewModel struct {
	UserName string `json:"user_id"`
	GroupsSkills []GroupSkillsResultModel `json:"directions" validate:"required"`
}

type GroupSkillsResultModel struct {
	Name string `json:"name"`
	Average  float64  `json:"average"`
	Skills []SkillResultModel `json:"competencies"`
}

type SkillResultModel struct {
	Name string `json:"name"`
	Score  float64  `json:"score"`
}


// type GetAdminResultsPerfReviewResponse struct {
// 	Matrixs []MatrixModel `json:"matrixs" validate:"required"`
// }

func (s ServerAdmin) GetAdminV1ResultsPerfReview(w http.ResponseWriter, r *http.Request) error {
	response := []ResultPerfRewiewModel{
		{
			UserName: "Анна Ларионова",
			GroupsSkills: []GroupSkillsResultModel{
		{
					Name:    "Security",
					Average: 2.5,
					Skills: []SkillResultModel{
						{Name: "WAF", Score: 2},
						{Name: "Firewall", Score: 2},
						{Name: "Security", Score: 3},
					},
				},
			},
		},
		{
			UserName: "Егор Гусяков",
			GroupsSkills: []GroupSkillsResultModel{
				{
					Name:    "DB",
					Average: 4.17,
					Skills: []SkillResultModel{
						{Name: "PostgreSQL", Score: 5},
						{Name: "MySQL/MariaDB", Score: 4},
						{Name: "ClickHouse", Score: 4},
						{Name: "MS SQL", Score: 3},
						{Name: "Redis", Score: 4},
						{Name: "MongoDB", Score: 5},
					},
				},
				{
					Name:    "Automation",
					Average: 4.29,
					Skills: []SkillResultModel{
						{Name: "Bash", Score: 4},
						{Name: "Python", Score: 5},
						{Name: "Terraform", Score: 5},
						{Name: "Ansible", Score: 5},
						{Name: "Helm", Score: 4},
						{Name: "Gitlab-ci", Score: 4},
						{Name: "Puppet", Score: 3},
					},
				},
				{
					Name:    "Other",
					Average: 4.6,
					Skills: []SkillResultModel{
						{Name: "Kubernetes", Score: 5},
						{Name: "Yandex Cloud", Score: 4},
						{Name: "Git-flow", Score: 5},
						{Name: "Jira-flow", Score: 4},
						{Name: "Docker", Score: 5},
					},
				},
				{
					Name:    "OS",
					Average: 4.2,
					Skills: []SkillResultModel{
						{Name: "RHEL", Score: 4},
						{Name: "CentOS", Score: 5},
						{Name: "Ubuntu", Score: 5},
						{Name: "Debian", Score: 4},
						{Name: "Windows Server", Score: 3},
					},
				},
				{
					Name:    "Logs/Monitoring",
					Average: 4.13,
					Skills: []SkillResultModel{
						{Name: "ELK", Score: 4},
						{Name: "Zabbix", Score: 5},
						{Name: "Prometheus", Score: 5},
						{Name: "Grafana", Score: 5},
						{Name: "Fluent-bit", Score: 4},
						{Name: "OpenTelemetry", Score: 3},
						{Name: "MicroMetrics", Score: 4},
						{Name: "VictoriaMetrics", Score: 3},
					},
				},
			},
		},
	}	
	

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}