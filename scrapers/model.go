package scrapers

import (
	"time"
)

type Deal struct {
	Id       int64
	Name     string
	Link     string
	ImgLink  string    `db:"img_link"`
	OldPrice float64   `db:"old_price"`
	NewPrice float64   `db:"new_price"`
	Sold     int64     `db:"items_sold"`
	Left     int64     `db:"items_left"`
	Start    time.Time `db:"start_date"`
	End      time.Time `db:"end_date"`
	Code     string    `db:"promo_code"`
	SiteName string    `db:"site_name"`
}

type Subscriber struct {
	Id string `db:"subscriber_chat"`
}

type Scrapper interface {
	Scrap() *Deal
}
