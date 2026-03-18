package commandresponse

type (
	Message struct {
		ResponseType string `json:"response_type"`
		Blocks       []any  `json:"blocks"`
	}

	Section struct {
		Type   string `json:"type"`
		Text   Text   `json:"text,omitempty,omitzero"`
		Fields []any  `json:"fields,omitempty"`
	}

	Text struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}

	Image struct {
		Type     string `json:"type"`
		ImageURL string `json:"image_url"`
		AltText  string `json:"alt_text"`
	}
)
