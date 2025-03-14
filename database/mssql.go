package database

import (
	"fmt"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func NewMsSqlConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
