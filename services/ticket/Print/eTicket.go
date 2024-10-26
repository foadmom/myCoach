package print

import (
	"fmt"

	"github.com/go-pdf/fpdf"
)

// ========================================================
//
// ========================================================
func Ticket_T1(ticket *ticketInfo) {
	var _pdf *fpdf.Fpdf = InitPage()
	DrawDividingCentreLines(_pdf)
	TopLeftQuadrant(_pdf, ticket)

	TopRightQuadrant(_pdf, ticket, 108, TOP_MARGIN)

	BottomRightQuadrant(_pdf)
	_err := _pdf.OutputFileAndClose("hello.pdf")
	if _err != nil {
		fmt.Printf("%s\n", _err)
	}
}
