package prompts

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed evaluate_job.system.txt
var evaluateJobSystem string

//go:embed evaluate_job.user.tmpl
var evaluateJobUserTmpl string

//go:embed cover_letter.system.tmpl
var coverLetterSystemTmpl string

//go:embed cover_letter.user.tmpl
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
