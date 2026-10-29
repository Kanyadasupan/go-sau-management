// ทำงานกับ DB-CRUD
package main

import "gorm.io/gorm"

// ฟังก์ชั่นเพิ่มข้อมูลลงตาราง
func createStudent(db *gorm.DB, student *Student) error {
	return db.Create(student).Error
}

// ฟังก์ชั่น ดึงข้อมูลทั้งหมดจากตาราง
func getAllStudent(db *gorm.DB) ([]Student, error) {
	//สรา้งตัวแปลเก็บข้อมูลที่ได้การดึง
	var students []Student
	result := db.Find(&students)
	return students, result.Error
}

// ฟังก์ชั่นดึงข้อมู,จากตาราง เงื่อนไข PK คือ ID
func getStudentByID(db *gorm.DB, id int) (*Student, error) {
	//สร้างตัวแปรเก้บข้อมูลที่ได้จากการดึง
	var student Student

	result := db.First(&student, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil
}

// ฟังก์ชั่นดึงข้อมู,จากตาราง เงื่อนไขไม่ใช่ PK เช่น ชื่อ ตรงกับที่ส่งมาบางส่วน
func getStudentByName(db *gorm.DB, StuName string) ([]Student, error) {
	var students []Student

	result := db.Where("stu_name LIKE ?", "%"+StuName+"%").Find(&students)

	return students, result.Error
}

func getStudentByLevel(db *gorm.DB, Stulevel int) ([]Student, error) {
	var students []Student

	result := db.Where("stu_level = ?", Stulevel).Find(&students)

	return students, result.Error
}
//ฟังก์ชั่นแก้ไขข้อมู,ในตาราง เงื่อนไข PK คือ ID
func updateStudent(db *gorm.DB, id int, student *Student) error {
	return db.Model(&Student{}).Where("id = ?", id).Updates(student).Error
}

//ฟังก์ชั่นลบข้อมู,ในตาราง เงื่อนไข PK คือ ID
func deleteStudent(db *gorm.DB, id int) error {
	return db.Delete(&Student{}, id).Error
}