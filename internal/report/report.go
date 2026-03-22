package report

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/creativoma/freelow/internal/timer"
)

type TaskSummary struct {
	Name     string
	Duration time.Duration
	Commits  int
	Messages []string
}

type Report struct {
	ClientName string
	Since      time.Time
	Until      time.Time
	Tasks      []TaskSummary
	Total      time.Duration
}

const reportTmpl = `## {{.PeriodLabel}} report · {{.ClientName}} · {{.DateRange}}

**Total hours:** {{.TotalStr}}

| Task | Time | Commits |
|---|---|---|
{{- range .Tasks}}
| {{.Name}} | {{.DurationStr}} | {{.CommitCount}} |
{{- end}}

### Changes
{{- range .Tasks}}
{{- range .Messages}}
- {{.}}
{{- end}}
{{- end}}`

type templateData struct {
	PeriodLabel string
	ClientName  string
	DateRange   string
	TotalStr    string
	Tasks       []templateTask
}

type templateTask struct {
	Name        string
	DurationStr string
	CommitCount int
	Messages    []string
}

func Generate(r *Report, periodLabel string) (string, error) {
	tasks := make([]templateTask, len(r.Tasks))
	for i, t := range r.Tasks {
		tasks[i] = templateTask{
			Name:        t.Name,
			DurationStr: timer.FormatDuration(t.Duration),
			CommitCount: t.Commits,
			Messages:    t.Messages,
		}
	}
	data := templateData{
		PeriodLabel: periodLabel,
		ClientName:  r.ClientName,
		DateRange:   fmt.Sprintf("%s–%s", r.Since.Format("2 Jan"), r.Until.Format("2 Jan 2006")),
		TotalStr:    timer.FormatDuration(r.Total),
		Tasks:       tasks,
	}
	tmpl, err := template.New("report").Parse(reportTmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func BuildFromSessions(sessions []timer.Session, clientName string, since, until time.Time) *Report {
	taskMap := make(map[string]*TaskSummary)
	var taskOrder []string

	for _, s := range sessions {
		if s.End == nil {
			continue
		}
		if s.End.Before(since) || s.Start.After(until) {
			continue
		}
		name := s.Task
		if _, ok := taskMap[name]; !ok {
			taskMap[name] = &TaskSummary{Name: name}
			taskOrder = append(taskOrder, name)
		}
		taskMap[name].Duration += s.ElapsedDuration()
		taskMap[name].Commits += len(s.Commits)
	}

	var tasks []TaskSummary
	var total time.Duration
	for _, name := range taskOrder {
		t := taskMap[name]
		tasks = append(tasks, *t)
		total += t.Duration
	}

	return &Report{
		ClientName: clientName,
		Since:      since,
		Until:      until,
		Tasks:      tasks,
		Total:      total,
	}
}
