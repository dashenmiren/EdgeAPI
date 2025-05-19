package acme

const DefaultProviderCode = "letsencrypt"

type Provider struct {
	Name           string `json:"name"`
	Code           string `json:"code"`
	Description    string `json:"description"`
	APIURL         string `json:"apiURL"`
	TestAPIURL     string `json:"testAPIURL"`
	RequireEAB     bool   `json:"requireEAB"`
	EABDescription string `json:"eabDescription"`
}

func FindProviderWithCode(code string) *Provider {
	for _, provider := range FindAllProviders() {
		if provider.Code == code {
			return provider
		}
	}
	return nil
}
