package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"smart_home/db"
	"strconv"
	"time"
)

type Service interface {
	Register(c *gin.Context)
	CookieLogin(c *gin.Context)
	PasswordLogin(c *gin.Context)
	GetAPIKey(ctx *gin.Context)
	Ping(ctx *gin.Context)
	Auth() gin.HandlerFunc
}

type service struct {
	db *gorm.DB
}

func MustNewService(database *gorm.DB) Service {
	return &service{database}
}

func updateCookie(usr *db.User) {
	h := md5.New()
	h.Write([]byte(usr.Name + strconv.FormatInt(time.Now().Unix(), 10) + "mogician"))
	usr.Cookie = hex.EncodeToString(h.Sum(nil))
}

func generateAPIKey(usr db.User) string {
	h := md5.New()
	h.Write([]byte(usr.Phone + strconv.FormatInt(time.Now().Unix(), 10) + "114514"))
	return hex.EncodeToString(h.Sum(nil))
}

func passwordEncrypt(pwd string) string {
	h := md5.New()
	h.Write([]byte(pwd + "1919810"))
	return hex.EncodeToString(h.Sum(nil))
}

func (t *service) UpdateInfo(usr *db.User) error {
	updateCookie(usr)
	res := t.db.Save(&usr)
	return res.Error
}

func checkName(s string) bool {
	return len(s) >= 6
}

func checkPassword(s string) bool {
	return len(s) >= 6
}

func checkPhone(s string) bool {
	if len(s) != 11 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (t *service) Register(c *gin.Context) {
	var usr db.User
	err := c.ShouldBind(&usr)
	fmt.Println(usr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "注册失败！"})
		return
	}
	if !checkName(usr.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "注册失败，用户名格式错误！"})
		return
	}
	if !checkPassword(usr.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "注册失败，密码格式错误！"})
		return
	}
	if !checkPhone(usr.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "注册失败，手机号格式错误！"})
		return
	}
	usr.Password = passwordEncrypt(usr.Password)
	res := t.db.Create(&usr)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "注册失败，该用户名或手机号已存在！"})
		return
	}
	err = t.UpdateInfo(&usr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "注册失败！"})
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("usr", usr.Cookie, 3600, "/", "mogician.cc", false, true)
	c.Status(http.StatusOK)
}

func (t *service) CookieLogin(c *gin.Context) {
	intf, _ := c.Get("usr")
	usr := intf.(db.User)
	c.JSON(http.StatusOK, gin.H{"Name": usr.Name})
}

func (t *service) PasswordLogin(c *gin.Context) {
	var usr, data db.User
	err := c.ShouldBind(&usr)
	if err != nil || usr.Name == "" || usr.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "登录失败！"})
		return
	}
	res := t.db.First(&data, "name = ?", usr.Name)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "登录失败，该用户不存在！"})
		return
	}
	usr.Password = passwordEncrypt(usr.Password)
	fmt.Println(usr.Password)
	if data.Password != usr.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "登录失败，密码错误！"})
		return
	}
	updateCookie(&data)
	t.db.Save(&data)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("usr", data.Cookie, 3600, "/", "mogician.cc", false, true)
	c.Status(http.StatusOK)
}

func (t *service) GetAPIKey(ctx *gin.Context) {
	intf, _ := ctx.Get("usr")
	usr := intf.(db.User)
	usr.Key = generateAPIKey(usr)
	t.db.Save(&usr)
	ctx.JSON(http.StatusOK, gin.H{"Key": usr.Key})
}

func (t *service) Ping(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (t *service) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("usr")
		if err != nil {
			k := c.Query("key")
			if k != "" {
				var usr db.User
				res := t.db.Where("`key` = ?", k).First(&usr)
				if res.Error == nil {
					c.Set("usr", usr)
					c.Next()
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "请先登录或注册！"})
			return
		}
		var usr = db.User{Cookie: cookie}
		res := t.db.Where("cookie = ?", cookie).First(&usr)
		if res.Error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "登录过期，请刷新！"})
			return
		}
		c.Set("usr", usr)
		fmt.Println(usr.Cookie)
		c.Next()
	}
}
