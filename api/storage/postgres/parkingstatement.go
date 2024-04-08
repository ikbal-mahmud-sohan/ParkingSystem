package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"parking/api/storage"
	"parking/utility/logging"
)

const getParkingReport = `
	SELECT
		SUM(COALESCE(fee, 0)) AS total_earning,
		COUNT(*) AS total_parked,
		SUM(COALESCE(parked_hours, 0)) AS total_hours
	FROM vehicle WHERE exit IS NOT NULL 
`

func (s *Storage) GetParkingStatement(ctx context.Context, startDate, endDate time.Time) (*storage.Statement, error) {
	log := logging.FromContext(ctx).WithField("method", "GetParkingStatement")
	data, statement, arguments := storage.Statement{}, getParkingReport, []interface{}{}
	if !startDate.IsZero() {
		statement += " AND DATE(entry) >= ?"
		arguments = append(arguments, startDate)
	}
	if !endDate.IsZero() {
		statement += " AND DATE(entry) <= ?"
		arguments = append(arguments, endDate)
	}

	completeQuery, args, err := sqlx.In(statement, arguments...)
	if err != nil {
		log.WithError(err).Error("failed to sqlx")
		return nil, storage.InvalidArgument
	}

	if err := s.db.GetContext(ctx, &data, s.db.Rebind(completeQuery), args...); err != nil {
		logging.WithError(err, log).Error("failed to get statement in storage")
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return &data, nil
}
