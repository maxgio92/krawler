package rpm

type Data struct {
	Type     string   `xml:"type,attr"`
	Location Location `xml:"location"`
}

type Location struct {
	Href string `xml:"href,attr"`
}

func (d *Data) GetLocation() string {
	return d.Location.Href
}
