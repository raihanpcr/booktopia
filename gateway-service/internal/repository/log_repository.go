package repository

import (
	"gateway-service/internal/model"
	"gorm.io/gorm"
)

type LogRepository interface {
	CreateLog(logEntry *model.ServiceLog)
}

type gormLogRepository struct {
	db *gorm.DB
}

func NewGormLogRepository(db *gorm.DB) LogRepository {
	return &gormLogRepository{db: db}
}

// CreateLog menyimpan log secara non-blocking
func (r *gormLogRepository) CreateLog(logEntry *model.ServiceLog) {
	go func() {
		r.db.Create(logEntry)
	}()
}