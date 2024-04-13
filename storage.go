package track

// Storage keeps a list of all activities locally available.
type Storage interface {
	GetActivities() ([]Activity, error)
	AddActivity(activity Activity) error
}
