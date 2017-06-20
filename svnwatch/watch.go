package svnwatch

import "encoding/xml"

type Watches struct {
	XMLName xml.Name `xml:"watches"`
	Watches []Watch  `xml:"watch"`
}

type Watch struct {
	XMLName    xml.Name  `xml:"watch"`
	Repository string    `xml:"repository"`
	Commands   []Command `xml:"command"`
}

func (w *Watch) Update(repo *Repository) (bool, error) {
	updated, err := repo.Update()

	if err != nil {
		return false, err
	}

	if updated {
		for _, cmd := range w.Commands {
			if cmd.Type == "revision" {
				cmd.Execute(repo)
			}
		}
	}

	return updated, nil
}
