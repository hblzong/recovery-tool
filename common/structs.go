package common

type Order struct {
}

type OrderRisk struct {
	Order
	RiskSign string
}

type OrderBusiness struct {
	Order
	BusinessSign string
}

type RHashInfo struct {
	RHashes []string
	Cursor  int
}
