package main

import (
	"os"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	Work		string `json:"work"  query:"work"`
	Pass		string `json:"pass"  query:"pass"`
	Genre		string `json:"genre" query:"genre"`
	Num			string `json:"num"   query:"num"`
	Flag		string `json:"flag"  query:"flag"`
}

func scoring(c echo.Context) error {
	db := db_connect()
	sqldb , _ := db.DB()
	defer sqldb.Close()


	mes := new(Message)
	if err := c.Bind(mes); err != nil {
		mes.Work = "scoring mes error!"
		return c.JSON(http.StatusOK, mes)
	}

	var quiz []Quiz
	result := db.Where("genre = ? AND num = ?", mes.Genre, mes.Num).Find(&quiz)
	if result.Error != nil {
		mes.Work = "scoring error!"
		return c.JSON(http.StatusOK, mes)
	}
	ans := "NG"
	if len(quiz) > 0 && quiz[0].Flag == mes.Flag {
		ans = "OK"
		old , _ := strconv.Atoi(quiz[0].Caught)
		new := strconv.Itoa(old + 1)
		db.Model(&Quiz{}).Where("genre = ? AND num = ?", mes.Genre, mes.Num).Update("caught", new)
	}
	return c.String(http.StatusOK, ans)
}

func show(c echo.Context) error {
	db := db_connect()
	sqldb , _ := db.DB()
	defer sqldb.Close()

	mes := new(Message)
	if err := c.Bind(mes); err != nil {
		mes.Work = "show mes error!"
		return c.JSON(http.StatusOK, mes)
	}

	var quiz []Quiz
	result := db.Find(&quiz)
	if result.Error != nil {
		mes.Work = "show error!"
		return c.JSON(http.StatusOK, mes)
	}
	return c.JSON(http.StatusOK, quiz)
}

func db_check(c echo.Context) error {
	db := db_connect()
	sqldb , err := db.DB()
	defer sqldb.Close()
	mes := new(Message)
	mes.Work = "DB conn Success!"
	if err != nil {
		mes.Work = "conn Error!"
	}
	return c.JSON(http.StatusOK, mes)
}

func table_make(c echo.Context) error {
	db := db_connect()
	sqldb , _ := db.DB()
	defer sqldb.Close()
	mes := new(Message)
	mes.Work = "table make Success!"
	result := db.Migrator().CreateTable(&Quiz{})
	if result.Error != nil {
		mes.Work = result.Error()
	}
	return c.JSON(http.StatusOK, mes)
}

func insert_row(c echo.Context) error {
	db := db_connect()
	sqldb , _ := db.DB()
	defer sqldb.Close()

	mes := new(Message)
	if err := c.Bind(mes); err != nil {
		mes.Work = "insert_row mes error!"
		return c.JSON(http.StatusOK, mes)
	}
	row := Quiz{Genre: mes.Genre, Num: mes.Num, Caught: "0", Flag: mes.Flag}
	result := db.Create(&row)
	if result.Error != nil {
		mes.Work = "insert_row error!"
	}
	return c.JSON(http.StatusOK, mes)
}

func delete_row(c echo.Context) error {
	db := db_connect()
	sqldb , _ := db.DB()
	defer sqldb.Close()

	mes := new(Message)
	if err := c.Bind(mes); err != nil {
		mes.Work = "delete_row mes error!"
		return c.JSON(http.StatusOK, mes)
	}

	var quiz []Quiz
	result := db.Where("genre = ? AND num = ?", mes.Genre, mes.Num).Delete(&quiz)
	mes.Work = "delete success!"
	if result.Error != nil {
		mes.Work = "delete error!"
	}
	return c.JSON(http.StatusOK, mes)
}

func get_row(c echo.Context) error {
	db := db_connect()
	sqldb , _ := db.DB()
	defer sqldb.Close()

	mes := new(Message)
	if err := c.Bind(mes); err != nil {
		mes.Work = "get_row mes error!"
		return c.JSON(http.StatusOK, mes)
	}

	var quiz []Quiz
	result := db.Where("genre = ? AND num = ?", mes.Genre, mes.Num).Find(&quiz)
	mes.Work = "get_row success!"
	if result.Error != nil {
		mes.Work = "get_row error!"
	}
	return c.JSON(http.StatusOK, quiz)
}

func main() {
	MY_PASS := os.Getenv("MY_PASS")
	e := echo.New()

	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})
	e.GET("/scoring", scoring)
	admin := e.Group("/admin")
	admin.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:pass",
		Validator: func(key string, c echo.Context) (bool, error) {
				  return key == MY_PASS , nil
				},
	}))
	admin.GET("/show", show)
	admin.GET("/db", db_check)
	admin.GET("/table", table_make)
	admin.GET("/insert_row", insert_row)
	admin.GET("/delete_row", delete_row)
	admin.GET("/get_row", get_row)
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
