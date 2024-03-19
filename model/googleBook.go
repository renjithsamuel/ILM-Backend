package model

type GoogleBookResponse struct {
	TotalItems int              `json:"totalItems"`
	Items      []GoogleBookItem `json:"items"`
}

type GoogleBookItem struct {
	VolumeInfo GoogleBook `json:"volumeInfo"`
}

type GoogleBook struct {
	Title               string              `json:"title"`
	Subtitle            string              `json:"subtitle,omitempty"`
	Authors             []string            `json:"authors"`
	Publisher           string              `json:"publisher"`
	PublishedDate       string              `json:"publishedDate"`
	Description         string              `json:"description"`
	IndustryIdentifiers []Identifier        `json:"industryIdentifiers"`
	PageCount           int                 `json:"pageCount"`
	PrintType           string              `json:"printType"`
	Categories          []string            `json:"categories"`
	MaturityRating      string              `json:"maturityRating"`
	AllowAnonLogging    bool                `json:"allowAnonLogging"`
	ContentVersion      string              `json:"contentVersion"`
	PanelizationSummary PanelizationSummary `json:"panelizationSummary"`
	ImageLinks          ImageLinks          `json:"imageLinks"`
	Language            string              `json:"language"`
	PreviewLink         string              `json:"previewLink"`
	InfoLink            string              `json:"infoLink"`
	CanonicalVolumeLink string              `json:"canonicalVolumeLink"`
	ReadingModes        ReadingModes        `json:"readingModes"`
}

type Identifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type PanelizationSummary struct {
	ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
	ContainsImageBubbles bool `json:"containsImageBubbles"`
}

type ImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

type ReadingModes struct {
	Text  bool `json:"text"`
	Image bool `json:"image"`
}
