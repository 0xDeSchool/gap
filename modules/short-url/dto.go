package short_url

type ShortUrlInput struct {
	Url string `json:"url" binding:"required"`
}

type ShortUrlOutput struct {
	Key string `json:"key"`
}
