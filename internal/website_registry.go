package internal

import (
	"path/filepath"
	"sync"
)

type WebsiteRegistry struct {
	SourceRoot string     `yaml:"source_root"`
	DocsRoot   string     `yaml:"docs_root"`
	WebSites   []*Website `yaml:"websites"`
}

var Registry *WebsiteRegistry

func (wr *WebsiteRegistry) Initialise() error {
	for _, website := range wr.WebSites {
		err := website.Initialise(wr.SourceRoot, wr.DocsRoot)
		if err != nil {
			return err
		}
		err = website.GitRepo.Initialise(filepath.Join(wr.SourceRoot, website.ID))
		if err != nil {
			return err
		}
	}
	return nil
}

func (wr *WebsiteRegistry) Add(w *Website) {
	wr.WebSites = append(wr.WebSites, w)
}

func (wr *WebsiteRegistry) GetWebsiteByID(id string) *Website {
	for _, w := range wr.WebSites {
		if w.ID == id {
			return w
		}
	}
	return nil
}

func (wr *WebsiteRegistry) ProcessAllWebsites() {
	var wg sync.WaitGroup
	for _, website := range wr.WebSites {
		wg.Add(1)
		go website.ProcessInSynchronousWorker(&wg)
	}
	wg.Wait()
}
