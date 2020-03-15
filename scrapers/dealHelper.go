package scrapers

func GetAllDeals() []*Deal {
	return append(GetDealsAt23(), GetDealsAt1022()...)
}

func GetFlexibleDeals() []*Deal {
	return []*Deal{ScrapMorele()}
}

func GetDealsAt1022() []*Deal {
	kom := ScrapXKomGroup("https://www.x-kom.pl/")
	alto := ScrapXKomGroup("https://www.al.to/")
	return []*Deal{kom, alto}
}

func GetDealsAt23() []*Deal {
	combat := ScrapCombat("https://www.combat.pl/")
	return []*Deal{combat}
}
