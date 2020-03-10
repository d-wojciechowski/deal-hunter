package scrapers

func GetAllDeals() []*Deal {
	return append(GetDealsAt23(), GetDealsAt1022()...)
}

func GetDealsAt1022() []*Deal {
	kom := ScrapXKomGroup("https://www.x-kom.pl/")
	alto := ScrapXKomGroup("https://www.al.to/")
	morele := ScrapMorele()
	return []*Deal{kom, alto, morele}
}

func GetDealsAt23() []*Deal {
	combat := ScrapCombat("https://www.combat.pl/")
	return []*Deal{combat}
}
