package networkentities

type Playlist struct {
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CoverUrl      string    `json:"coverUrl"`
	SmallCoverUrl string    `json:"smallCoverUrl"`
	PlaylistType  string    `json:"playlistType"`
	Id            int       `json:"id"`
	VODs          []Episode `json:"vods"`
	Paging        Paging    `json:"paging"`
}
