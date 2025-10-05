package handlers

import (
	"log"
	"net/http"
	"strconv"

	"scavenger/internal/alerts"
	"scavenger/internal/auth"
	"scavenger/internal/models"
	"scavenger/views"
)

func (h *Handler) StudentsManager(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroups()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при получении групп")
		log.Printf("Failed to get groups: %v", err)
		groups = []models.Group{}
	}

	students, err := h.db.GetAllStudents()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при получении студентов")
		log.Printf("Failed to get users: %v", err)
		students = []models.User{}
	}

	views.StudentsManager(students, groups).Render(r.Context(), w)
}

func (h *Handler) AddStudents(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	var student struct {
		name     string
		username string
		password string
	}

	student.name = r.FormValue("name")
	student.username = r.FormValue("username")
	student.password = r.FormValue("password")

	var groupID int
	groupID, err = strconv.Atoi(r.FormValue("group"))
	if err != nil {
		alerts.FlashError(w, r, "Ошибка определения группы: Группа не выбрана")
		log.Printf("Failed to atoi groupID: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	group, err := h.db.GetGroupByID(groupID)
	if err != nil {
		alerts.FlashError(w, r, "Группа не найдена")
		log.Printf("Failed to get group from db: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	newStudent, err := auth.RegisterUser(
		student.username,
		student.name,
		student.password,
		string(models.StudentRole),
	)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при регистрации пользователя")
		log.Printf("Failed to register user (student): %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	newStudent.GroupName = group.Name
	newStudent.GroupID = group.ID
	log.Printf("Create user (student): %+v", newStudent)

	err = h.db.CreateStudent(newStudent)
	if err != nil {
		alerts.FlashWarning(w, r, "Студент создан, но не привязан к группе")
		log.Printf("Failed to create student: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Студент "+student.name+" создан")
	http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
}

func (h *Handler) EditStudent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	var student struct {
		name     string
		username string
		password string
	}

	student.name = r.FormValue("name")
	student.username = r.FormValue("username")
	student.password = r.FormValue("password")

	log.Printf("Edit student: %+v", student)

	var groupID int
	groupID, err = strconv.Atoi(r.FormValue("group"))
	if err != nil {
		alerts.FlashError(w, r, "Ошибка определения группы")
		log.Printf("Failed to atoi groupID: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	editStudent, err := h.db.GetStudentByUsername(student.username)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка: студент не найден")
		log.Printf("Failed to get student from db: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	if student.name != "" {
		editStudent.Name = student.name
	}

	if student.password != "" {
		passHash, err := auth.GeneratePassHash(student.password)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка хэширования пароля")
			log.Printf("Failed to hash password: %v", err)
			http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
			return
		}
		editStudent.PasswordHash = passHash
	}

	editStudent.GroupID = groupID

	err = h.db.UpadateStudent(editStudent)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при обновлении данных студента")
		log.Printf("Failed to update student: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Данные обновлены")
	http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
}

func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	err := r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	student, err := h.db.GetStudentByUsername(username)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения пользователя БД")
		log.Printf("Failed to get student from db: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	if student.Username != r.FormValue("username") {
		alerts.FlashWarning(w, r, "Введен неверный логин")
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	err = h.db.DeleteStudent(student)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка удаления студента из БД")
		log.Printf("Failed to delete student from db: %v", err)
		http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Студент удален")
	http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
}
