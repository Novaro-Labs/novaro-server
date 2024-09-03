package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type EventsDao struct {
	db *gorm.DB
}

func NewEventsDao(db *gorm.DB) *EventsDao {
	return &EventsDao{
		db: db,
	}
}

func (d *EventsDao) Create(r *model.Events) error {
	r.Status = 0
	tx := d.db.Create(&r)
	return tx.Error
}

func (d *EventsDao) Delete(id string) error {
	tx := d.db.Where("id = ?", id).Delete(&model.Events{})
	return tx.Error
}

func (d *EventsDao) Updates(r *model.Events) error {
	tx := d.db.Updates(&r)
	return tx.Error
}

func (d *EventsDao) Get(id string) (model.Events, error) {
	var events model.Events
	tx := d.db.Preload("Img").Where("id = ? and status = '0'", id).First(&events)

	return events, tx.Error
}

func (d *EventsDao) GetList(r *model.Events) ([]model.Events, error) {

	var events []model.Events
	tx := d.db.Where("status = '0'")

	if r.Title != "" {
		tx = tx.Where("title like ?", "%"+r.Title+"%")
	}

	tx.Preload("Img").Find(&events)
	return events, tx.Error
}
