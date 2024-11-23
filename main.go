package main

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"
)

type InMemoryStudentManager struct {
	students         []model.Student
	loginAttempts    map[string]int
	maxLoginAttempts int
}

func NewInMemoryStudentManager() *InMemoryStudentManager {
	return &InMemoryStudentManager{
		students: []model.Student{
			{ID: "A12345", Name: "Aditira", StudyProgram: "TI"},
			{ID: "B21313", Name: "Dito", StudyProgram: "TK"},
			{ID: "A34555", Name: "Afis", StudyProgram: "MI"},
		},
		loginAttempts:    make(map[string]int),
		maxLoginAttempts: 3,
	}
}

func (sm *InMemoryStudentManager) GetStudents() []model.Student {
	return sm.students
}

func (sm *InMemoryStudentManager) Login(id, name string) (string, error) {
	if id == "" || name == "" {
		return "", errors.New("ID and Name must not be empty")
	}

	
	if sm.loginAttempts[id] >= sm.maxLoginAttempts {
		return "", errors.New("Login gagal: Batas maksimum login terlampaui")
	}

	for _, student := range sm.students {
		if student.ID == id && student.Name == name {
			sm.loginAttempts[id] = 0 
			return fmt.Sprintf("Login berhasil: Selamat datang %s! Kamu terdaftar di program studi: %s", name, sm.GetStudyProgramName(student.StudyProgram)), nil
		}
	}

	
	sm.loginAttempts[id]++
	return "", errors.New("Login gagal: data mahasiswa tidak ditemukan")
}

func (sm *InMemoryStudentManager) Register(id, name, studyProgram string) (string, error) {
	if id == "" || name == "" || studyProgram == "" {
		return "", errors.New("ID, Name or StudyProgram is undefined!")
	}

	for _, student := range sm.students {
		if student.ID == id {
			return "", errors.New("Registrasi gagal: id sudah digunakan")
		}
	}

	if _, err := sm.GetStudyProgram(studyProgram); err != nil {
		return "", fmt.Errorf("Study program %s is not found", studyProgram)
	}

	sm.students = append(sm.students, model.Student{ID: id, Name: name, StudyProgram: studyProgram})
	return fmt.Sprintf("Registrasi berhasil: %s (%s)", name, studyProgram), nil
}

func (sm *InMemoryStudentManager) GetStudyProgram(code string) (string, error) {
	programs := map[string]string{
		"TI": "Teknik Informatika",
		"TK": "Teknik Komputer",
		"MI": "Manajemen Informatika",
		"SI": "Sistem Informasi",
	}
	if name, ok := programs[code]; ok {
		return name, nil
	}
	return "", errors.New("Study program not found")
}

func (sm *InMemoryStudentManager) GetStudyProgramName(code string) string {
	name, _ := sm.GetStudyProgram(code)
	return name
}

func (sm *InMemoryStudentManager) ModifyStudent(name string, modifier func(*model.Student)) (string, error) {
	for i, student := range sm.students {
		if student.Name == name {
			modifier(&sm.students[i])
			return "Program studi mahasiswa berhasil diubah.", nil
		}
	}
	return "", errors.New("Mahasiswa tidak ditemukan")
}

func (sm *InMemoryStudentManager) ChangeStudyProgram(newProgram string) func(*model.Student) {
	return func(s *model.Student) {
		s.StudyProgram = newProgram
	}
}

func (sm *InMemoryStudentManager) ImportStudents(filepaths []string) error {
	ch := make(chan error, len(filepaths)) 

	for _, filepath := range filepaths {
		go func(path string) {
			file, err := os.Open(path)
			if err != nil {
				ch <- err
				return
			}
			defer file.Close()

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				ch <- err
				return
			}

			time.Sleep(50 * time.Millisecond) 

			sm.addStudents(records)
			ch <- nil
		}(filepath)
	}

	
	for range filepaths {
		if err := <-ch; err != nil {
			return err
		}
	}

	return nil
}

func (sm *InMemoryStudentManager) addStudents(records [][]string) {
	sm.students = append(sm.students, sm.parseRecords(records)...)
}

func (sm *InMemoryStudentManager) parseRecords(records [][]string) []model.Student {
	var newStudents []model.Student
	for _, record := range records {
		if len(record) >= 3 {
			newStudents = append(newStudents, model.Student{
				ID:           record[0],
				Name:         record[1],
				StudyProgram: record[2],
			})
		}
	}
	return newStudents
}

func (sm *InMemoryStudentManager) SubmitAssignments(count int) {
	time.Sleep(150 * time.Millisecond)
}