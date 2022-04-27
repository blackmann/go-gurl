package lib

import (
	"encoding/json"
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

type PersistenceObserver interface {
	OnChange(persistence Persistence)
}

type Persistence interface {
	SaveHistory(history History)
	GetHistory() []History
	AddListener(observer PersistenceObserver)
}

type DbPersistence struct {
	db        *gorm.DB
	observers []PersistenceObserver
}

func (d *DbPersistence) AddListener(observer PersistenceObserver) {
	d.observers = append(d.observers, observer)
}

func (d DbPersistence) SaveHistory(history History) {
	d.db.Create(&history)
}

func (d DbPersistence) GetHistory() []History {
	var res []History
	d.db.Order("date desc").Find(&res)

	return res
}

func (d DbPersistence) OnChange() {
	for _, observer := range d.observers {
		observer.OnChange(&d)
	}
}

func NewDbPersistence() (DbPersistence, error) {
	if db, err := gorm.Open(sqlite.Open("config.db"), &gorm.Config{}); err == nil {
		db.AutoMigrate(&History{})

		return DbPersistence{db: db}, nil
	} else {
		return DbPersistence{}, err
	}
}
