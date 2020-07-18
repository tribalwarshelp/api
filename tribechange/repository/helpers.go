package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/tribalwarshelp/shared/models"
)

func appendTribeChangeFilterOr(or *models.TribeChangeFilterOr) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if or != nil {
			if len(or.NewTribeID) > 0 {
				q = q.WhereOr("new_tribe_id IN (?)", pg.In(or.NewTribeID))
			}
			if len(or.OldTribeID) > 0 {
				q = q.WhereOr("old_tribe_id IN (?)", pg.In(or.OldTribeID))
			}
		}
		return q, nil
	}
}
