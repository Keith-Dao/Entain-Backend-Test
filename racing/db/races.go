package db

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(in *racing.ListRacesRequest) ([]*racing.Race, error)

	// Get will return a race with the given id if it exists.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(in *racing.ListRacesRequest) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, in.Filter)
	query, err = r.applyOrderBy(query, in.Sort)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

// Get will return a race with the given id if it exists.
func (r *racesRepo) Get(id int64) (*racing.Race, error) {
	var (
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]
	query, args = r.applyIdFilter(query, id)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	races, err := r.scanRaces(rows)
	if err != nil {
		return nil, err
	}
	if len(races) == 0 {
		return nil, fmt.Errorf(`A race with id "%d" does not exist.`, id)
	}

	return races[0], nil
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if filter.OnlyShowVisible {
		clauses = append(clauses, "visible = 1")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

// Applies the order by statement to the query.
// If no columns are provided in sort, sorting is defaults to advertised_start_time in ascending order.
func (r *racesRepo) applyOrderBy(query string, sort []*racing.ListRacesRequestSort) (string, error) {
	if len(sort) == 0 {
		defaultSort := racing.ListRacesRequestSort{
			Column:       "advertised_start_time",
			IsDescending: false,
		}
		sort = append(sort, &defaultSort)
	}

	var clauses []string
	isValidColumn := regexp.MustCompile("^[A-Za-z0-9_]+$")
	for _, sortDetails := range sort {
		if len(sortDetails.Column) == 0 {
			return query, errors.New("A sort request object is missing the column value.")
		}

		// Check column is valid (prevent SQL injection)
		if !isValidColumn.MatchString(sortDetails.Column) {
			return query, fmt.Errorf("%q is not a valid column name.", sortDetails.Column)
		}

		// Build order clause
		orderDirection := "ASC"
		if sortDetails.IsDescending {
			orderDirection = "DESC"
		}
		clauses = append(clauses, sortDetails.Column+" "+orderDirection)
	}

	query += " ORDER BY " + strings.Join(clauses, ",")
	return query, nil
}

// applyIdFilter applies the id filter to the query.
func (r *racesRepo) applyIdFilter(query string, id int64) (string, []interface{}) {
	return query + "WHERE id = ?", []interface{}{id}
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}
