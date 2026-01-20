package models

import (
	"net/http"
	"strconv"
)

type ReportFilterParams struct {
	DisciplineID  int
	LabID         int
	Status        string
	Grade         int
	StudentSearch string
	Period        string
	SortBy        string
	SortOrder     string
	Page          int
	PageSize      int
	TotalPages    int
}

func (f *ReportFilterParams) Parse(r *http.Request) {
	f.DisciplineID, _ = strconv.Atoi(r.URL.Query().Get("discipline_id"))	
	f.LabID, _ = strconv.Atoi(r.URL.Query().Get("lab_id"))
	f.Status = r.URL.Query().Get("status")
	f.Grade, _ =strconv.Atoi(r.URL.Query().Get("grade"))
	f.StudentSearch = r.URL.Query().Get("student_search")
	f.Period = r.URL.Query().Get("period")

	f.SortBy = r.URL.Query().Get("sort_by")
	f.SortOrder = r.URL.Query().Get("sort_order")
	if f.SortOrder == "" {
		f.SortOrder = "desk"
	}

	f.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if f.Page == 0 {
		f.Page = 1
	}

	f.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	if f.PageSize == 0 {
		f.PageSize = 20
	}
}

func (f *ReportFilterParams) IsEmpty() bool {
	return f.DisciplineID == 0 &&
		   f.LabID == 0 &&
		   f.Status == "" &&
		   f.Grade == 0 &&
		   f.StudentSearch == "" &&
		   f.Period == "" &&
		   f.SortBy == ""
}
