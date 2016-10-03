package main

import (
	"fmt"
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

func (fmtr *Formater) Parse(wh *Webhook) {
	fm := template.FuncMap{
		"join":           strings.Join,
		"formatDatetime": formatDatetime,
	}
	tmpl, err := template.New("main").Funcs(fm).Parse(fmtr.Format)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
	if err := tmpl.Execute(os.Stdout, wh); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}

func formatDatetime(date string) string {
	return strings.Replace(date, " ", "T", -1)
}
