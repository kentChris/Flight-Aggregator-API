package entity

type SearchCache struct {
	Flights            []Flight `json:"flights"`
	ProvidersSucceeded int      `json:"providers_succeeded"`
	ProvidersFailed    int      `json:"providers_failed"`
}
