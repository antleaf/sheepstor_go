package internal

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-github/v71/github"
)

func GitHubWebHookHandler(resp http.ResponseWriter, req *http.Request) {
	Log.Debug("Handling GitHUb webhook post....")
	websitePtr := Registry.GetWebsiteByID(chi.URLParam(req, "website_id"))
	payload, err := github.ValidatePayload(req, []byte(os.Getenv(websitePtr.GitHubWebHookSecretEnvKey)))
	if err != nil {
		Log.Error(err.Error())
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		Log.Error(err.Error())
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	switch e := event.(type) {
	case *github.PushEvent:
		Log.Debug("Github push event received")
		if websitePtr != nil {
			website := *websitePtr
			Log.Debugf("SheepstorWebsite identified from GitHub push event; '%s'", website.ID)
			gitRepo := website.GitRepo
			localCommitID := gitRepo.GetHeadCommitID()
			pushCommitID := *e.HeadCommit.ID
			if localCommitID != pushCommitID {
				Log.Debugf("Attempting to build website '%s'", website.ID)
				err = website.Process()
				if err != nil {
					Log.Error(err.Error())
					http.Error(resp, err.Error(), http.StatusInternalServerError)
					return
				} else {
					Log.Infof("Built website '%s'", website.ID)
				}
			}
		} else {
			Log.Errorf("SheepstorWebsite with repo name '%s' and branch ref '%s' not found", e.GetRepo().GetFullName(), e.GetRef())
		}
	default:
		return
	}
}
