package types

type Metadata struct {
	Name        string `json:"name"`
	Website     string `json:"website"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Twitter     string `json:"twitter"`
	IPAddress   string `json:"ip_address"`
}

type Operator struct {
	EthAddress   string `json:"eth_address"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	LogoURI      string `json:"logo_uri"`
	IPAddress    string `json:"ip_address"`
	BtcPublicKey string `json:"btc_public_key"`
	Website      string `json:"website"`
	Twitter      string `json:"twitter"`
}
