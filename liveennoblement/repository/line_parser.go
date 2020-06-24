package repository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type lineParser struct {
	location *time.Location
}

func newLineParser() *lineParser {
	return &lineParser{}
}

func (parser *lineParser) parse(line []string) (*models.LiveEnnoblement, error) {
	if len(line) != 4 {
		return nil, fmt.Errorf("Invalid line format (should be village_id,timestamp,new_owner,old_owner)")
	}
	var err error
	e := &models.LiveEnnoblement{}
	e.VillageID, err = strconv.Atoi(line[0])
	if err != nil {
		return nil, errors.Wrap(err, "*models.LiveEnnoblement.VillageID")
	}
	timestamp, err := strconv.Atoi(line[1])
	if err != nil {
		return nil, errors.Wrap(err, "timestamp")
	}
	e.EnnobledAt = time.Unix(int64(timestamp), 0)
	e.NewOwnerID, err = strconv.Atoi(line[2])
	if err != nil {
		return nil, errors.Wrap(err, "*models.LiveEnnoblement.NewOwnerID")
	}
	e.OldOwnerID, err = strconv.Atoi(line[3])
	if err != nil {
		return nil, errors.Wrap(err, "*models.LiveEnnoblement.OldOwnerID")
	}

	return e, nil
}
