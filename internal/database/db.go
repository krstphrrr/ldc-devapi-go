package database

import (
	"database/sql"
	"fmt"
	"go-api-app/config"

	_ "github.com/lib/pq"
)

// GetTenantDB connects to the database using the tenant's credentials.
func GetTenantDB(tenant string) (*sql.DB, error) {
	// Get the tenant credentials from the configuration
	tenantConfig, ok := config.Config.Database.Tenants[tenant]
	if !ok {
		return nil, fmt.Errorf("invalid tenant: %s", tenant)
	}

	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Name,
		tenantConfig.User,
		tenantConfig.Password,
	)

	// Open a new connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
