package prompts

import (
	"bytes"
	"text/template"
)

var evaluateJobSystem string
var evaluateJobUserTmpl string
var coverLetterSystemTmpl string
var coverLetterUserTmpl string

func EvaluateJobSystem() string {
	return evaluateJobSystem
}

type EvaluateJobData struct {
	ProfileJSON string
	JobJSON     string
}

type CoverLetterData struct {
	ProfileJSON string
	JobJSON     string
	Tone        string
}

func EvaluateJobUser(data EvaluateJobData) (string, error) {
	return render(evaluateJobUserTmpl, data)
}

func CoverLetterSystem(data CoverLetterData) (string, error) {
	return render(coverLetterSystemTmpl, data)
}

func CoverLetterUser(data CoverLetterData) (string, error) {
	return render(coverLetterUserTmpl, data)
}

func render(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
