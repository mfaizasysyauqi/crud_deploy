package main

import (
	// "database/sql"
	// "reflect"
	"fmt"
	"log"
	"net/http"
	"os"

	"simple-api/auth"
	"simple-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq" // add this
)

type newStudent struct {
	Student_id       uint64 `json:"student_id" binding:"required"`
	Student_name     string `json:"student_name" binding:"required"`
	Student_age      uint64 `json:"student_age" binding:"required"`
	Student_address  string `json:"student_address" binding:"required"`
	Student_phone_no string `json:"student_phone_no" binding:"required"`
}

// func rowToStruct(rows *sql.Rows, dest interface{}) error {
// 	destv := reflect.ValueOf(dest).Elem()

// 	args := make([]interface{}, destv.Type().Elem().NumField())

// 	for rows.Next() {
// 		rowp := reflect.New(destv.Type().Elem())
// 		rowv := rowp.Elem()

// 		for i := 0; i < rowv.NumField(); i++ {
// 			args[i] = rowv.Field(i).Addr().Interface()
// 		}

// 		if err := rows.Scan(args...); err != nil {
// 			return err
// 		}

// 		destv.Set(reflect.Append(destv, rowv))
// 	}

// 	return nil
// }

func postHandler(c *gin.Context, db *gorm.DB) {
	// if c.Bind(&newStudent) == nil {
	// 	_, err := db.Exec("insert into students values ($1,$2,$3,$4,$5)", newStudent.Student_id, newStudent.Student_name, newStudent.Student_age, newStudent.Student_address, newStudent.Student_phone_no)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"message": "success create"})
	// }

	// c.JSON(http.StatusBadRequest, gin.H{"message": "error"})

	// ==================================================================================

	var newStudent newStudent
	c.Bind(&newStudent)
	db.Create(&newStudent)
	c.JSON(http.StatusOK, gin.H{"message": "success create", "data": newStudent})
}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	// row, err := db.Query("select * from students")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// }

	// rowToStruct(row, &newStudent)

	// if newStudent == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"data": newStudent})

	// ==================================================================================
	var newStudent []newStudent
	db.Find(&newStudent)
	c.JSON(http.StatusOK, gin.H{"message": "succes find all", "data": newStudent})

}

func getHandler(c *gin.Context, db *gorm.DB) {
	// var newStudent []newStudent

	// row, err := db.Query("select * from students where student_id = $1", studentId)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// }

	// rowToStruct(row, &newStudent)

	// if newStudent == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"data": newStudent})

	// ==================================================================================

	var newStudent newStudent

	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success find by id", "data": newStudent})
}

func putHandler(c *gin.Context, db *gorm.DB) {
	// var newStudent newStudent

	// studentId := c.Param("student_id")

	// if c.Bind(&newStudent) == nil {
	// 	_, err := db.Exec("update students set student_name=$1 where student_id=$2", newStudent.Student_name, studentId)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"message": "success update"})
	// }

	// ==================================================================================

	var newStudent newStudent

	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id=?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}
	var reqStudent = newStudent

	c.Bind(&reqStudent)

	db.Model(&newStudent).Update(reqStudent)

	c.JSON(http.StatusOK, gin.H{
		"message": "success update",
		"data":    reqStudent,
	})
}

func delHandler(c *gin.Context, db *gorm.DB) {
	// studentId := c.Param("student_id")

	// _,err := db.Exec("delete from students where student_id=$1",studentId)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"message": "success delete"})

	// ==================================================================================

	var newStudent newStudent

	studentId := c.Param("student_id")
	db.Delete(&newStudent, "student_id=?", studentId)

	c.JSON(http.StatusOK, gin.H{
		"message": "success delete",
	})
}

func setupRouter() *gin.Engine {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error load env")
	}

	conn := os.Getenv("POSTGRES_URL")
	gormDB, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	// Pass the gorm.DB instance to your Migrate function
	Migrate(gormDB)
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	r.POST("/login", auth.LoginHandler)

	r.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, gormDB)
	})

	r.GET("/student", middleware.AuthValid, func(ctx *gin.Context) {
		getAllHandler(ctx, gormDB)
	})

	r.GET("/student/:student_id", middleware.AuthValid, func(ctx *gin.Context) {
		getHandler(ctx, gormDB)
	})

	r.PUT("/student/:student_id", func(ctx *gin.Context) {
		putHandler(ctx, gormDB)
	})

	r.DELETE("/student/:student_id", func(ctx *gin.Context) {
		delHandler(ctx, gormDB)
	})

	return r
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&newStudent{})

	data := newStudent{}
	if db.Find(&data).RecordNotFound() {
		fmt.Println("======================== run seeder user ========================")
		seederUser(db)
	}
}

func seederUser(db *gorm.DB) {
	data := newStudent{
		Student_id:       1,
		Student_name:     "Yayat",
		Student_age:      20,
		Student_address:  "Jakarta",
		Student_phone_no: "012345689",
	}

	db.Create(&data)
}

func main() {
	r := setupRouter()

	r.Run()

}
