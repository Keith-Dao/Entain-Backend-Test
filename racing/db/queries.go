package db

import "time"

const (
	racesList = "list"
)

func getRaceQueries() map[string]string {
	return map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time,
				CASE
					WHEN advertised_start_time < '` + time.Now().Format(time.RFC3339) + `' THEN 'CLOSED'
					ELSE 'OPEN'
				END as status
			FROM races
		`,
	}
}
