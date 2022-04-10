package main

import (
	"log"
	"os"
	"text/template"
)

type Diagram struct {
	Name      string
	Date      string
	Entities  []*Entity
	Relations []*ForeignKey
}

type Entity struct {
	Name       string
	Attributes []*Column
}

type Column struct {
	Name    string
	Primary bool
	Key     bool
	Type    string
	NotNull bool
}

type ForeignKey struct {
	TargetTable  string
	TargetColumn string
	SourceTable  string
	SourceColumn string
}

func RenderFromTemplate(diagram *Diagram, tplFile string) {
	template, err := template.ParseFiles(tplFile)
	if err != nil {
		log.Fatalf("Error loading template: %s\n", err)
	}
	err = template.Execute(os.Stdout, diagram)
	if err != nil {
		log.Fatalf("Error generating digraph: %s\n", err)
	}
}

func setKey(entities []*Entity, table string, col string) {
	for _, entity := range entities {
		if entity.Name == table {
			for i := range entity.Attributes {
				if entity.Attributes[i].Name == col {
					entity.Attributes[i].Key = true
					return
				}
			}
		}
	}
}
