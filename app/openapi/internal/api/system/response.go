package system

type currentTimeEntity struct {
	Zone      string `json:"zone"`
	Time      string `json:"time"`
	Timestamp int64  `json:"timestamp"`
}
