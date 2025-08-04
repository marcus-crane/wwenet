package networkentities

type Season struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	LongDescription string          `json:"longDescription"`
	SmallCoverUrl   string          `json:"smallCoverUrl"`
	CoverUrl        string          `json:"coverUrl"`
	TitleUrl        string          `json:"titleUrl"`
	PosterUrl       string          `json:"posterUrl"`
	SeasonNumber    int             `json:"seasonNumber"`
	EpisodeCount    int             `json:"episodeCount"`
	Id              int             `json:"id"`
	Series          TruncatedSeries `json:"series"`
	Episodes        []Episode       `json:"episodes"`
	Paging          Paging          `json:"paging"`
	Licenses        []string        `json:"licenses"`
}
