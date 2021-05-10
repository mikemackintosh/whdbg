package main

import (
	"html/template"
	"os"
	"strings"
)

// Formater is a format required string
type Formater struct {
	Format string
	Vars   interface{}
}

// NewFormater will create a new format object using the provided string
func NewFormater(format string) (*Formater, error) {
	fmtr := &Formater{Format: format}
	return fmtr, nil
}

// Parse will parse the format string using the provided Webhook
func (fmtr *Formater) Parse(wh *Webhook) error {
	fm := template.FuncMap{
		"join":           strings.Join,
		"formatDatetime": formatDatetime,
	}
	tmpl, err := template.New("main").Funcs(fm).Parse(fmtr.Format)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(os.Stdout, wh); err != nil {
		return err
	}

	return nil
}

func formatDatetime(date string) string {
	return strings.Replace(date, " ", "T", -1)
}
