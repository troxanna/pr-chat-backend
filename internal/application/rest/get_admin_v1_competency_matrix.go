package rest

import (
	"encoding/json"
	"net/http"

	// "github.com/troxanna/pr-chat-backend/internal/domain/entity"
	// "github.com/google/uuid"
)

type GetAdminCompetencyMatrixResponse struct {
	Matrixs []MatrixModel `json:"matrixs" validate:"required"`
}

func (s ServerAdmin) GetAdminV1CompetencyMatrix(w http.ResponseWriter, r *http.Request) error {
	response := GetAdminCompetencyMatrixResponse{
		Matrixs: []MatrixModel{
			{
				Name: "DevOps Competency Matrix",
				GroupsSkills: []GroupSkillsModel{
					{
						Name:        "Automation",
						Description: "Инструменты для автоматизации",
						Type:        "hard",
						Skills: []SkillModel{
							{Name: "Ansible"},
							{Name: "Terraform"},
							{Name: "Gitlab-CI"},
						},
					},
					{
						Name:        "Monitoring",
						Description: "Средства мониторинга и логирования",
						Type:        "hard",
						Skills: []SkillModel{
							{Name: "Prometheus"},
							{Name: "Grafana"},
							{Name: "Zabbix"},
						},
					},
					{
						Name: "Soft Skills",
						Type: "soft",
						Skills: []SkillModel{
							{Name: "Коммуникация"},
							{Name: "Презентация задач"},
						},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return nil
}