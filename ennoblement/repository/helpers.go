package repository

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

func getCSVData(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return csv.NewReader(resp.Body).ReadAll()
}

func parseLine(line []string, timezone string) (*models.Ennoblement, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("Invalid timezone: %s", timezone)
	}
	if len(line) != 4 {
		return nil, fmt.Errorf("Invalid line format (should be village_id,timestamp,new_owner,old_owner)")
	}
	e := &models.Ennoblement{}
	e.VillageID, err = strconv.Atoi(line[0])
	if err != nil {
		return nil, errors.Wrap(err, "*models.Ennoblement.VillageID")
	}
	timestamp, err := strconv.Atoi(line[1])
	if err != nil {
		return nil, errors.Wrap(err, "timestamp")
	}
	e.EnnobledAt = time.Unix(int64(timestamp), 0).In(loc)
	e.NewOwnerID, err = strconv.Atoi(line[2])
	if err != nil {
		return nil, errors.Wrap(err, "*models.Ennoblement.NewOwnerID")
	}
	e.OldOwnerID, err = strconv.Atoi(line[2])
	if err != nil {
		return nil, errors.Wrap(err, "*models.Ennoblement.OldOwnerID")
	}

	return e, nil
}
