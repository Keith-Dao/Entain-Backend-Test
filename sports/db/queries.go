package db

import "time"

const (
	sportsList = "list"
)

func getSportsQueries() map[string]string {
	return map[string]string{
		sportsList: `
			SELECT 
				id, 
				sport, 
				name, 
				number, 
				visible, 
				advertised_start_time,
				CASE
					WHEN advertised_start_time < '` + time.Now().Format(time.RFC3339) + `' THEN 'CLOSED'
					ELSE 'OPEN'
				END as status
			FROM events
		`,
	}
}
