package utils

import (
	"html/template"

	"github.com/gogf/gf/v2/os/gtime"
)

var (
	filenameTemplateMap = map[int]string{
		0: `[{{ currentTime }}][{{ .Anchor | sanitize }}][{{ .RoomName | sanitize }}].`,
		1: `[{{ .Anchor | sanitize }}][{{ .RoomName | sanitize }}][{{ currentTime }}].`,
		2: `{{ currentTime }}_{{ .Anchor | sanitize }}_{{ .RoomName | sanitize }}.`,
		3: `{{ .Anchor | sanitize }}_{{ .RoomName | sanitize }}_{{ currentTime }}.`,
		4: `[{{ currentTime }}][{{ .Anchor | sanitize }}].`,
		5: `[{{ .Anchor | sanitize }}][{{ currentTime }}].`,
		6: `{{ currentTime }}_{{ .Anchor | sanitize }}.`,
		7: `{{ .Anchor | sanitize }}_{{ currentTime }}.`,
	}
	outputPathTemplateMap = map[int]string{
		0: `{{ outputPath }}/{{ .Platform | sanitize }}/{{ .Anchor | sanitize }}/{{ currentMonth }}/`,
		1: `{{ outputPath }}/{{ .Platform | sanitize }}/{{ .Anchor | sanitize }}/{{ currentDate }}/`,
		2: `{{ outputPath }}/{{ .Platform | sanitize }}/{{ .Anchor | sanitize }}/{{ currentMonth }}/{{ currentDate }}/`,
		3: `{{ outputPath }}/{{ .Platform | sanitize }}/{{ .Anchor | sanitize }}/`,
	}
)

func GetOutputPathTemplate(index int) *template.Template {
	var t string
	t = outputPathTemplateMap[0]
	if _, ok := outputPathTemplateMap[index]; ok {
		t = outputPathTemplateMap[index]
	}
	return template.
		Must(template.
			New("outputPathTemplate").
			Funcs(getFuncsMap()).
			Parse(t))
}

func GetFilenameTemplate(outputPath, format string, index int) *template.Template {
	var t string
	t = filenameTemplateMap[0]
	if _, ok := filenameTemplateMap[index]; ok {
		t = filenameTemplateMap[index]
	}
	return template.
		Must(template.
			New("filenameTemplate").
			Funcs(getFuncsMap()).
			Parse(outputPath + t + format))
}

func GetDownloadPathTemplate(isTemp bool) *template.Template {
	templateStr := `{{ downloadPath }}/{{ .Platform }}/{{ currentMonth }}/`
	if isTemp {
		templateStr = `{{ tempDownloadPath }}/{{ .Platform }}/{{ currentMonth }}/`
	}
	return template.
		Must(template.
			New("downloadPathTemplate").
			Funcs(getFuncsMap()).
			Parse(templateStr))
}

func GetDownloadFilenameTemplate(outputPath, format string, random string) *template.Template {
	var templateStr string
	if random == "" {
		templateStr = outputPath + `[{{ if gt (runeCount .Title) 20 }}{{ runeSubString .Title 0 20 }}...{{ else }}{{ .Title }}{{ end }}].` + format
	} else {
		templateStr = outputPath + `[{{ if gt (runeCount .Title) 20 }}{{ runeSubString .Title 0 20 }}...{{ else }}{{ .Title }}{{ end }}]-` + random + "." + format
	}
	return template.
		Must(template.
			New("filenameTemplate").
			Funcs(getFuncsMap()).
			Parse(templateStr))
}

func getFuncsMap() template.FuncMap {
	return template.FuncMap{
		"currentTime": func() string {
			return gtime.Datetime()
		},
		"currentDate": func() string {
			return gtime.Date()
		},
		"currentMonth": func() string {
			return gtime.Now().Format("Y-m")
		},
		"runeCount": func(s string) int {
			return len([]rune(s))
		},
		"runeSubString": func(s string, start, length int) string {
			runes := []rune(s)
			if start >= len(runes) {
				return ""
			}
			end := start + length
			end = min(end, len(runes))
			return string(runes[start:end])
		},
		"outputPath":       GetOutputPath,
		"downloadPath":     GetDownloadPath,
		"tempDownloadPath": GetTempDownloadPath,
		"sanitize":         SanitizeFilename,
	}
}
