package models

type StatisticFilter struct {
	UserUID      *string
	AIServiceUID *string
	Order        *string
	PageSize     *uint32
	PageNumber   *uint32
}

type StatisticRead struct {
	UserUID      string
	AIServiceUID string
	Count        uint32
	FullAmount   float64
}

type StatisticWrite struct {
	UserUID      string
	AIServiceUID string
}
