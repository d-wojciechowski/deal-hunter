package scrapers

func GetAllDeals() []*Deal {
	kom := ScrapXKomGroup("https://www.x-kom.pl/")
	alto := ScrapXKomGroup("https://www.al.to/")
	combat := ScrapCombat("https://www.combat.pl/")
	morele := ScrapMorele()
	return []*Deal{kom, alto, morele, combat}
}
