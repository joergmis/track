package track

import "github.com/pkg/errors"

var (
	ErrNoActivities       = errors.New("no activities found")
	ErrNoMatchingActivity = errors.New("no matching activity found")
)

// Storage keeps a list of all activities locally available.
type Storage interface {
	// GetLastActivity checks the stored activities and returns the activity
	// which has the oldest start time. In case there is no previous data,
	// 'ErrNoActivities' is returned.
	GetLastActivity() (Activity, error)

	// GetActivities returns all stored activities, sorted according to the
	// start time.
	GetActivities() ([]Activity, error)

	// AddActivity stores the activity.
	AddActivity(activity Activity) error

	// UpdateActivity updates (replaces) the stored activity with the given
	// one. Returns [ErrNoMatchingActivity] if no activity matches the ID of
	// the given activity.
	UpdateActivity(activity Activity) error
}
