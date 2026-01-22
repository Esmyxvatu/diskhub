package render

import (
	"fmt"
	"html/template"
)

type DumpRow struct {
	Offset string
	Bytes  []string
	ASCII  string
}

func DumpToHTML(data []byte) template.HTML {
	rows := make([]DumpRow, 0)
	for i := 0; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}

		byteStr := make([]string, 0, end-i)
		for _, b := range data[i:end] {
			byteStr = append(byteStr, fmt.Sprintf("%02x", b))
		}

		ascii := ""
		for _, b := range data[i:end] {
			if b >= 32 && b <= 126 {
				ascii += string(b)
			} else {
				ascii += "•"
			}
		}

		rows = append(rows, DumpRow{
			Offset: fmt.Sprintf("%08x", i),
			Bytes:  byteStr,
			ASCII:  ascii,
		})
	}

	tmpl := "<table id='dump'><tr><th>Offset</th><th>Bytes</th><th>ASCII</th></tr>"

	for _, row := range rows {
		tmpl += "<tr><td>" + row.Offset + "</td><td>"
		for _, bytes := range row.Bytes {
			tmpl += bytes + " "
		}
		tmpl += "</td><td>" + row.ASCII + "</td>"
	}

	tmpl += "</table>"

	return template.HTML(tmpl)
}
