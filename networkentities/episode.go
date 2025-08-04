package networkentities

type Episode struct {
	Type               string             `json:"vod"`
	Id                 int                `json:"id"`
	Title              string             `json:"title"`
	Description        string             `json:"description"`
	ThumbnailUrl       string             `json:"thumbnailUrl"`
	PosterUrl          string             `json:"posterUrl"`
	CoverUrl           string             `json:"coverUrl"`
	Duration           int                `json:"duration"`
	ExternalAssetId    string             `json:"externalAssetId"`
	Favorite           bool               `json:"favourite"`
	EpisodeInformation EpisodeInformation `json:"episodeInformation"` // When embedded in season payload
	PlayerUrlCallback  string             `json:"playerUrlCallback"`
	OnlinePlayback     string             `json:"onlinePlayback"`
}

type Rating struct {
	Rating      string   `json:"rating"`
	Descriptors []string `json:"descriptors"`
}

type EpisodeInformation struct {
	SeasonNumber  int    `json:"seasonNumber"`
	EpisodeNumber int    `json:"episodeNumber"`
	Season        int    `json:"season"`
	SeasonTitle   string `json:"seasonTitle"`
}
