package env

import "bytes"

type CI struct {
	Name     string
	BuildID  string
	BuildURL string
}

func (c CI) String() string {
	out := &bytes.Buffer{}
	out.WriteString("CI_NAME=")
	out.WriteString(c.Name)
	out.WriteString("\nCI_BUILD_ID=")
	out.WriteString(c.BuildID)
	out.WriteString("\nCI_BUILD_URL=")
	out.WriteString(c.BuildURL)
	return out.String()
}

func loadCIInfo() CI {
	return CI{
		Name:     findVar(ciNameVars),
		BuildID:  findVar(ciBuildIDVars),
		BuildURL: findVar(ciBuildURLVars),
	}
}

var ciNameVars = []string{"CI_NAME", "APPVEYOR", "BUILDKITE", "CIRCLECI", "GITLAB_CI", "JENKINS_URL", "SEMAPHORE", "TDDIUM", "TRAVIS", "WERCKER"}

var ciBuildIDVars = []string{"CI_BUILD_ID", "APPVEYOR_BUILD_ID", "BUILDKITE_JOB_ID", "BUILD_NUMBER", "CIRCLE_BUILD_NUM", "CI_BUILD_NUMBER", "SEMAPHORE_BUILD_NUMBER", "TDDIUM_SESSION_ID", "TRAVIS_JOB_ID", "WERCKER_BUILD_ID"}

var ciBuildURLVars = []string{"CI_BUILD_URL", "APPVEYOR_API_URL", "BUILDKITE_BUILD_URL", "BUILD_URL", "CIRCLE_BUILD_NUM", "WERCKER_BUILD_URL"}
