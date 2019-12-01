package scrapers

import "time"

type Deal struct {
	Name     string
	Link     string
	ImgLink  string
	OldPrice float64
	NewPrice float64
	Sold     int64
	Left     int64
	Start    time.Time
	End      time.Time
	Code     string
}
