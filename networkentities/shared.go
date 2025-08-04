package networkentities

type Paging struct {
	MoreDataAvailable bool `json:"moreDataAvailable"`
	LastSeen          int  `json:"lastSeen"`
}
