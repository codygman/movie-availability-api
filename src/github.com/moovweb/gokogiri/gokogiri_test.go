package gokogiri

import (
	"github.com/moovweb/gokogiri/help"
	"testing"
)

func TestParseHtml(t *testing.T) {
	input := "<html><body><div><h1></div>"
	expected := `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">
<html><body><div><h1></h1></div></body></html>
`
	doc, err := ParseHtml([]byte(input))
	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}
	if doc.String() != expected {
		t.Error("the output of the html doc does not match the expected")
	}

	expected = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">
<html>
<head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head>
<body><div><h1></h1></div></body>
</html>
`
	doc.Root().FirstChild().AddPreviousSibling("<head></head>")

	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the html doc does not match the expected")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestParseXml(t *testing.T) {
	input := "<foo></foo>"
	expected := `<?xml version="1.0" encoding="utf-8"?>
<foo/>
`
	doc, err := ParseXml([]byte(input))
	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	if doc.String() != expected {
		t.Error("the output of the xml doc does not match the expected")
	}

	expected = `<?xml version="1.0" encoding="utf-8"?>
<foo>
  <bar/>
</foo>
`
	doc.Root().AddChild("<bar/>")
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match the expected")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}
