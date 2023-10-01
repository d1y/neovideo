package typekit

type ImportDataForm struct {
	URL  string `json:"url" form:"url"`
	Data string `json:"data" form:"data"`
}
