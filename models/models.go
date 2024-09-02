package models

type Blog struct {
	Id int `json:"id"`
	Author string `json:"author"`
	Title string `json:"title"`
	Text string `json:"text"`
	CreationTime string `json:"creation_time"`
}
