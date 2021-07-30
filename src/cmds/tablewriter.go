package cmds

import (
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

func NewTableWriter(w io.Writer, data interface{}, format string) {

	if reflect.ValueOf(data).Len() == 0 {
		fmt.Fprintln(w, "List is empty")
		return
	}

	// initialize tabwriter
	table := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	table.Init(w, 0, 0, 3, ' ', 0)
	defer table.Flush()

	var (
		dataStruct = reflect.Indirect(reflect.ValueOf(data).Index(0))
		tableRows  = dataStruct.NumField()
	)

	if format != "wide" {
		if dataStruct.NumField() > 6 {
			tableRows = 6
		}
	}

	columns := ""
	for i := 0; i < tableRows; i++ {
		field := dataStruct.Type().Field(i)
		if field.Tag.Get("json") != "-" {
			columns += fmt.Sprintf("%v\t", field.Name)
		}
	}
	fmt.Fprintln(table, columns)

	datarow := ""
	for i := 0; i < tableRows; i++ {
		field := dataStruct.Type().Field(i)
		if field.Tag.Get("json") != "-" {
			datarow += fmt.Sprintf("%v\t", dataStruct.Field(i).Interface())
		}
	}
	fmt.Fprintln(table, datarow)
}
