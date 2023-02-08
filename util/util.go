package util

import (
	"gorm.io/gorm"
	"smart_home/db"
)

func DeleteLog(database *gorm.DB, id uint) {
	database.Unscoped().Where("id = ?", id).Delete(&db.Log{})
}

func DeleteDevice(database *gorm.DB, id uint) {
	database.Unscoped().Where("did = ?", id).Delete(&db.Log{})
	database.Unscoped().Where("id = ?", id).Delete(&db.Device{})
}

func DeleteSite(database *gorm.DB, id uint) {
	var d []db.Device
	database.Where("sid = ?", id).Find(&d)
	for _, o := range d {
		DeleteDevice(database, o.ID)
	}
	database.Unscoped().Where("id = ?", id).Delete(&db.Site{})
}
