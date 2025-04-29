package messenger

type Alerter interface {
	AlertBudgetClose(categoryID uint)
	AlertBudgetReached(categoryID uint)
}
