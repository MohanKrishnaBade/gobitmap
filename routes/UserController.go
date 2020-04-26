package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/gobitmap/orm"
	"github.com/gobitmap/redisStorage"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	user, err := convertResponseToStruct(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, user)
		return
	}
	user.Password = passwordHash([]byte(user.Password))
	db := orm.Connection{}.Connect()
	defer db.Close()
	user.Create(db)

	c.JSON(http.StatusOK, user)
}

func Update(c *gin.Context) {
	user, err := convertResponseToStruct(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, user)
		return
	}
	u64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
	}
	user.ID = uint(u64)
	db := orm.Connection{}.Connect()
	defer db.Close()
	user.Update(db)
	c.JSON(http.StatusOK, user)
}
func Login(c *gin.Context) {
	user, err := convertResponseToStruct(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	//db connection stuff
	db := orm.Connection{}.Connect()
	defer db.Close()

	//actual logic starts from here.
	plainPwd := user.Password
	user.SearchByEmail(db)
	if user.ID > 0 {
		if comparePasswords(user.Password, []byte(plainPwd)) {
			setUserCookie(c, user)
			c.JSON(http.StatusOK, nil)
			return
		}
	}
	c.JSON(http.StatusUnprocessableEntity, nil)
}
func FindAll(c *gin.Context) {
	db := orm.Connection{}.Connect()
	defer db.Close()
	var users []orm.User
	db.Table("user").Find(&users)
	c.JSON(http.StatusOK, users)
}

func DeleteAll(c *gin.Context) {
	db := orm.Connection{}.Connect()
	defer db.Close()

	var users []orm.User
	db.Table("user").Find(&users)
	for _, v := range users {
		v.Delete(db)
	}
	c.JSON(http.StatusOK, users)
}

func convertResponseToStruct(c *gin.Context) (orm.User, error) {
	var user orm.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		return user, nil
	}
	return user, nil
}

func passwordHash(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func setUserCookie(c *gin.Context, u orm.User) {
	// create Cookie
	sID := uuid.NewV4().String()
	c.SetCookie("gin_cookie", sID, 3600, "/", "localhost", false, true)

	rc := redisStorage.RedisConnect()
	err := rc.Set(sID, u.Email, 0).Err()
	if err != nil {
		panic(err)
	}
}

func GetUserFromCookie(c *gin.Context) {

	var user orm.User
	currentCookie, err := c.Cookie("gin_cookie")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "couldn't find the cookie",
		})
		return
	}
	rc := redisStorage.RedisConnect()
	value, err := rc.Get(currentCookie).Result()
	if err == redis.Nil || err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "couldn't find the user in redis",
		})
		return
	}
	//db connection stuff
	db := orm.Connection{}.Connect()
	defer db.Close()

	user.Email = value
	user.SearchByEmail(db)
	if user.ID > 0 {
		c.JSON(http.StatusOK, user)
		return
	}

	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"error": "user not found",
	})
}
