package main

import (
	"html/template"
	"os"
	"strings"
)

type Formater struct {
	Format string
	Vars   interface{}
}

func NewFormater(format string) (*Formater, error) {
	fmtr := &Formater{Format: format}
	return fmtr, nil
}

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
