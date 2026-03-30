package github

import (
	client "api-workflow/internal/client/github"
	"math"
	"regexp"
	"slices"
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

	slices.Reverse(reviews)
	var me string
	var notMe = []string{}
	var approveCount int

	for _, review := range reviews {
		// MY APPROVED
		if review.User.Login == MYUSER {
			me = "ME: " + mapLabel(review)
			approveCount++
			continue
		}

		firstL := strings.ToUpper(string(review.User.Login[0]))

		if review.User.Login != MYUSER {
			switch review.State {
			case "APPROVED":
				approveCount++
				notMe = append(notMe, firstL+": "+mapLabel(review))
			case "COMMENTED":
				notMe = append(notMe, firstL+": "+mapLabel(review))
			case "DISMISSED":
				notMe = append(notMe, firstL+": "+mapLabel(review))
			}
		}
	}

	var result = []string{}

	if len(notMe) > 1 {
		result = notMe[0:2]
	}

	if me != "" {
		result = append(result, me)
	}

	if len(result) < 1 {
		return ReviewData{"No reviews", "yellow"}
	}

	size := float64(len(result))
	start := int(math.Max(size-2, 0))
	result = result[start:]

	var color string
	if approveCount > 2 {
		color = "green"
	} else {
		color = "yellow"
	}

	return ReviewData{strings.Join(result, " "), color}
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
