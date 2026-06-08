package models

type Country struct {
	Name       string   `json:"name"`
	Official   string   `json:"official"`
	Slug       string   `json:"slug"`
	Capital    string   `json:"capital"`
	Region     string   `json:"region"`
	Subregion  string   `json:"subregion"`
	Population string   `json:"population"`
	Flag       string   `json:"flag"`
	Currency   string   `json:"currency"`
	Languages  []string `json:"languages"`
	Lat        float64  `json:"lat"`
	Lng        float64  `json:"lng"`
}
