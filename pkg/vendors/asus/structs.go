package asus

type downloadsASUS struct {
	Result struct {
		Objects []struct {
			Name  string `json:"Name"`
			Files []struct {
				Version     string `json:"Version"`
				Title       string `json:"Title"`
				Description string `json:"Description"`
				FileSize    string `json:"FileSize"`
				ReleaseDate string `json:"ReleaseDate"`
				DownloadURL struct {
					Global string `json:"Global"`
				} `json:"DownloadUrl"`
			} `json:"Files"`
		} `json:"Obj"`
	} `json:"Result"`
}
