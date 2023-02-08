package device

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"smart_home/db"
	"smart_home/site"
	"smart_home/util"
)

type Service interface {
	AddDevice(ctx *gin.Context)
	GetDevice(ctx *gin.Context)
	SaveLayout(ctx *gin.Context)
	DeleteDevice(ctx *gin.Context)
	DeviceLog(ctx *gin.Context)
}

type service struct {
	db *gorm.DB
}

func MustNewService(database *gorm.DB) Service {
	return &service{database}
}

func newDeng() interface{} {
	return gin.H{
		"开关": Entry{Switch, "关"},
		"亮度": Entry{Slider, "50"},
	}
}

func newKaiguan() interface{} {
	return gin.H{
		"开关": Entry{Switch, "关"},
	}
}

func newChuanganqi() interface{} {
	return gin.H{
		"温度": Entry{Text, "20℃"},
		"湿度": Entry{Text, "60%"},
	}
}

func newMen() interface{} {
	return gin.H{
		"开关": Entry{Switch, "关"},
	}
}

func deviceInit(t *db.Device) {
	t.PosX = 0
	t.PosY = 0
	t.Zoom = 1
	var b []byte
	switch t.Type {
	case Deng:
		b, _ = json.Marshal(newDeng())
	case Kaiguan:
		b, _ = json.Marshal(newKaiguan())
	case Chuanganqi:
		b, _ = json.Marshal(newChuanganqi())
	case Men:
		b, _ = json.Marshal(newMen())
	}
	t.Status = string(b)
}

func (t *service) searchDevice(uid uint) []Resp {
	var res []db.Device
	t.db.Where("sid = ?", uid).Find(&res)
	var resp []Resp
	for _, o := range res {
		var s interface{}
		err := json.Unmarshal([]byte(o.Status), &s)
		if err != nil {
			continue
		}
		resp = append(resp, Resp{o.ID, o.PosX, o.PosY, o.Zoom, o.Img, o.Name, s, o.Type, o.Sid})
	}
	return resp
}

func (t *service) AddDevice(ctx *gin.Context) {
	var device db.Device
	err := ctx.ShouldBind(&device)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "添加失败！"})
		return
	}
	deviceInit(&device)
	t.db.Create(&device)
	ctx.JSON(http.StatusOK, t.searchDevice(device.Sid))
}

func (t *service) GetDevice(ctx *gin.Context) {
	var s site.Resp
	err := ctx.ShouldBind(&s)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "查询设备失败！"})
		return
	}
	ctx.JSON(http.StatusOK, t.searchDevice(s.ID))
}

func (t *service) SaveLayout(ctx *gin.Context) {
	var rq []db.Device
	err := ctx.ShouldBind(&rq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "保存失败！"})
		return
	}
	for _, o := range rq {
		t.log(o)
	}
	for _, o := range rq {
		t.db.Omit("created_at").Save(&o)
	}
	ctx.Status(http.StatusOK)
}

func (t *service) log(d db.Device) {
	var o db.Device
	o.ID = d.ID
	t.db.First(&o)
	if o.Status != d.Status {
		r := db.Log{Sid: o.Sid, Did: o.ID, StatusBefore: o.Status, StatusAfter: d.Status, Name: o.Name}
		t.db.Create(&r)
	}
}

func (t *service) DeviceLog(ctx *gin.Context) {
	var d db.Device
	err := ctx.ShouldBind(&d)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "添加记录失败！"})
		return
	}
	var o db.Device
	t.db.First(&o, d.ID)
	o.Status = d.Status
	t.log(o)
	ctx.Status(http.StatusOK)
}

func (t *service) DeleteDevice(ctx *gin.Context) {
	var d db.Device
	err := ctx.ShouldBind(&d)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "删除设备失败！"})
		return
	}
	util.DeleteDevice(t.db, d.ID)
}
