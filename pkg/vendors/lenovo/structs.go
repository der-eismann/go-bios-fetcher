package lenovo

type downloads struct {
	Message string `json:"message"`
	Body    struct {
		DownloadItems []struct {
			Files []struct {
				Name       string `json:"Name"`
				TypeString string `json:"TypeString"`
				Version    string `json:"Version"`
				URL        string `json:"URL"`
				Date       struct {
					Unix int64 `json:"Unix"`
				} `json:"Date"`
			} `json:"Files"`
			Title string `json:"Title"`
		} `json:"DownloadItems"`
	} `json:"body"`
}

type windowConfig struct {
	DynamicItems struct {
		ProductID string `json:"PRODUCTID"`
	} `json:"DynamicItems"`
}
