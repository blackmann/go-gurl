package lib

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type History struct {
	ID     uint
	Date   time.Time
	Method string
	Url    string
	//Headers    map[string]string
	Body       string
	Annotation string
}

type Persistence interface {
	SaveHistory(history History)
	GetHistory() []History
}

type DbPersistence struct {
	db *gorm.DB
}

func (d DbPersistence) SaveHistory(history History) {
	d.db.Create(&history)
}

func (d DbPersistence) GetHistory() []History {
	var res []History
	d.db.Order("date desc").Find(&res)

	return res
}

func NewDbPersistence() (DbPersistence, error) {
	if db, err := gorm.Open(sqlite.Open("config.db"), &gorm.Config{}); err == nil {
		db.AutoMigrate(&History{})

		return DbPersistence{db: db}, nil
	} else {
		return DbPersistence{}, err
	}
}
