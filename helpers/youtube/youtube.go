package youtube

import (
	"encoding/xml"
	"github.com/curt-labs/GoAdmin/helpers/rest"
	//"log"
	"math"
	"strconv"
)

type Feed struct {
	ID      string `xml:"id"`
	Total   int    `xml:"totalResults"`
	Page    int
	Start   int     `xml:"startIndex"`
	PerPage int     `xml:"itemsPerPage"`
	Videos  []Video `xml:"entry"`
	Pages   int
}

type Video struct {
	Published string    `xml:"published"`
	Updated   string    `xml:"updated"`
	Player    VideoURL  `xml:"content"`
	Title     string    `xml:"title"`
	Details   VideoInfo `xml:"group"`
}

type VideoURL struct {
	URL  string `xml:"src,attr"`
	Page string `xml:"url,attr"`
}

type VideoInfo struct {
	Description string       `xml:"description"`
	Title       string       `xml:"title"`
	VideoID     string       `xml:"videoid"`
	Images      []Thumbnails `xml:"thumbnail"`
	WatchPage   VideoURL     `xml:"player"`
}

type Thumbnails struct {
	URL    string `xml:"url,attr"`
	Height int    `xml:"height,attr"`
	Width  int    `xml:"width,attr"`
	Size   string `xml:"name,attr"`
}

func GetAll(page int, perpage int) (Feed, error) {
	var videos Feed
	start := ((page - 1) * perpage) + 1
	url := "https://gdata.youtube.com/feeds/api/videos?author=curtmfg&orderby=published&start-index=" + strconv.Itoa(start) + "&max-results=" + strconv.Itoa(perpage) + "&v=2"
	res, err := rest.Get(url)
	if err != nil {
		return videos, err
	}
	err = xml.Unmarshal(res, &videos)
	if err != nil {
		return videos, err
	}
	videos.Page = page
	videos.GetPageCount()
	return videos, nil
}

func Get(ytID string) (Video, error) {
	var video Video
	url := "https://gdata.youtube.com/feeds/api/videos/" + ytID + "?v=2"
	res, err := rest.Get(url)
	if err != nil {
		return video, err
	}
	err = xml.Unmarshal(res, &video)
	return video, err
}

func (v *Video) GetThumb() string {
	for _, img := range v.Details.Images {
		if img.Size == "default" {
			return img.URL
		}
	}
	return ""
}

func (v *Video) GetScreenshot() string {
	for _, img := range v.Details.Images {
		if img.Size == "hqdefault" {
			return img.URL
		}
	}
	return ""
}

func (f *Feed) GetPageCount() {
	f.Pages = int(math.Ceil(float64(f.Total) / float64(f.PerPage)))
}
