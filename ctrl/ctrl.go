package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"gorm.io/gorm"
	"net/http"
	"smart_home/device"
	"smart_home/log"
	"smart_home/site"
	"smart_home/user"
	"time"
)

type Service interface {
	SetRouter() *gin.Engine
}

type service struct {
	User   user.Service
	Site   site.Service
	Device device.Service
	log    log.Service
}

func MustNewService(db *gorm.DB) Service {
	return &service{user.MustNewService(db), site.MustNewService(db), device.MustNewService(db), log.MustNewService(db)}
}

func RateLimit(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "rate limit!"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (t service) SetRouter() *gin.Engine {
	r := gin.Default()
	r.Use(RateLimit(time.Millisecond*200, 10, 1))
	r.POST("/smart_home/api/user/register", t.User.Register)
	r.POST("/smart_home/api/user/passwordLogin", t.User.PasswordLogin)
	r.POST("/smart_home/api/user/ping", t.User.Ping)
	AuthGroup := r.Group("/smart_home/api", t.User.Auth())
	AuthGroup.POST("/user/cookieLogin", t.User.CookieLogin)
	AuthGroup.POST("/user/getAPIKey", t.User.GetAPIKey)
	AuthGroup.POST("/site/addHouse", t.Site.AddHouse)
	AuthGroup.POST("/site/getHouse", t.Site.GetHouse)
	AuthGroup.POST("/site/deleteHouse", t.Site.DeleteHouse)
	AuthGroup.POST("/device/addDevice", t.Device.AddDevice)
	AuthGroup.POST("/device/getDevice", t.Device.GetDevice)
	AuthGroup.POST("/device/saveLayout", t.Device.SaveLayout)
	AuthGroup.POST("/device/deviceLog", t.Device.DeviceLog)
	AuthGroup.POST("/device/deleteDevice", t.Device.DeleteDevice)
	AuthGroup.POST("/log/getLog", t.log.GetLog)
	AuthGroup.POST("/log/deleteLog", t.log.DeleteLog)
	return r
}
