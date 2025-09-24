package utils

import (
	"html/template"

	"github.com/gogf/gf/v2/os/gtime"
)

func GetOutputPathTemplate() *template.Template {
	return template.
		Must(template.
			New("outputPathTemplate").
			Funcs(getFuncsMap()).
			Parse(`{{ outputPath }}/{{ .Platform | sanitize }}/{{ .Anchor | sanitize }}/{{ currentMonth }}/`))
}

func GetFilenameTemplate(outputPath, format string) *template.Template {
	return template.
		Must(template.
			New("filenameTemplate").
			Funcs(getFuncsMap()).
			Parse(outputPath + `[{{ currentTime }}][{{ .Anchor | sanitize }}][{{ .RoomName | sanitize }}].` + format))
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
			if end > len(runes) {
				end = len(runes)
			}
			return string(runes[start:end])
		},
		"outputPath":       GetOutputPath,
		"downloadPath":     GetDownloadPath,
		"tempDownloadPath": GetTempDownloadPath,
		"sanitize":         SanitizeFilename,
	}
}
