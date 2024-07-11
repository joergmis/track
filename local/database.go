package local

import (
	"database/sql"
	"path/filepath"
	"time"

	"github.com/joergmis/track"
	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseName = "track.db"
	dateFormat   = time.RFC3339
)

type databaseStore struct {
	database *sql.DB
}

func NewDatabaseStorage(path string, version track.Version) (track.Storage, error) {
	storage := &databaseStore{}

	db, err := sql.Open("sqlite3", filepath.Join(path, databaseName))
	if err != nil {
		return storage, err
	}

	storage.database = db

	stmt := `
	create table if not exists activities (id text not null primary key, backend text, customer text, project text, service text, description text, startTime text, endTime text);
	`

	if _, err := storage.database.Exec(stmt); err != nil {
		return storage, err
	}

	return storage, nil
}

func (d *databaseStore) GetLastActivity() (track.Activity, error) {
	activity := track.Activity{}

	rows, err := d.database.Query("select * from activities order by startTime desc limit 1")
	if err != nil {
		return activity, err
	}
	defer rows.Close()

	for rows.Next() {
		var backend, start, end string
		if err := rows.Scan(&activity.ID, &backend, &activity.Customer, &activity.Project, &activity.Service, &activity.Description, &start, &end); err != nil {
			return activity, err
		}

		activity.Backend = track.BackendType(backend)

		activity.StartTime, err = time.Parse(dateFormat, start)
		if err != nil {
			return activity, err
		}

		activity.EndTime, err = time.Parse(dateFormat, end)
		if err != nil {
			return activity, err
		}
	}

	err = rows.Err()
	if err != nil {
		return activity, err
	}

	return activity, nil
}

func (d *databaseStore) GetAllActivities() ([]track.Activity, error) {
	activities := []track.Activity{}

	rows, err := d.database.Query("select * from activities")
	if err != nil {
		return activities, err
	}
	defer rows.Close()

	for rows.Next() {
		var backend, start, end string
		activity := track.Activity{}

		if err := rows.Scan(&activity.ID, &backend, &activity.Customer, &activity.Project, &activity.Service, &activity.Description, &start, &end); err != nil {
			return activities, err
		}

		activity.Backend = track.BackendType(backend)

		activity.StartTime, err = time.Parse(dateFormat, start)
		if err != nil {
			return activities, err
		}

		activity.EndTime, err = time.Parse(dateFormat, end)
		if err != nil {
			return activities, err
		}

		activities = append(activities, activity)
	}

	err = rows.Err()
	if err != nil {
		return activities, err
	}

	return activities, nil
}

func (d *databaseStore) AddActivity(activity track.Activity) error {
	tx, err := d.database.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into activities(id, backend, customer, project, service, description, startTime, endTime) values(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		activity.ID,
		activity.Backend,
		activity.Customer,
		activity.Project,
		activity.Service,
		activity.Description,
		activity.StartTime.Format(dateFormat),
		activity.EndTime.Format(dateFormat),
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *databaseStore) UpdateActivity(activity track.Activity) error {
	tx, err := d.database.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`update activities
	set
		backend=?,
		customer=?,
		project=?,
		service=?,
		description=?,
		startTime=?,
		endTime=?
	where
		id=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		activity.Backend,
		activity.Customer,
		activity.Project,
		activity.Service,
		activity.Description,
		activity.StartTime.Format(dateFormat),
		activity.EndTime.Format(dateFormat),
		activity.ID,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *databaseStore) GetUnsyncedActivities() ([]track.Activity, error) {
	return []track.Activity{}, nil
}

func (d *databaseStore) MarkActivityAsSynced(activity track.Activity) error {
	return nil
}

func (d *databaseStore) DeleteActivity(activity track.Activity) error {
	return nil
}
