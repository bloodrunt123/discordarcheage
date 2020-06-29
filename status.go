package archeage

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const serverStatusURL = "https://archeage.xlgames.com/serverstatus"

const (
	serverStatusRowQuery = `table tr`
)

type ServerStatus map[string]bool

func (old ServerStatus) DiffString(new ServerStatus) (diff ServerStatus, isDiff bool) {
	for k := range new {
		if old[k] != new[k] {
			diff[k] = (old[k] == false && new[k] == true)
			isDiff = true
		}
	}
	return
}

func (a *ArcheAge) FetchServerStatus() (serverStatus ServerStatus, err error) {
	doc, err := a.get(serverStatusURL)
	if err != nil {
		return
	}
	serverStatus = ServerStatus{}
	doc.Find(serverStatusRowQuery).Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".server").Text())
		if name == "" {
			return
		}
		status := s.Find(".stats span").HasClass("on")
		serverStatus[name] = status
	})
	return
}
