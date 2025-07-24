package structs

import "encoding/xml"

type Generic struct {
	XMLName xml.Name `xml:"Event"`
	System  struct {
		EventID     int `xml:"EventID"`
		TimeCreated struct {
			SystemTime string `xml:"SystemTime,attr"`
		} `xml:"TimeCreated"`
	} `xml:"System"`
}

type Event1149 struct {
	XMLName xml.Name `xml:"Event"`
	System  struct {
		TimeCreated struct {
			SystemTime string `xml:"SystemTime,attr"`
		} `xml:"TimeCreated"`
	} `xml:"System"`
	UserData struct {
		EventXML struct {
			Param1 string `xml:"Param1"`
			Param3 string `xml:"Param3"`
		} `xml:"EventXML"`
	} `xml:"UserData"`
}

type Event21 struct {
	XMLName xml.Name `xml:"Event"`
	System  struct {
		TimeCreated struct {
			SystemTime string `xml:"SystemTime,attr"`
		} `xml:"TimeCreated"`
	} `xml:"System"`
	UserData struct {
		EventXML struct {
			User      string `xml:"User"`
			Address   string `xml:"Address"`
			SessionID int    `xml:"SessionID"`
		} `xml:"EventXML"`
	} `xml:"UserData"`
}
