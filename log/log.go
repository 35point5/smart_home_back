package log

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"smart_home/db"
	"smart_home/device"
	"smart_home/util"
)

type Service interface {
	GetLog(ctx *gin.Context)
	DeleteLog(ctx *gin.Context)
}

type service struct {
	db *gorm.DB
}

func MustNewService(database *gorm.DB) Service {
	return &service{database}
}

func record(o device.Entry) string {
	t := o.Type
	v := o.Value
	switch t {
	case 0:
		return v
	case 1:
		return v + "%"
	case 2:
		return v
	}
	return ""
}

func (t *service) GetLog(ctx *gin.Context) {
	var s db.Site
	err := ctx.ShouldBind(&s)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "获取日志失败！"})
		return
	}
	var r []db.Log
	var resp []Resp
	fmt.Println(s.ID)
	t.db.Order("updated_at desc").Where("sid = ?", s.ID).Find(&r)
	fmt.Println(r)
	for _, o := range r {
		var ma, mb map[string]device.Entry
		err := json.Unmarshal([]byte(o.StatusBefore), &mb)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "获取日志失败！"})
			return
		}
		err = json.Unmarshal([]byte(o.StatusAfter), &ma)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "获取日志失败！"})
			return
		}
		var i []string
		for k := range mb {
			if mb[k].Value != ma[k].Value {
				i = append(i, k+"从"+record(mb[k])+"变为"+record(ma[k]))
			}
		}
		resp = append(resp, Resp{o.ID, o.Name, o.UpdatedAt.Format("2006-01-02 15:04:05"), i})
	}
	ctx.JSON(http.StatusOK, resp)
}

func (t *service) DeleteLog(ctx *gin.Context) {
	var l db.Log
	err := ctx.ShouldBind(&l)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除日志失败！"})
		return
	}
	util.DeleteLog(t.db, l.ID)
}
