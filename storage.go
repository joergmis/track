package track

import "github.com/pkg/errors"

var (
	ErrNoActivities       = errors.New("no activities found")
	ErrNoMatchingActivity = errors.New("no matching activity found")
)

// Storage keeps a list of all activities locally available.
type ActivityRepository interface {
	// GetLastActivity checks the stored activities and returns the activity
	// which has the oldest start time. In case there is no previous data,
	// 'ErrNoActivities' is returned.
	GetLastActivity(backend ProjectBackendType) (Activity, error)

	// GetActivities returns all stored activities, sorted according to the
	// start time.
	GetActivities() (map[ProjectBackendType][]Activity, error)

	// AddActivity stores the activity.
	AddActivity(backend ProjectBackendType, activity Activity) error

	// UpdateActivity updates (replaces) the stored activity with the given
	// one. Returns [ErrNoMatchingActivity] if no activity matches the ID of
	// the given activity.
	UpdateActivity(backend ProjectBackendType, activity Activity) error

	// DeleteActivity deletes a given activity from the storage. Returns
	// [ErrNoMatchingActivity] if no activity matches the ID of the given
	// activity.
	DeleteActivity(backend ProjectBackendType, activity Activity) error
}
