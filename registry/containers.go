package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/marcinhlybin/docker-env/docker"
	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
)

// Table options
const (
	minWidth int  = 1
	tabWidth int  = 1
	padding  int  = 5
	padChar  byte = ' '
	flags    uint = 0
)

func (registry *DockerProjectRegistry) ListContainers() error {
	containers, err := registry.fetchContainers()
	if err != nil {
		return err
	}
	if containers == nil {
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, minWidth, tabWidth, padding, padChar, flags)

	// Print header
	format := "%s\t%s\t%s\t%s\t%s\n"
	fmt.Fprintf(w, format,
		helpers.BoldText("Container name"),
		helpers.BoldText("Project"),
		helpers.BoldText("Service"),
		helpers.BoldText("State"),
		helpers.BoldText("Created at"),
	)

	// Print columns
	for _, c := range containers {
		fmt.Fprintf(w, format,
			helpers.NormalText(c.Name),
			helpers.NormalText(c.ProjectName()),
			helpers.NormalText(c.ServiceName()),
			colorState(c.State),
			helpers.NormalText(formatTime(c.CreatedAt)),
		)
	}

	w.Flush()

	return nil
}

func (registry *DockerProjectRegistry) fetchContainers() ([]docker.Container, error) {
	logger.Debug("Fetching container names")
	dc := registry.dockerCmd.FetchContainersCommand()
	jsonRecords, err := dc.ExecuteWithOutput()
	if err != nil {
		return nil, err
	}

	return registry.createContainersFromJson(jsonRecords)
}

func (registry *DockerProjectRegistry) createContainersFromJson(jsonRecords []string) ([]docker.Container, error) {
	var containers []docker.Container

	for _, jsonRecord := range jsonRecords {
		var c docker.Container

		err := json.Unmarshal([]byte(jsonRecord), &c)
		if err != nil {
			return nil, err
		}

		containers = append(containers, c)
	}

	return containers, nil
}

func colorState(state string) string {
	switch state {
	case "running":
		return helpers.GreenText(state)
	default:
		return helpers.NormalText(state)
	}
}

func formatTime(timeStr string) string {
	const layout = "2006-01-02 15:04:05 -0700 MST"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return timeStr // Return the original string if parsing fails
	}
	return parsedTime.Format("2006-01-02 15:04")
}
