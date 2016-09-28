package xmlparser

import (
	"encoding/xml"
)

type VolumeList struct {
	XMLName xml.Name `xml:"volumeList"`
	General uint8    `xml:"general,omitempty"`
	Bv      uint8    `xml:"bv,omitempty"`
	lead1   uint8    `xml:"lead1,omitempty"`
	lead2   uint8    `xml:"lead2,omitempty"`
}
type Status struct {
	XMLName  xml.Name   `xml:"status"`
	State    string     `xml:"state,attr,omitempty"`
	Position float32    `xml:"position,omitempty"`
	Volumes  VolumeList `xml:"volumeList"`
	Pitch    uint8      `xml:"pitch,omitempty"`
	Tempo    uint8      `xml:"tempo,omitempty"`
	Queue    Queue      `xml:"queue"`
}

type Queue struct {
	XMLName xml.Name `xml:"queue"`
	Items   []Item   `xml:"item"`
}
type Item struct {
	XMLName  xml.Name `xml:"item"`
	State    string   `xml:"status,attr,omitempty"`
	Title    string   `xml:"title,omitempty"`
	Artist   string   `xml:"artist,omitempty"`
	Year     string   `xml:"year,omitempty"`
	Duration float32  `xml:"duration,omitempty"`
	Singer   string   `xml:"singer,omitempty"`
}


func ParseXML(data string) (*Status,error) {
	bdata := []byte(data)

	var state = Status{}
	err := xml.Unmarshal([]byte(bdata), &state)

	return &state,err
}