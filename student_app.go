package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"
)

type student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Marks int    `json:"marks"`
}

func GetMySQLDB() *sql.DB {
	var db *sql.DB
	var server = "sql-db-student.database.windows.net"
	var port = 1433
	var user = "azureuser"
	var password = "Sh@nmukh1234"
	var database = "student-db"

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error

	db, err = sql.Open("sqlserver", connString)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getstudents(context *gin.Context) {
	students := []student{}
	std_db := GetMySQLDB()
	defer std_db.Close()
	rows, err := std_db.Query("select * from student")
	if err != nil {
		log.Fatal(err)
	} else {
		for rows.Next() {
			s := student{}
			rows.Scan(&s.ID, &s.Name, &s.Marks)
			students = append(students, s)
		}
	}
	context.IndentedJSON(http.StatusOK, students)
}

func addstudent(context *gin.Context) {
	std_db := GetMySQLDB()
	defer std_db.Close()
	s := student{}
	context.BindJSON(&s)
	_, err := std_db.Exec(fmt.Sprintf("insert into student values(%d, '%s', %d)", s.ID, s.Name, s.Marks))
	if err != nil {
		log.Fatal(err)
	}
	context.IndentedJSON(http.StatusCreated, s)
}

func Get_student_by_id(context *gin.Context) {
	id := context.Param("id")
	s := student{}
	std_db := GetMySQLDB()
	defer std_db.Close()
	rows, err := std_db.Query(fmt.Sprintf("select * from student where id=%s", id))
	if err != nil {
		log.Fatal(err)
	} else {
		rows.Next()
		rows.Scan(&s.ID, &s.Name, &s.Marks)
	}
	context.IndentedJSON(http.StatusOK, s)
}

func update_student(context *gin.Context) {
	std_db := GetMySQLDB()
	defer std_db.Close()
	s := student{}
	context.BindJSON(&s)
	result, err := std_db.Exec(fmt.Sprintf("UPDATE student SET name = '%s', marks = %d WHERE id=%d", s.Name, s.Marks, s.ID))
	if err != nil {
		log.Fatal(err)
	} else {
		r, _ := result.RowsAffected()
		if r == 0 {
			context.IndentedJSON(http.StatusNotModified, "{error: student not updated}")
		} else {
			context.IndentedJSON(http.StatusCreated, s)
		}
	}

}

func deletestudent(context *gin.Context) {
	std_db := GetMySQLDB()
	defer std_db.Close()
	ID, _ := strconv.Atoi(string(context.Param("id")))
	result, err := std_db.Exec(fmt.Sprintf("delete student where id=%d", ID))
	if err != nil {
		log.Fatal(err)
	} else {
		r, _ := result.RowsAffected()
		if r == 0 {
			context.IndentedJSON(http.StatusNotModified, "{error: student not deleted}")
		} else {
			context.IndentedJSON(http.StatusOK, "{message: student deleted}")
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/students", getstudents)
	router.GET("/students/:id", Get_student_by_id)
	router.PUT("/students/:id", update_student)
	router.POST("/students", addstudent)
	router.DELETE("/students/:id", deletestudent)
	router.Run(":9090")

}
