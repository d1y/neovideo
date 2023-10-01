package typekit

type ImportDataForm struct {
	URL  string `json:"url,omitempty" form:"url"`
	Data string `json:"data,omitempty" form:"data"`
}
