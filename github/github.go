package github

import (
	"regexp"
	"strings"
)

func getLastVersions(text string, count int, term string) string {
	if !strings.Contains(text, term) {
		return ""
	}

	taskIdReg := regexp.MustCompile(`(?m)\[` + term + `\-[0-9]+.+$`)
	versionReg := regexp.MustCompile(`# (v\d+\.\d+\.\d+)`)

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
