package entity

type Category string

const (
	FootballCategory = "football"
	HistoryCategory  = "history"
)

func (c Category) IsValid() bool {
	switch c {

	case FootballCategory:
		return true
	}

	return false
}

func CategoryList() []Category {
	return []Category{FootballCategory, HistoryCategory}
}
