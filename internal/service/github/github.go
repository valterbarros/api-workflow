package github

import (
	client "api-workflow/internal/client/github"
	"api-workflow/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	text string
	// green or yellow
	color string
}

func ReviewsResult() ReviewData {
	const MYUSER = "valter-clk"
	// maybe temp here
	ser := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 16 comments
		// 1 approved
		// 0 dismissed
		data := utils.GetFileData(utils.GetFromRoot("fixtures/reviews.json"))
		w.Write([]byte(data))
	}))

	reviews := client.GetReviews(ser.URL)

	slices.Reverse(reviews)
	var me string
	data := []string{}
	var approveCount int

	for _, review := range reviews {
		// MY APPROVED
		if review.User.Login == MYUSER && review.State == "APPROVED" {
			me = "ME: OK"
			approveCount++
		}

		firstL := strings.ToUpper(string(review.User.Login[0]))

		if review.User.Login != MYUSER {
			switch review.State {
			case "APPROVED":
				approveCount++
				data = append(data, firstL+": OK")
			case "COMMENTED":
				data = append(data, firstL+": CO")
			case "DISMISSED":
				data = append(data, firstL+": DI")
			default:
				fmt.Println("default case")
			}
		}
	}

	var last []string

	if me != "" {
		last = []string{data[0]}
		last = append(last, me)
	} else {
		last = data[0:2]
	}

	fmt.Println(last)

	var color string
	if approveCount > 2 {
		color = "green"
	} else {
		color = "yellow"
	}

	return ReviewData{strings.Join(last, " "), color}
}
