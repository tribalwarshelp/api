package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/tribalwarshelp/shared/models"
)

func appendWhereClauseForEnnoblementFilterOr(or *models.EnnoblementFilterOr) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if or != nil {
			if len(or.NewOwnerID) > 0 {
				q = q.WhereOr("new_owner_id IN (?)", pg.In(or.NewOwnerID))
			}
			if len(or.NewOwnerTribeID) > 0 {
				q = q.WhereOr("new_owner_tribe_id IN (?)", pg.In(or.NewOwnerTribeID))
			}
			if len(or.OldOwnerID) > 0 {
				q = q.WhereOr("old_owner_id IN (?)", pg.In(or.OldOwnerID))
			}
			if len(or.OldOwnerTribeID) > 0 {
				q = q.WhereOr("old_owner_tribe_id IN (?)", pg.In(or.OldOwnerTribeID))
			}
		}
		return q, nil
	}
}
