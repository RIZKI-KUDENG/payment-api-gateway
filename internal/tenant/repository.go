package tenant

import (
	"context"
	"payment-api-gateway/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Repository interface{
	GetTenantByID (ctx context.Context, tenantID string) (*models.Tenant, error)
}


type pgRepository struct {
	db *pgxpool.Pool
}

func (r *pgRepository) GetTenantByID(ctx context.Context, tenantID string) (*models.Tenant, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM tenants WHERE id = $1 AND status = 'active'", tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenant models.Tenant
	if rows.Next() {
		err := rows.Scan(&tenant)
		if err != nil {
			return nil, err
		}
	}
	return &tenant, nil
}
