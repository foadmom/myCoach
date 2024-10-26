package print

import (
	"bytes"
	"fmt"

	"github.com/go-pdf/fpdf"
)

// ============================================================================
// margints and sizes
// ============================================================================
var A4_WIDTH float64 = 210
var A4_HEIGHT float64 = 297

var DEFAULT_FONT string = "Arial"
var FONT_SIZE_XXS float64 = 6
var FONT_SIZE_XS float64 = 8
var FONT_SIZE_S float64 = 10
var FONT_SIZE_M float64 = 12
var FONT_SIZE_XM float64 = 13
var FONT_SIZE_L float64 = 14
var FONT_SIZE_XL float64 = 16
var FONT_SIZE_XXL float64 = 18

var QUADRANT_MARGIN float64 = MARGIN
var QUADRANT_WIDTH float64 = A4_WIDTH / 2
var QUADRANT_HEIGHT float64 = A4_HEIGHT / 2
var QUADRANT_TEXT_WIDTH float64 = QUADRANT_WIDTH - (QUADRANT_MARGIN * 2)
var QUADRANT_TEXT_HEIGHT float64 = QUADRANT_HEIGHT / (QUADRANT_MARGIN * 2)

var LINE_EXTENSION float64 = 2

var MARGIN float64 = 4
var LEFT_MARGIN float64 = MARGIN
var RIGHT_MARGIN float64 = MARGIN
var TOP_MARGIN float64 = MARGIN
var BOTTOM_MARGIN float64 = MARGIN

var TLQ_ORIGIN_X float64 = 0
var TLQ_ORIGIN_Y float64 = 0
var TLQ_TEXT_ORIGIN_X float64 = 0 + LEFT_MARGIN
var TLQ_TEXT_ORIGIN_Y float64 = 0 + TOP_MARGIN
var TLQ_TEXT_USABLE_W float64 = QUADRANT_WIDTH - (QUADRANT_MARGIN * 2)
var TLQ_TEXT_USABLE_H float64 = QUADRANT_HEIGHT - (QUADRANT_MARGIN * 2)

var TRQ_ORIGIN_X float64 = V_CENTRE
var TRQ_ORIGIN_Y float64 = TLQ_TEXT_ORIGIN_Y
var TRQ_TEXT_ORIGIN_X float64 = TRQ_ORIGIN_X + LEFT_MARGIN
var TRQ_TEXT_ORIGIN_Y float64 = TRQ_ORIGIN_Y + TOP_MARGIN

var BRQ_ORIGIN_X float64 = V_CENTRE
var BRQ_ORIGIN_Y float64 = H_CENTRE
var BRQ_TEXT_ORIGIN_X float64 = BRQ_ORIGIN_X + LEFT_MARGIN
var BRQ_TEXT_ORIGIN_Y float64 = BRQ_ORIGIN_Y + TOP_MARGIN

var TRAVEL_TIPS_HEIGHT float64 = 71

var YOUR_EXTRAS_Y float64 = 120

var MAX_X float64 = A4_WIDTH - RIGHT_MARGIN
var MAX_Y float64 = A4_HEIGHT - BOTTOM_MARGIN
var V_CENTRE float64 = A4_WIDTH / 2
var H_CENTRE float64 = A4_HEIGHT / 2

var FILL_COLOR []int = []int{230, 230, 230}      // dark grey
var DEFAULT_TEXT_COLOR []int = []int{64, 64, 64} // dark grey
var GRAY_TEXT_COLOR []int = []int{126, 126, 126} // dark grey
var NATIONAL_BLUE []int = []int{39, 109, 156}    // blue
var EXPRESS_RED []int = []int{153, 31, 0}        // red

type ticketInfo struct {
}

// ========================================================
//
// ========================================================
func InitPage() *fpdf.Fpdf {
	_pdf := fpdf.New("P", "mm", "A4", "")
	addPage(_pdf)
	return _pdf
}

// ========================================================
//
// ========================================================
func addPage(pdf *fpdf.Fpdf) {
	pdf.AddPage()
	SetMargins(pdf)
	pageNumber(pdf)
}

// ========================================================
//
// ========================================================
func pageNumber(pdf *fpdf.Fpdf) {
	pdf.SetFont(DEFAULT_FONT, "B", 8)
	pdf.SetXY(195, 5)
	_pageInfo := fmt.Sprintf("%d of %d", pdf.PageNo(), pdf.PageCount())
	pdf.Cell(12, 3, _pageInfo)
}

// ========================================================
// This is the dividing lines at the centre of the page to
// divide the page into 4 quadrant and mark the suggested
// folding position.
// ========================================================
func DrawDividingCentreLines(pdf *fpdf.Fpdf) {
	pdf.SetDashPattern([]float64{1, 1}, 1)
	pdf.SetDrawColor(195, 195, 195)
	pdf.Line(LEFT_MARGIN-LINE_EXTENSION, H_CENTRE, MAX_X+LINE_EXTENSION, H_CENTRE)
	pdf.Line(V_CENTRE, TOP_MARGIN-LINE_EXTENSION, V_CENTRE, MAX_Y+LINE_EXTENSION)
}

// ========================================================
//
// ========================================================
func SetFillColor(pdf *fpdf.Fpdf) {
	pdf.SetFillColor(FILL_COLOR[0], FILL_COLOR[1], FILL_COLOR[2])
}

// ========================================================
//
// ========================================================
func AddImage(pdf *fpdf.Fpdf, imageType string, imageName string, x, y, w, h float64) {
	var opt fpdf.ImageOptions
	opt.ImageType = imageType
	pdf.ImageOptions(imageName, x, y, w, h, false, opt, 0, "")
}

// ========================================================
//
// ========================================================
func SetMargins(pdf *fpdf.Fpdf) {
	pdf.SetMargins(LEFT_MARGIN, TOP_MARGIN, RIGHT_MARGIN)
}

// ========================================================
//
// ========================================================
func SetTextColor(pdf *fpdf.Fpdf, rgb []int) {
	pdf.SetTextColor(rgb[0], rgb[1], rgb[2])
}

// ========================================================
//
// ========================================================
func writeBasicMultiCell(pdf *fpdf.Fpdf, x, y, w, h float64, text string) float64 {
	pdf.SetXY(x, y)
	pdf.MultiCell(w, h, text, "0", "L", false)

	return pdf.GetY()
}

// ========================================================
//
// ========================================================
func NationalExpressHeading(pdf *fpdf.Fpdf, size, x, y, w, h float64) float64 {
	pdf.SetY(y)
	pdf.SetX(x)
	pdf.SetFont(DEFAULT_FONT, "B", size)
	pdf.SetTextColor(NATIONAL_BLUE[0], NATIONAL_BLUE[1], NATIONAL_BLUE[2])
	pdf.Cell(40, 12, "national ")
	pdf.SetTextColor(EXPRESS_RED[0], NATIONAL_BLUE[1], NATIONAL_BLUE[2])
	pdf.MultiCell(w, h, "express", "0", "L", false)
	return pdf.GetY()
}

// ========================================================
//
// ========================================================
func AddQRCodeFromBytes(pdf *fpdf.Fpdf, data []byte, x, y, w, h float64) float64 {
	_reader := bytes.NewReader(data)
	var _options fpdf.ImageOptions = fpdf.ImageOptions{ImageType: "PNG",
		ReadDpi: false, AllowNegativePosition: false}
	_ = pdf.RegisterImageReader("qrCode", "png", _reader)

	pdf.ImageOptions("qrCode", x, y, w, h, false, _options, 0, "")

	return pdf.GetY()
}
