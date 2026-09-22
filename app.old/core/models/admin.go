package models

type AdminDashboard struct {
	TotalStudents    int `json:"total_students"`
	TotalTeachers    int `json:"total_teachers"`
	TotalGroups      int `json:"total_groups"`
	TotalDisciplines int `json:"total_disciplines"`
	TotalLabs        int `json:"total_labs"`
	TotalCodes       int `json:"total_reg_codes"`
}
