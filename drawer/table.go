package drawer

import (
	"log"
	"strings"
)

//minPrintItemLen - minimal content length in row (for example - 0.00% is 5 symbols)
//used to calculate spaces between cells
const minPrintItemLen = 5

//Table is used to draw info in "table" form to stdout
type Table struct {
	Header []string
}

//DrawHeader will print table header
func (t *Table) DrawHeader() {
	log.Printf(strings.Join(t.Header, "    "))
}

//DrawRow will print table row
func (t *Table) DrawRow(cells []string) {
	log.Println(t.formatRow(cells))
}

//formatRow is used to prepare row data (calculate spaces count between cells)
func (t *Table) formatRow(cells []string) string {
	emptyString := func(len int) string {
		return strings.Join(make([]string, len), " ")
	}

	//build `        ` separator based on items len (to make sure that each element will shown one-by-one)
	result := make([]string, 0)
	for i := 0; i < len(cells)-1; i++ {
		sepLen := len(t.Header[i]) - (len(cells[i]) - len(cells[i+1]))
		sepLen -= len(cells[i+1]) - minPrintItemLen
		result = append(result, cells[i]+emptyString(sepLen))
	}

	//append last item
	result = append(result, cells[len(cells)-1])
	return strings.Join(result, "")
}
