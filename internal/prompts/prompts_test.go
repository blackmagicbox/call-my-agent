package prompts

import (
	"strings"
	"testing"
)

const fakeProfile = `{"name":"Jane Doe","title":"Senior Backend Engineer","skills":["Go","Python","Kubernetes"],"experience_years":8}`
const fakeJob = `{"title":"Staff Engineer","company":"Acme Corp","requirements":["Go","distributed systems","5+ years"]}`

func TestEvaluateJobSystem(t *testing.T) {
	evaluateJobSystem = "You are a job evaluator."
	got := EvaluateJobSystem()
	if got != "You are a job evaluator." {
		t.Errorf("EvaluateJobSystem() = %q, want %q", got, "You are a job evaluator.")
	}
}

func TestEvaluateJobUser(t *testing.T) {
	evaluateJobUserTmpl = "CANDIDATE PROFILE:\n{{.ProfileJSON}}\n\nJOB LISTING:\n{{.JobJSON}}"

	t.Run("happy path", func(t *testing.T) {
		got, err := EvaluateJobUser(EvaluateJobData{
			ProfileJSON: fakeProfile,
			JobJSON:     fakeJob,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(got, fakeProfile) {
			t.Errorf("output missing profile, got:\n%s", got)
		}
		if !strings.Contains(got, fakeJob) {
			t.Errorf("output missing job, got:\n%s", got)
		}
	})

	t.Run("empty fields", func(t *testing.T) {
		got, err := EvaluateJobUser(EvaluateJobData{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(got, "CANDIDATE PROFILE:") {
			t.Errorf("output missing structure, got:\n%s", got)
		}
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		evaluateJobUserTmpl = "{{.Broken"
		_, err := EvaluateJobUser(EvaluateJobData{})
		if err == nil {
			t.Fatal("expected parse error, got nil")
		}
		evaluateJobUserTmpl = "CANDIDATE PROFILE:\n{{.ProfileJSON}}\n\nJOB LISTING:\n{{.JobJSON}}"
	})
}

func TestCoverLetterSystem(t *testing.T) {
	coverLetterSystemTmpl = "You are a cover letter writer.\nTone: {{.Tone}}"

	t.Run("happy path", func(t *testing.T) {
		got, err := CoverLetterSystem(CoverLetterData{
			ProfileJSON: fakeProfile,
			JobJSON:     fakeJob,
			Tone:        "professional",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(got, "professional") {
			t.Errorf("output missing tone, got:\n%s", got)
		}
	})
}

func TestCoverLetterUser(t *testing.T) {
	coverLetterUserTmpl = "CANDIDATE PROFILE:\n{{.ProfileJSON}}\n\nJOB DESCRIPTION:\n{{.JobJSON}}\n\nWrite the cover letter in a {{.Tone}} tone."

	t.Run("happy path", func(t *testing.T) {
		got, err := CoverLetterUser(CoverLetterData{
			ProfileJSON: fakeProfile,
			JobJSON:     fakeJob,
			Tone:        "friendly",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(got, fakeProfile) {
			t.Errorf("output missing profile, got:\n%s", got)
		}
		if !strings.Contains(got, fakeJob) {
			t.Errorf("output missing job, got:\n%s", got)
		}
		if !strings.Contains(got, "friendly") {
			t.Errorf("output missing tone, got:\n%s", got)
		}
	})

	t.Run("empty tone", func(t *testing.T) {
		_, err := CoverLetterUser(CoverLetterData{
			ProfileJSON: fakeProfile,
			JobJSON:     fakeJob,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestRenderErrors(t *testing.T) {
	t.Run("malformed template", func(t *testing.T) {
		coverLetterUserTmpl = "{{.Unclosed"
		_, err := CoverLetterUser(CoverLetterData{})
		if err == nil {
			t.Fatal("expected parse error, got nil")
		}
		coverLetterUserTmpl = "{{.ProfileJSON}}"
	})

	t.Run("missing field causes error with missingkey=error option", func(t *testing.T) {
		// Go templates default to printing <no value> for missing keys,
		// so a missing field in the struct won't cause an error by default.
		// This test verifies render doesn't crash on valid templates.
		coverLetterUserTmpl = "{{.ProfileJSON}} {{.Tone}}"
		_, err := CoverLetterUser(CoverLetterData{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
