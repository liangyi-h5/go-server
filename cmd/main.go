package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

var readonlyDB *sql.DB // 声明全局 db 变量
var writeDB *sql.DB    // 声明全局 db 变量

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func helloHandler(c *fiber.Ctx) error {
	data := User{
		Name: "li san",
		Age:  12,
	}
	response := Response{
		Message: "Success !",
		Code:    200,
		Data:    data,
	}
	// json.NewEncoder(w).Encode(response)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		errorResponse := Response{
			Message: "Hello, World!",
			Code:    400,
			Data:    err,
		}
		errorJsonResponse, _ := json.Marshal(errorResponse)
		errStr := string(errorJsonResponse)
		log.Error(errStr)
		c.SendString(errStr)
	}
	return c.SendString(string(jsonResponse))
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Info(err.Error())
		return
	}
	port := os.Getenv("PORT")
	log.Info(port, "port")
	readonlyDB = initDb()
	writeDB = createWriteDB()
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Get("/oms/emails", getOmsEmails)
	app.Get("/oms/email_detail", getOmsEmailDetail)
	app.Get("/oms/updata_email", updataOmsEmail)
	app.Get("/oms/set_email", setOmsEmail)
	app.Get("/oms/delete_email", deleteOmsEmailByEmail)
	app.Get("/", helloHandler)
	app.Use(func(c *fiber.Ctx) error {
		// 404处理
		errorJsonResponse, _ := json.Marshal(Response{
			Message: "Sorry can't find that!",
			Code:    fiber.StatusNotFound,
			Data:    nil,
		})
		log.Error(c.Context().URI())
		return c.Status(fiber.StatusNotFound).SendString(string(errorJsonResponse))
	})
	app.Listen(":" + port)
}
