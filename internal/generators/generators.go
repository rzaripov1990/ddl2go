package generators

import (
	"ddl2go/internal/entity"
	"log"
	"os"
	"text/template"
)

func GoStruct(packageName string, tableName string, table entity.Table) {
	tmpl := `type {{ .Name }} struct {
{{- range .Columns }} 
	{{ .Name }} {{ .GoType }} {{ .Tag }} {{if .IsComment }}//{{ .Comment }} {{ end }}
{{- end }}
}

type {{ .Name }}Arr []{{ .Name }}`
	t, err := template.New("struct").Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(packageName + "/" + tableName + ".go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	uniqPackages := map[string]bool{}
	imports := ""
	for i := range table.Columns {
		if _, ok := uniqPackages[table.Columns[i].GoPackage]; !ok && table.Columns[i].GoPackage != "" {
			uniqPackages[table.Columns[i].GoPackage] = true
			if imports == "" {
				imports += "import ("
			}
			imports += "\n    \"" + table.Columns[i].GoPackage + "\""
		}
	}

	_, _ = file.WriteString("package " + packageName + "\n\n")
	if imports != "" {
		_, _ = file.WriteString(imports + "\n)\n\n")
	}

	err = t.Execute(file, table)
	if err != nil {
		log.Fatal(err)
	}
}

func GraphQLSchema(packageName string, tableName string, table entity.Table) {
	tmpl := `type {{ .Name }} {
{{- range .Columns }} 
	{{ .Name }}: {{ .GraphQLType }} {{if .IsComment }} # {{ .Comment }} {{ end }}
{{- end }}
}

extend type Query {
    get{{ .Name }}(id: ID!): {{ .Name }}
    list{{ .Name }}: [{{ .Name }}!]!
}`
	t, err := template.New("graphql").Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(packageName + "/" + tableName + "_schema.graphql")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = t.Execute(file, table)
	if err != nil {
		log.Fatal(err)
	}
}
