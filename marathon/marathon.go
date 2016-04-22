package marathon

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Response struct {
	Apps  []App  `json:"apps"`
	Tasks []Task `json:"tasks"`
}

type App struct {
	ID     string            `json:"id"`
	Labels map[string]string `json:"labels"`
}

type Task struct {
	ID     string `json:"id"`
	AppID  string `json:"appId"`
	Host   string `json:"host"`
	Ports  []int  `json:"ports"`
	Labels map[string]string
}

func FetchApps(nodes []string) (*Response, error) {
	labels, err := FetchAppsLabels(nodes)
	if err != nil {
		logrus.Error("could not retrieve app labels")
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/v2/tasks?status=running", nodes[0]))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var apps Response
	json.NewDecoder(resp.Body).Decode(&apps)

	for idx, task := range apps.Tasks {
		apps.Tasks[idx].Labels = labels[task.AppID]
	}

	return &apps, nil
}

func FetchAppsLabels(nodes []string) (map[string]map[string]string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/v2/apps", nodes[0]))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var apps Response
	json.NewDecoder(resp.Body).Decode(&apps)

	labels := make(map[string]map[string]string)
	for _, app := range apps.Apps {
		labels[app.ID] = app.Labels
	}

	return labels, nil
}
