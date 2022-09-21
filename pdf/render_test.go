package pdf

import (
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	pdf := NewPdf()

	pdf.AddCell()

	testFile, err := os.Create("test.pdf")
	if err != nil {
		t.Error(err)
		return
	}
	defer func(testFile *os.File) {
		_ = testFile.Close()
	}(testFile)
	_, err = pdf.WriteTo(testFile)
	if err != nil {
		t.Error(err)
		return
	}
}
