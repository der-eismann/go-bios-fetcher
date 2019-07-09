package lenovo

type downloads struct {
	DownloadItems []download
}

type download struct {
	Title string
	Files []file
}

type file struct {
	Name       string
	TypeString string
	Version    string
	URL        string
	Date       date
}

type date struct {
	Unix int64
}
