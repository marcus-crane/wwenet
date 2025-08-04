package networkentities

type Series struct {
	Id              int      `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	LongDescription string   `json:"longDescription"`
	SmallCoverUrl   string   `json:"smallCoverUrl"`
	CoverUrl        string   `json:"coverUrl"`
	TitleUrl        string   `json:"titleUrl"`
	PosterUrl       string   `json:"posterUrl"`
	LogoUrl         string   `json:"logoUrl"`
	Seasons         []Season `json:"seasons"` // Truncated
	Paging          Paging   `json:"paging"`
}

// Embedded in Season payload
type TruncatedSeries struct {
	SeriesId        int    `json:"seriesId"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	LongDescription string `json:"longDescription"`
}
