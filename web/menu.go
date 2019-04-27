package web

type Menu struct {
	DefaultTitle string
	ItemGroups   [][]*MenuItem
}

func (m Menu) GetActiveTitle() string {
	for _, itemGroup := range m.ItemGroups {
		for _, item := range itemGroup {
			if item.Active {
				return item.Title
			}
		}
	}
	return m.DefaultTitle
}

func (m Menu) SetActive(url string) {
	for _, itemGroup := range m.ItemGroups {
		for _, item := range itemGroup {
			if item.Url == url {
				item.Active = true
			}
		}
	}
}

type MenuItem struct {
	Url    string
	Title  string
	Active bool
}
