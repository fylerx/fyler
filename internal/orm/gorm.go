package orm

import (
	"fmt"

	"github.com/fylerx/fyler/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	dsn = fmt.Sprintf(dsn,
		cfg.DB.Host,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Database,
		cfg.DB.Port,
		cfg.DB.Sslmode,
		cfg.DB.Timezone,
	)

	return gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{PrepareStmt: false})
}
