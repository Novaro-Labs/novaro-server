package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"novaro-server/config"
	"strings"
)

type Events struct {
	Id        string `json:"id"`
	SourceId  string `json:"sourceId"`
	UserId    string `json:"userId"`
	Title     string `json:"title"`
	Website   string `json:"website"`
	ExpiredAt string `json:"expiredAt"`
	Status    int    `json:"status"`
	CreatedAt string `json:"createdAt"`
	Img       Imgs   `json:"img" gorm:"foreignKey:SourceId;references:SourceId;"`
}

func (Events) TableName() string {
	return "events"
}

func (e *Events) BeforeCreate(tx *gorm.DB) error {
	e.Id = strings.ReplaceAll(uuid.New().String(), "-", "")
	return nil
}

func (e *Events) Create(r *Events) error {
	r.Status = 0
	tx := config.DB.Create(&r)
	return tx.Error
}

func (e *Events) Delete(id string) error {
	tx := config.DB.Where("id = ?", id).Delete(&Events{})
	return tx.Error
}

func (e *Events) Update(r *Events) error {
	if r.Id == "" {
		err := e.Create(r)
		return err
	}

	tx := config.DB.Updates(&r)
	return tx.Error
}

func (e *Events) Get(id string) (Events, error) {
	var events Events
	tx := config.DB.Preload("Img").Where("id = ? and status = '0'", id).First(&events)

	return events, tx.Error
}

func (e *Events) GetList(r *Events) ([]Events, error) {

	var events []Events
	tx := config.DB.Where("status = '0'")

	if r.Title != "" {
		tx = tx.Where("title like ?", "%"+r.Title+"%")
	}

	tx.Preload("Img").Find(&events)
	return events, tx.Error
}
