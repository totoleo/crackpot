package pdf

import (
	"io"
	"math"

	"github.com/phpdave11/gofpdf"
)

type Pdf struct {
	printer  *gofpdf.Fpdf
	fontSize float64
}

func NewPdf() *Pdf {
	f := gofpdf.New("Portrait", "mm", "A4", "./fonts")
	f.AddUTF8Font("WenDingPLJianBaoSong", "", "")
	fontSize := float64(14)
	f.SetFont("WenDingPLJianBaoSong", "", fontSize)
	return &Pdf{printer: f, fontSize: fontSize}
}

func (p *Pdf) AddCell() {
	p.printer.AddPage()
	p.helperLine()

	lineHeight := p.printer.PointConvert(p.fontSize)

	p.printer.MultiCell(190, lineHeight, "测试", "", "C", false)
	p.printer.MultiCell(190, lineHeight, "gopdf is a simple library for generating PDF document written in Go lang. A minimum version of Go 1.13 is required.Unicode subfont embedding. (汉字, Japanese, Korean, etc.) Features: Draw line, oval, rect, curve; Draw image (jpg, png); Set image mask; Password protection; Font kerning\n新的一行", "", "L", false)
	p.printer.CellFormat(190, lineHeight, "居左", "", -1, "L", false, 0, "")
	p.printer.CellFormat(0, lineHeight, "居右", "", -1, "R", false, 0, "")
	p.printer.CellFormat(-190, lineHeight, "居中", "", 1, "C", false, 0, "")

	ptSize, _ := p.printer.GetFontSize()
	ratio := 2.7
	newSize := ptSize * ratio
	leftRatio := (math.Ceil(ratio) - ratio) / 2
	p.printer.Ln(leftRatio * lineHeight)
	p.printer.SetFontSize(newSize)
	p.printer.CellFormat(190, lineHeight*ratio, "新的征程", "", 2, "C", false, 0, "")
	p.printer.Ln(leftRatio * lineHeight)
	p.printer.SetFontSize(ptSize)

	p.printer.MultiCell(0, lineHeight, "gofpdf supports UTF-8 TrueType fonts and “right-to-left” languages. Note that Chinese, Japanese, and Korean characters may not be included in many general purpose fonts. For these languages, a specialized font (for example, NotoSansSC for simplified Chinese) can be used.", "", "", false)

}
func (p *Pdf) helperLine() {
	lineHeight := p.printer.PointConvert(p.fontSize)
	p.printer.SetLineWidth(0.1)
	p.printer.SetLineJoinStyle("bevel")
	_, ht, _ := p.printer.PageSize(p.printer.PageNo())
	for i := lineHeight; i < ht; i += lineHeight {
		p.printer.Line(0, i, 585, i)
	}
}

func (p *Pdf) WriteTo(w io.Writer) (int64, error) {
	err := p.printer.Output(w)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
