package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// เรียกใช้งานฟังก์ชันเชื่อมต่อฐานข้อมูล
	db, err := ConnectDatabase()
	if err != nil {
		fmt.Println("Connect database error:", err)
		return
	}

	// Initialize Fiber
	app := fiber.New()

	// Routes Activate
	studentRoute(app, db)

	// สั่งรัน Server ที่ Port 5599
	if err := app.Listen(":5599"); err != nil {
		log.Fatal(err)
	}
}
