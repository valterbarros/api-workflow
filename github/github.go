package github

import (
	"regexp"
	"strings"
)

var (
	versionReg = regexp.MustCompile(`# (v\d+\.\d+\.\d+)`)
)

func getLastVersions(text string, count int, term string) (allVersion string) {
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

func getVersion(text, taskId string) (lastVersion string) {
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
		return lastVersion
	}

	if version != "" {
		return findVersion(lines, index+1, taskId, version)
	} else {
		return findVersion(lines, index+1, taskId, lastVersion)
	}
}
