package render

import (
	"strings"
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
		isLastLine := false
		if end > len(data) {
			isLastLine = true
			end = len(data)
		}

		byteStr := make([]string, 0, end-i)
		for _, b := range data[i:end] {
			byteStr = append(byteStr, fmt.Sprintf("%02x", b))
		}

		var ascii strings.Builder
		for _, b := range data[i:end] {
			if b >= 32 && b <= 126 {
				ascii.WriteString(string(b))
			} else {
				ascii.WriteString("•")
			}
		}

		if isLastLine {
			missing := 16 - (end - i)
			for range missing {
				byteStr = append(byteStr, "&nbsp;&nbsp;")
				ascii.WriteString("&nbsp;")
			}
		}

		rows = append(rows, DumpRow{
			Offset: fmt.Sprintf("%08x", i),
			Bytes:  byteStr,
			ASCII:  ascii.String(),
		})
	}

	var tmpl strings.Builder; tmpl.WriteString("<table id='dump'><tr><th>Offset</th><th>Bytes</th><th>ASCII</th></tr>")

	for _, row := range rows {
		tmpl.WriteString("<tr><td>" + row.Offset + "</td><td>")
		for _, bytes := range row.Bytes {
			tmpl.WriteString(bytes + " ")
		}
		tmpl.WriteString("</td><td>" + row.ASCII + "</td></tr>")
	}

	tmpl.WriteString("</table>")

	return template.HTML(tmpl.String())
}
