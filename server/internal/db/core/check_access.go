package core

import (
	"context"
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/models"
)

func (r *BaseRepo[T]) CheckAccess(ctx context.Context, uc *models.UserContext, id any) (bool, error) {
	switch uc.Role {
	case models.RoleAdmin:
		return true, nil // admin can access everything
	case models.RoleUser:
		// assume tables have owner_id column to restrict ownership
		query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE id = ? AND owner_id = ?", r.table)
		var count int
		err := r.GetExec().GetContext(ctx, &count, query, id, uc.ID)
		if err != nil {
			return false, err
		}
		return count > 0, nil
	case models.RoleGuest:
		return false, nil // guests can access nothing
	default:
		return false, fmt.Errorf("unknown role: %d", uc.Role)
	}
}
