package route

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"srv/src/module/hero"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type heroInfo struct {
	HeroId   int    `json:"id" db:"id"`
	HeroName string `json:"name" db:"name"`
}

type userInfo struct {
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

func createHeroTable() {
	schema := `CREATE TABLE  if not exists hero (
	id int AUTO_INCREMENT primary key unique NOT NULL ,
  name varchar(50) NOT NULL)
	;`
	// 调用Exec函数执行sql语句，创建表
	_, err := db.Exec(schema)
	//错误处理
	if err != nil {
		panic(err)
	}
}

func createUserTable() {
	schema := `CREATE TABLE  if not exists user (
	id int AUTO_INCREMENT primary key NOT NULL ,
  name varchar(50) unique NOT NULL,
  password varchar(20) NOT NULL,
	status varchar(20) NOT NULL DEFAULT 'activation')
	;`
	// 调用Exec函数执行sql语句，创建表
	_, err := db.Exec(schema)
	//错误处理
	if err != nil {
		panic(err)
	}
}

func insertHero(heroName string) int {

	judgeNameSql := `SELECT COUNT(*) FROM hero WHERE name=?`
	var total int
	db.Get(&total, judgeNameSql, heroName)
	if total == 1 {
		return 0
	}
	insertSql := `INSERT INTO hero (name) VALUES(?)`
	db.MustExec(insertSql, heroName)
	return 9
}
func insertUser(userName string, userPwd string) int {

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
func getHeroAll() []heroInfo {
	var res []heroInfo
	getHeroSql := `SELECT * FROM hero`
	if err := db.Select(&res, getHeroSql); err != nil {
		var errHero []heroInfo
		return errHero
	}
	return res
}
func getHero(heroID int) heroInfo {
	var res heroInfo
	getHeroSql := `SELECT id,name FROM hero WHERE id=? `
	if err := db.Get(&res, getHeroSql, heroID); err != nil {
		var errHero heroInfo
		return errHero
	}

	return res
}

func updateHero(heroId int, heroName string) int {
	updateHeroSql := `UPDATE hero SET name=? WHERE id =?`
	ret, err := db.Exec(updateHeroSql, heroName, heroId)
	if err != nil {
		return 0
	}
	n, err := ret.RowsAffected()
	if err != nil || n == 0 {
		return 0
	}
	return 9
}

func delHero(heroId int) int {
	delHeroSql := `DELETE FROM hero WHERE id=?`
	_, err := db.Exec(delHeroSql, heroId)
	if err != nil {
		return 0
	}
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
	createHeroTable()
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
		projApiRoute.POST("/user/login", func(c *gin.Context) {

			var resStr string
			var userLogin userInfo
			if err := c.ShouldBindJSON(&userLogin); err != nil {
				log.Printf("Error: %s", err.Error())
				return
			}

			ok, err := regexp.MatchString("^[a-zA-Z0-9]+$", userLogin.UserName)
			log.Println("ok,err", ok, err)
			if ok {
				ok, _ = regexp.MatchString("^[a-zA-Z0-9]+$", userLogin.UserPwd)
				if ok {
					strBytes := []byte(userLogin.UserPwd)
					enconded := base64.StdEncoding.EncodeToString(strBytes)
					if res := judgeUser(userLogin.UserName, enconded); res != 9 {
						if res == 0 {
							resStr = "Wrong Account"
						} else {
							resStr = "Wrong Password"
						}
					}
				} else {
					resStr = "Illegal Password"
				}
			} else {
				resStr = "Illegal Name"
			}
			if resStr != "" {
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
			var userLogin userInfo
			var resStr string
			if err := c.ShouldBindJSON(&userLogin); err != nil {
				log.Printf("Error: %s", err.Error())
				return
			}
			ok, _ := regexp.MatchString("^[a-zA-Z0-9]+$", userLogin.UserName)
			if ok {
				ok, _ = regexp.MatchString("^[a-zA-Z0-9]+$", userLogin.UserPwd)
				if ok {
					strBytes := []byte(userLogin.UserPwd)
					enconded := base64.StdEncoding.EncodeToString(strBytes)

					if res := insertUser(userLogin.UserName, enconded); res != 9 {
						resStr = "Insert Fail"
					} else {
						resStr = "Insert seccess"
					}
				} else {
					resStr = "Illegal Password"
				}
			} else {
				resStr = "Illegal Name"
			}

			c.JSON(200, gin.H{
				"Register": resStr,
			})
		})

		//hero
		projApiRoute.GET("/hero", func(c *gin.Context) {

			strId := c.DefaultQuery("id", "-1")
			var resHero heroInfo
			intId, err := strconv.Atoi(strId)
			if err == nil {
				resHero = getHero(intId)
			}

			c.JSON(200, gin.H{
				"Res": resHero,
			})
		})
		//

		projApiRoute.GET("/hero/getall", func(c *gin.Context) {
			resHero := getHeroAll()

			c.JSON(200, gin.H{
				"Res": resHero,
			})
		})

		projApiRoute.POST("/hero", func(c *gin.Context) {

			var hero heroInfo
			if err := c.ShouldBindJSON(&hero); err != nil {
				log.Printf("Error: %s", err.Error())
				return
			}
			var resStr string
			ok, _ := regexp.MatchString("^[a-zA-Z0-9]+$", hero.HeroName)
			if ok {
				intRes := insertHero(hero.HeroName)

				if intRes != 9 {
					resStr = "Fail"
				} else {
					resStr = "Success"
				}
			} else {
				resStr = "Illegal Name"
			}
			c.JSON(200, gin.H{
				"Res": resStr,
			})
		})

		projApiRoute.DELETE("/hero", func(c *gin.Context) {
			resStr := "Fail"
			heroId := c.DefaultQuery("id", "-1")
			intId, err := strconv.Atoi(heroId)
			if err == nil {
				intRes := delHero(intId)
				if intRes == 9 {
					resStr = "Success"
				}
			}
			c.JSON(200, gin.H{
				"Res": resStr,
			})
		})

		projApiRoute.PUT("/hero", func(c *gin.Context) {
			resStr := "Fail"
			var hero heroInfo
			if err := c.ShouldBindJSON(&hero); err != nil {
				log.Printf("Error: %s", err.Error())
				return
			}
			ok, _ := regexp.MatchString("^[a-zA-Z0-9]+$", hero.HeroName)
			if ok {

				intRes := updateHero(hero.HeroId, hero.HeroName)
				if intRes == 9 {
					resStr = "Success"
				}
			} else {
				resStr = "Illegal Name"
			}
			c.JSON(200, gin.H{
				"Res": resStr,
			})
		})
	}
	return r
}
