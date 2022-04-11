package main

import (
	"fmt"
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
	Attributes []*Attribute
}

type Attribute struct {
	Name     string
	Type     string
	Primary  bool
	Key      bool
	Nullable bool
	IsIndex  bool
}

type ForeignKey struct {
	TargetTable  string
	TargetColumn string
	SourceTable  string
	SourceColumn string
}

func RenderFromTemplate(diagram *Diagram, tplFile string) error {
	template, err := template.ParseFiles(tplFile)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	err = template.Execute(os.Stdout, diagram)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func addIndex(entities []*Entity, table, key string, unique bool) {
	typ := "INDEX"
	if unique {
		typ = "UNIQUE INDEX"
	}

	for _, entity := range entities {
		if entity.Name == table {
			entity.Attributes = append(entity.Attributes, &Attribute{
				Name:     key,
				Type:     typ,
				Nullable: false,
				IsIndex:  true,
			})

			return
		}
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
