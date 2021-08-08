package route

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"srv/src/module/hero"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type shuju struct {
	UserName string `json:"userName"`
	UserPwd  string `json:"userPwd"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func initMySQL(databaseName string) (err error) {
	dsn := fmt.Sprintf("root:1234@tcp(localhost:3306)/%s?charset=utf8&parseTime=True&loc=Local", databaseName)
	db, err = sqlx.Open("mysql", dsn)
	return err
}

func createUserTable() {
	schema := `CREATE TABLE  if not exists user (
	id int AUTO_INCREMENT primary key NOT NULL ,
  name varchar(50) unique NOT NULL,
  password varchar(20) NOT NULL)
	;`
	// 调用Exec函数执行sql语句，创建表
	_, err := db.Exec(schema)
	//错误处理
	if err != nil {
		panic(err)
	}
}

func insertUser(userName string, userPwd string) (res int) {

	judgeNameSql := `SELECT COUNT(*) FROM user WHERE name=?`
	var total int
	db.Get(&total, judgeNameSql, userName)
	if total == 1 {
		return 0
	}
	insertSql := `INSERT INTO user (name,password) VALUES(?,?)`
	db.MustExec(insertSql, userName, userPwd)
	return 9
}

func judgeUser(userName string, userPwd string) (res int) {
	judgeNameSql := `SELECT COUNT(*) FROM user WHERE name=?`
	var total int
	db.Get(&total, judgeNameSql, userName)
	if total == 0 {
		return 0
	}
	judgePwdSql := `SELECT COUNT(*) FROM user WHERE name=? AND password=?`
	db.Get(&total, judgePwdSql, userName, userPwd)
	if total == 0 {
		return 1
	} else {
		return 9
	}
}

// Init ...
func Init() *gin.Engine {
	// Init database
	if initMySQL("proj") != nil {
		println("error")
		return nil
	} else {
		println("Connect DAtabases success")
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)

	createUserTable()

	// Init route
	r := gin.Default()

	rootRoute := r.Group("/")
	{

		rootRoute.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/demo/proj_out/")
		})

		rootRoute.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong12",
			})
		})
	}

	// /api
	apiRoute := r.Group("/api")
	{
		apiRoute.GET("/ListHero", hero.SrvListHero)
	}

	// /demo
	projRoute := rootRoute.Group("/demo")
	projRoute.StaticFS("/proj_out", http.Dir("./proj_out"))
	projRoute.StaticFS("/heroes_out", http.Dir("./heroes_out"))

	projApiRoute := projRoute.Group("/api")
	{
		projApiRoute.GET("/user/login", func(c *gin.Context) {

			name := c.DefaultQuery("name", "")
			pwd := c.DefaultQuery("pwd", "")
			// log.Println(name)
			// log.Println(pwd)
			var resStr string
			if res := judgeUser(name, pwd); res != 9 {
				if res == 0 {
					resStr = "Wrong Account"
				} else {
					resStr = "Wrong Password"
				}
				log.Println("rseStr", resStr)
				c.JSON(200, gin.H{
					"Login": resStr,
				})
				return
			}
			log.Println("account-password corret")
			resStr = "account-password corret"
			token := RandStringRunes(16)
			log.Printf("userLogin token: %s", token)
			// OK
			c.SetCookie("userSession", token, 3600, "/", "", false, true)

			c.JSON(200, gin.H{
				"Login": resStr,
			})
		})

		projApiRoute.POST("/user/register", func(c *gin.Context) {
			var userLogin shuju
			var resStr string
			if err := c.ShouldBindJSON(&userLogin); err != nil {
				log.Printf("Error: %s", err.Error())
				return
			}
			if res := insertUser(userLogin.UserName, userLogin.UserPwd); res != 9 {
				resStr = "Insert Fail"
				return
			} else {
				resStr = "Insert seccess"
			}

			c.JSON(200, gin.H{
				"Login": resStr,
			})
		})

	}

	return r
}