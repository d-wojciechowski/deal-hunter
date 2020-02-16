package scrapers

func GetAllDeals() []*Deal {
	kom := ScrapXKomGroup("https://www.x-kom.pl/")
	alto := ScrapXKomGroup("https://www.al.to/")
	morele := ScrapMorele()
	return []*Deal{kom, alto, morele}
}
