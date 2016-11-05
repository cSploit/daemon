package jobs

import (
	"strings"
	"unicode/utf8"
)

type (
	OutputHolder struct {
		Output []*outputLine `gorm:"-" json:"output"`
		//TODO: hide merged
		MergedOutput *string `json:"-"`
	}

	outputLine struct {
		dest outputDestination `json:"destination"`
		line string            `json:"line"`
	}

	outputDestination rune
)

const (
	stdoutDest outputDestination = 'O'
	stderrDest outputDestination = 'E'
)

func parseOutputLine(line string) *outputLine {

	first, i := utf8.DecodeRuneInString(line)

	return &outputLine{
		dest: outputDestination(first),
		line: line[i:],
	}
}

func parseText(text string) (res []*outputLine) {
	for _, line := range strings.Split(text, "\n") {
		res = append(res, parseOutputLine(line))
	}
	return
}

func (oh *OutputHolder) AddToStdout(line string) {
	if oh.MergedOutput == nil {
		t := ""
		oh.MergedOutput = &t
	}

	s := *oh.MergedOutput
	s += string(stdoutDest) + line + "\n"
}

func (oh *OutputHolder) AddToStderr(line string) {
	if oh.MergedOutput == nil {
		t := ""
		oh.MergedOutput = &t
	}

	s := *oh.MergedOutput
	s += string(stderrDest) + line + "\n"
}

func (oh *OutputHolder) Load() {
	if oh.MergedOutput == nil {
		return
	}

	oh.Output = parseText(*oh.MergedOutput)
}
