package database

import (
	"log"
	"scavenger/internal/auth"
	"scavenger/internal/models"
)

func (db *Database) SetTestData(cfg *models.Config) {
	log.Printf("Test data: %+v", cfg.TestData)
	if err := db.CreateRole("admin", "Admin role"); err != nil {
		log.Printf("Failed to create admin role")
	} else {
		log.Printf("Role admin created")
	}

	if err := db.CreateRole("student", "Student role"); err != nil {
		log.Printf("Failed to create student role")
	} else {
		log.Printf("Role student created")
	}

	for _, admin := range cfg.TestData.Roles.Admin {
		admin, err := auth.RegisterUser(admin.Username, admin.Name, admin.PasswordHash, "admin")
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

	for groupName, students := range cfg.TestData.Roles.Student {
		group := models.Group{Name: groupName}
		if err := db.CreateGroup(&group); err != nil {
			log.Printf("Failed to create group %s: %v", groupName, err)
		} else {
			log.Printf("Group %s created", groupName)
		}
		for _, student := range students {
			student, err := auth.RegisterUser(student.Username, student.Name, student.PasswordHash, "student")
			if err != nil {
				log.Printf("Failed to register student %v: %v", student, err)
				continue
			}
			if err := db.CreateUser(student); err != nil {
				log.Printf("Failed to create user(student) %v: %v", student, err)
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
