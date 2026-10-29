// ทำ API
package main

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func studentRoute(app *fiber.App, db *gorm.DB) {
	//api post for crate (Endpoint:/api/students)
	app.Post("/api/students", func(c fiber.Ctx) error {
		//ตัวแปรเก็บข้อมูลที่ส่งมาจาก client-user
		student := new(Student)

		//นำข้อมูลที่ส่งมาเก็ฐลงในตัวแปล พร้อมทั้งตรวจ
		if error := c.Bind().Body(student); error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		//นำข้อมู,บันทึกลงตารางในdatabase พร้อมทั้งตรวจสอบ
		if error := createStudent(db, student); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		//กรณีบันทึกข้อมูลลงเรียบร้อยแล้ว
		return c.Status(201).JSON(fiber.Map{
			"message": "Student created successfully",
			"data":    student,
		})
	})

	//api Get for Read (Endpoint:/api/students)
	app.Get("/api/students", func(c fiber.Ctx) error {
		//ตัวแปรเก็บข้อมูลที่จะส่งกลับไปยังclient-user
		var students []Student

		//สรา้งตัวแปลเก็บข้อมูลที่เป้ฯเงื่อนไขที่ส่งมาแบบ query string
		StuName := c.Query("StuName")
		StuLevel := c.Query("StuLevel")

		var error error

		if StuName != "" {
			students, error = getStudentByName(db, StuName)
		} else if StuLevel != "" {
			StuLevelInt, err := strconv.Atoi(StuLevel)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{
					"message": "Invalid request parameter",
				})
			}
			students, error = getStudentByLevel(db, StuLevelInt)
		} else {
			//ดึงข้อมูลจากตาราง เก้ฐในสตัวแปร
			students, error = getAllStudent(db)
		}

		//ตรวจสอบว่ามีข้อมูลหรือไม่
		if errors.Is(error, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"message": "Student not found",
			})
		}

		//ตรวจสอบว่ามีการทำงานผิดพลาดมั้ย
		if error != nil {
			return c.Status(500).JSON(fiber.Map{
				"Massage": "Internal server error",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Success",
			"data":    students,
		})
	})

	//api Get for Read (Endpoint:/api/students/:id)
	//เงื่อนไขส่งมาแบบ path parameter ส่งต่อท้าย endpoint
	app.Get("/api/students/:id", func(c fiber.Ctx) error {
		//สร้างตัวแปรเก็บ PK คือ ID ที่ส่งมา (ใช้ Params ตัว P ใหญ่)
		id := c.Params("id")

		// แปลง id เป็นตัวเลข
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid student ID",
			})
		}

		// ดึงข้อมูลนักเรียนจาก Database ตาม ID
		student, err := getStudentByID(db, idInt)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "Student not found",
			})
		}

		// ดึงข้อมูลสำเร็จ ส่งข้อมูลกลับไป
		return c.Status(200).JSON(fiber.Map{
			"message": "Success",
			"data":    student,
		})
	})

	//API Put for Update (Endpoint:/api/students/:id)
	app.Put("/api/students/:id", func(c fiber.Ctx) error {
		// สร้างตัวเก็บ PK คือ ID ที่ส่งมา
		id := c.Params("id")

		// แปลง id เป็นตัวเลข
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		//ตัวแปรเก็บข้อมูลที่ส่งมาจาก client-user
		student := new(Student)

		//นำข้อมูลที่ส่งมาเก็บในตัวแปร พร้อมทั้งตรวจสอบ
		if error := c.Bind().Body(student); error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		// ดึงข้อมูลเก็บในตัวแปร เพื่อตรวจสอบว่ามี ID นี้ในระบบหรือไม่
		if _, error := getStudentByID(db, idInt); error != nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "Student not found",
			})
		}

		// แก้ไขข้อมูลลงตารางในฐานข้อมูล
		if error := updateStudent(db, idInt, student); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		// กรณีแก้ไขข้อมูลเรียบร้อยแล้ว
		return c.Status(200).JSON(fiber.Map{
			"message": "Student updated successfully",
			"data":    student,
		})
	})

	//API Delete for Delete (Endpoint:/api/students/:id)
	app.Delete("/api/students/:id", func(c fiber.Ctx) error {
		// สร้างตัวเก็บ PK คือ ID ที่ส่งมา
		id := c.Params("id")

		// แปลง id เป็นตัวเลข
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		//ลบข้อมูลออกจากตารางใน database พร้อมทั้งตรวจสอบ
		if error := deleteStudent(db, idInt); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		//กรณีบันทึกลงเรียบร้อยแล้ว
		return c.Status(200).JSON(fiber.Map{
			"message": "Student delete successfully",
		})
	})
}
