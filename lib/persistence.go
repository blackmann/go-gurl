package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type History struct {
	ID                uint `gorm:"primaryKey"`
	Date              time.Time
	Method            string
	Url               string
	HeadersSerialized string
	Body              string
	Annotation        string
	Status            int

	Headers map[string][]string `gorm:"-"`
}

func (h *History) BeforeCreate(tx *gorm.DB) error {
	if serialized, err := json.Marshal(h.Headers); err == nil {
		h.HeadersSerialized = string(serialized)
	}

	return nil
}

func (h *History) AfterFind(tx *gorm.DB) error {
	headers := make(map[string][]string)
	err := json.Unmarshal([]byte(h.HeadersSerialized), &headers)

	if err == nil {
		h.Headers = headers
	}

	return err
}

type Bookmark struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
	Url  string
}

type Cookie struct {
	ID  uint `gorm:"primaryKey"`
	Raw string
}

type Persistence interface {
	SaveHistory(history History)
	GetHistory() []History
	AnnotateHistory(id uint, annotation string) error
	SaveBookmark(bookmark Bookmark)
	GetBookmarks() []Bookmark
	GetBookmark(name string) (Bookmark, error)
	SaveRawCookies(rawCookie string)
	GetRawCookies() string
}

type DbPersistence struct {
	db *gorm.DB
}

func (d DbPersistence) SaveBookmark(bookmark Bookmark) {
	if d.db.Model(&bookmark).Where("name = ?", bookmark.Name).Updates(&bookmark).RowsAffected == 0 {
		d.db.Create(&bookmark)
	}
}

func (d DbPersistence) GetBookmarks() []Bookmark {
	var res []Bookmark
	d.db.Order("name asc").Find(&res)

	return res
}

func (d DbPersistence) GetBookmark(name string) (Bookmark, error) {
	var res Bookmark
	tx := d.db.Where("name = ?", name).Find(&res)

	if tx.RowsAffected == 0 {
		return Bookmark{}, errors.New(fmt.Sprintf("no bookmark with this name %s", name))
	}

	return res, nil
}

func (d DbPersistence) AnnotateHistory(id uint, annotation string) error {
	var existing History
	if tx := d.db.Where("annotation = ?", annotation).Find(&existing); tx.RowsAffected == 0 {
		// if none exists, go ahead and set it
		d.db.Model(History{}).Where("id = ?", id).Update("annotation", annotation)
	} else if existing.ID == id {
		// setting on itself so ignore
		return nil
	}

	return errors.New(fmt.Sprintf("history exists with this annotation: %s", annotation))
}

func (d DbPersistence) SaveHistory(history History) {
	d.db.Create(&history)
}

func (d DbPersistence) GetHistory() []History {
	var res []History
	d.db.Order("date desc").Find(&res)

	return res
}

func (d DbPersistence) SaveRawCookies(rawCookie string) {
	var count int64
	if d.db.Model(&Cookie{}).Count(&count); count == 0 {
		d.db.Create(&Cookie{Raw: rawCookie})
	} else {
		var last Cookie
		_ = d.db.Last(&last)
		d.db.Model(&last).Update("raw", rawCookie)
	}
}

func (d DbPersistence) GetRawCookies() string {
	var res Cookie
	d.db.First(&res)

	return res.Raw
}

func NewDbPersistence() (DbPersistence, error) {
	if db, err := gorm.Open(sqlite.Open("config.db"), &gorm.Config{}); err == nil {
		db.AutoMigrate(&History{}, &Bookmark{}, &Cookie{})

		return DbPersistence{db: db}, nil
	} else {
		return DbPersistence{}, err
	}
}
