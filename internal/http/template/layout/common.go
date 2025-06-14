package layout

import "fmt" 

type Props struct {
	Title string
}

func TitleTemplate(val string) string {
	if val == "" {
		return "Personal"
	}

	return fmt.Sprintf("%s | Personal", val)
}

func Fallback(val string, fallback string) string {
	if val == "" {
		return fallback
	}

	return val
}

