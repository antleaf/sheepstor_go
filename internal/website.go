package internal

import (
	"os"
	"path/filepath"
	"sync"
)

type Website struct {
	ID                        string  `yaml:"id"`
	ContentProcessor          string  `yaml:"content_processor"` // either 'hugo' or nil
	ProcessorRoot             string  `yaml:"processor_root"`    // e.g. a sub-folder in the repo called 'webroot'
	WebRoot                   string  `yaml:"-"`
	IndexForSearch            bool    `yaml:"index"` // run the pagefind executable to create a search index
	GitHubWebHookSecretEnvKey string  `yaml:"github_webhook_secret_env_key"`
	GitRepo                   GitRepo `yaml:"git"`
}

func (w *Website) Initialise(sourcesRoot, docsRoot string) error {
	w.WebRoot = filepath.Join(docsRoot, w.ID)
	w.ProcessorRoot = filepath.Join(sourcesRoot, w.ID, w.ProcessorRoot)
	return os.MkdirAll(w.WebRoot, os.ModePerm)
}

func (w *Website) Process() error {
	err := w.provisionSources()
	if err != nil {
		Log.Error(err.Error())
		return err
	}
	err = w.build()
	if err != nil {
		Log.Error(err.Error())
	} else {
		Log.Infof("Built website: '%s'", w.ID)
	}
	return err
}

func (w *Website) ProcessInSynchronousWorker(wg *sync.WaitGroup) {
	_ = w.Process()
	wg.Done()
}

func (w *Website) build() error {
	var err error
	targetFolderPathForBuild := filepath.Join(w.WebRoot, "public_1")
	symbolicLinkPath := filepath.Join(w.WebRoot, "public")
	currentSymLinkTargetPath, readlinkErr := os.Readlink(symbolicLinkPath)
	if readlinkErr == nil {
		if currentSymLinkTargetPath == filepath.Join(w.WebRoot, "public_1") {
			targetFolderPathForBuild = filepath.Join(w.WebRoot, "public_2")
		}
	}
	if _, statErr := os.Stat(targetFolderPathForBuild); statErr == nil {
		os.RemoveAll(targetFolderPathForBuild)
	}
	err = os.MkdirAll(targetFolderPathForBuild, os.ModePerm)
	err = os.MkdirAll(filepath.Join(w.WebRoot, "logs"), os.ModePerm)
	if err != nil {
		return err
	}
	switch w.ContentProcessor {
	case "Hugo":
		err = HugoProcessor(w.ProcessorRoot, targetFolderPathForBuild)
		if err != nil {
			Log.Error("Hugo Processor error")
			return err
		}
	default:
		err = DefaultProcessor(w.ProcessorRoot, targetFolderPathForBuild)
		if err != nil {
			Log.Error("Default Processor error")
			return err
		}
	}
	if w.IndexForSearch {
		err = IndexForSearch(targetFolderPathForBuild)
		if err != nil {
			return err
		}
	}

	if _, err = os.Lstat(symbolicLinkPath); err == nil {
		if err = os.Remove(symbolicLinkPath); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		// do nothing?
	}
	err = os.Symlink(targetFolderPathForBuild, symbolicLinkPath) // Only switch if successful
	if err != nil {
		return err
	}
	return err
}

func (w *Website) provisionSources() error {
	var err error
	gitFolderPath := filepath.Join(w.GitRepo.WorkingDir, ".git")
	if _, err = os.Stat(gitFolderPath); os.IsNotExist(err) {
		err = os.MkdirAll(w.GitRepo.WorkingDir, os.ModePerm)
		if err != nil {
			return err
		}
		err = w.GitRepo.Clone()
		if err != nil {
			return err
		}
	} else {
		err = w.GitRepo.Pull()
		if err != nil {
			return err
		}
	}
	return err
}

func (w *Website) commitAndPush(message string) error {
	err := w.GitRepo.Pull()
	if err != nil {
		return err
	}
	err = w.GitRepo.CommitAndPush(message)
	if err != nil {
		return err
	}
	return err
}
