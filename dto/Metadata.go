package dto

type MetadataExternal struct {
	Metadata Metadata `json:"metadata"`
}
type Metadata struct {
	OgTitle       string `json:"ogTitle"`
	OgDescription string `json:"ogDescription"`
	Domain        string `json:"domain"`
	FullURL       string `json:"fullUrl"`
	Image         string `json:"image"`
}

type MetadataRequest struct {
	Url string `json:"url"`
}
