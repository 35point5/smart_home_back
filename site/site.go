package site

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"smart_home/db"
	"smart_home/util"
)

type Service interface {
	AddHouse(ctx *gin.Context)
	GetHouse(ctx *gin.Context)
	DeleteHouse(ctx *gin.Context)
}

type service struct {
	db *gorm.DB
}

func MustNewService(database *gorm.DB) Service {
	return &service{database}
}

func (t *service) searchSite(uid uint) []Resp {
	var res []db.Site
	t.db.Where("uid = ?", uid).Find(&res)
	var resp []Resp
	for _, o := range res {
		resp = append(resp, Resp{o.ID, o.Name, o.Img})
	}
	return resp
}

func (t *service) AddHouse(ctx *gin.Context) {
	intf, _ := ctx.Get("usr")
	usr := intf.(db.User)
	var site db.Site
	err := ctx.ShouldBind(&site)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "添加失败！"})
		return
	}
	site.Uid = usr.ID
	t.db.Create(&site)
	ctx.JSON(http.StatusOK, t.searchSite(site.Uid))
}

func (t *service) GetHouse(ctx *gin.Context) {
	intf, _ := ctx.Get("usr")
	usr := intf.(db.User)
	ctx.JSON(http.StatusOK, t.searchSite(usr.ID))
}

func (t *service) DeleteHouse(ctx *gin.Context) {
	var s db.Site
	err := ctx.ShouldBind(&s)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除设备失败！"})
		return
	}
	util.DeleteSite(t.db, s.ID)
}
