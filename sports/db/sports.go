package db

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"sports/proto/sports"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
)

// SportsRepo provides repository access to sports events.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of sports events.
	List(in *sports.ListEventsRequest) ([]*sports.Event, error)

	// Get will return a sports event with the given id if it exists.
	Get(id int64) (*sports.Event, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sport events repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sport event repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy events.
		err = s.seed()
	})

	return err
}

func (s *sportsRepo) List(in *sports.ListEventsRequest) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportsQueries()[sportsList]

	query, args = s.applyFilter(query, in.Filter)
	query, err = s.applyOrderBy(query, in.Sort)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanEvents(rows)
}

// Get will return a sport event with the given id if it exists.
func (s *sportsRepo) Get(id int64) (*sports.Event, error) {
	var (
		query string
		args  []interface{}
	)

	query = getSportsQueries()[sportsList]
	query, args = s.applyIdFilter(query, id)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	events, err := s.scanEvents(rows)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, fmt.Errorf(`A sports event with id "%d" does not exist.`, id)
	}

	return events[0], nil
}

func (s *sportsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.Sports) > 0 {
		clauses = append(clauses, "sport IN ("+strings.Repeat("?,", len(filter.Sports)-1)+"?)")

		for _, meetingID := range filter.Sports {
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
func (r *sportsRepo) applyOrderBy(query string, sort []*sports.ListEventsRequestSort) (string, error) {
	if len(sort) == 0 {
		defaultSort := sports.ListEventsRequestSort{
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
func (s *sportsRepo) applyIdFilter(query string, id int64) (string, []interface{}) {
	return query + "WHERE id = ?", []interface{}{id}
}

func (m *sportsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Sport, &event.Name, &event.Number, &event.Visible, &advertisedStart, &event.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}

	return events, nil
}
