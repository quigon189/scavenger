package database

import (
	"log"
	"scavenger/internal/auth"
	"scavenger/internal/models"
	"slices"
)

func (db *Database) SetTestData(cfg *models.Config) {

	roles, err := db.GetRoles()
	if err != nil {
		log.Fatalf("Failed to get roles: %v", err)
	}

	if !slices.Contains(roles, string(models.AdminRole)) {
		if err := db.CreateRole(string(models.AdminRole), "Admin role"); err != nil {
			log.Fatalf("Failed to create admin role")
		} else {
			log.Printf("Role admin created")
		}
	}

	if !slices.Contains(roles, string(models.StudentRole)) {
		if err := db.CreateRole(string(models.StudentRole), "Student role"); err != nil {
			log.Fatalf("Failed to create student role")
		} else {
			log.Printf("Role student created")
		}
	}

	for _, admin := range cfg.TestData.Roles.Admin {
		_, err := db.GetUserByUsername(admin.Username)
		if err != nil {
			admin, err := auth.RegisterUser(admin.Username, admin.Name, admin.PasswordHash, string(models.AdminRole))
			if err != nil {
				log.Printf("Failed to register admin: %v", err)
				continue
			}
			if err := db.CreateUser(admin); err != nil {
				log.Printf("Failed to create admin %s: %v", admin.Name, err)
			} else {
				log.Printf("User %s created", admin.Name)
			}
		}
	}

	groups, err := db.GetAllGroups()
	if err != nil {
		log.Fatalf("Failed to get groups: %v", err)
	}
	groupsNames := make([]string, 1)
	for _, g := range groups {
		groupsNames = append(groupsNames, g.Name)
	}

	for groupName, students := range cfg.TestData.Roles.Student {
		if !slices.Contains(groupsNames, groupName) {
			group := models.Group{Name: groupName}
			if err := db.CreateGroup(&group); err != nil {
				log.Printf("Failed to create group %s: %v", groupName, err)
				continue
			} else {
				log.Printf("Group %s created", groupName)
			}
			for _, student := range students {
				student, err := auth.RegisterUser(student.Username, student.Name, student.PasswordHash, "student")
				if err != nil {
					log.Printf("Failed to register student %v: %v", student, err)
					continue
				}

				student.GroupName = groupName
				if err := db.CreateStudent(student); err != nil {
					log.Printf("Failed to create student %+v: %v", student, err)
				} else {
					log.Printf("Student %+v created", student)
				}
			}
		}

	}
}
