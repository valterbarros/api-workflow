package github

import (
	client "api-workflow/internal/client/github"
	"fmt"
	"regexp"
	"strings"
)

var (
	versionReg = regexp.MustCompile(`# (v\d+\.\d+\.\d+)`)
)

func GetLastVersions(text string, count int, term string) (allVersion string) {
	if !strings.Contains(text, term) {
		return ""
	}

	taskIdReg := regexp.MustCompile(`(?m)\[` + term + `\-[0-9]+.+$`)

	lines := strings.Split(text, "\n")
	var out strings.Builder

	for _, line := range lines {
		if count < 0 {
			return out.String()
		}

		version := versionReg.FindString(line)

		if version != "" {
			out.WriteString(version + "\n")
			count--
		}

		taskTitle := taskIdReg.FindString(line)

		if taskTitle != "" {
			out.WriteString("- " + taskTitle + "\n")
		}
	}

	return out.String()
}

func GetVersion(text, taskId string) (lastVersion string) {
	if !strings.Contains(text, taskId) {
		return ""
	}

	lines := strings.Split(text, "\n")

	return findVersion(lines, 0, taskId, "")
}

func findVersion(lines []string, index int, taskId, lastVersion string) string {
	line := lines[index]

	version := versionReg.FindString(line)

	if strings.Contains(line, taskId) {
		return lastVersion + "\n" + line
	}

	if version != "" {
		return findVersion(lines, index+1, taskId, version)
	} else {
		return findVersion(lines, index+1, taskId, lastVersion)
	}
}

func GetPullRequests() client.PullRequests {
	return client.GetPullRequests("")
}

type ReviewData struct {
	Text string `json:"text"`
	// green or yellow
	Color string `json:"color"`
}

const MYUSER = "valter-clk"

func ReviewsResult(url string) ReviewData {
	reviews := client.GetReviews(url)

	var approveCount int
	var meType string
	mapper := map[string]int{}

	for _, review := range reviews {
		if review.User.Login == MYUSER {
			meType = mapLabel(review)
		}

		switch review.State {
		case "APPROVED":
			mapper["AP"]++
			approveCount++
		case "COMMENTED":
			mapper["CO"]++
		case "DISMISSED":
		}
	}

	var text string
	// Map iteration order is random in Go; use a fixed key order for stable text/tests.
	for _, k := range []string{"CO", "AP"} {
		if v := mapper[k]; v > 0 {
			text += fmt.Sprintf("%s: %d ", k, v)
		}
	}

	if len(mapper) == 0 {
		return ReviewData{"No reviews", "yellow"}
	}

	if meType != "" {
		text += "ME: " + meType
	}

	var color string
	if approveCount > 2 {
		color = "green"
	} else {
		color = "yellow"
	}

	return ReviewData{strings.TrimSpace(text), color}
}

func mapLabel(review client.User) string {
	switch review.State {
	case "APPROVED":
		return "OK"
	case "COMMENTED":
		return "CO"
	case "DISMISSED":
		return "DI"
	default:
		return ""
	}
}
