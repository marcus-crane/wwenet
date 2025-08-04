package networkentities

type Playback struct {
	Annotations     Annotations `json:"annotations"`
	SmoothStreaming []Stream    `json:"smoothStreaming"`
	Dash            []Stream    `json:"dash"`
	HLS             []Stream    `json:"hls"`
}

type Annotations struct {
	Titles     string `json:"titles"`
	Thumbnails string `json:"thumbnails"`
}

type Stream struct {
	Subtitles []Subtitle `json:"subtitles"`
	Url       string     `json:"url"`
}

type Subtitle struct {
	Format   string `json:"format"`
	Language string `json:"language"`
	Url      string `json:"url"`
}
