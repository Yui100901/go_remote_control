package parser

import (
	"go_remote_control/collector/data"
	"regexp"
)

type Rule struct {
	Regex []string
}

type Parse interface {
	Preprocess(d *data.Data)
	GetLink(d *data.Data)
	GetImageLink(d *data.Data)
}

type Parser struct {
	*Rule
}

func (p *Parser) Preprocess(d *data.Data) {

}

func (p *Parser) GetImageLink(d *data.Data) {
	reImageLink := regexp.MustCompile(`src="(.*?)"`)
	d.ParsedData = reImageLink.FindAllString(d.RawHtml, -1)
}

func (p *Parser) GetLink(d *data.Data) {
	reLink := regexp.MustCompile(`href="(.*?)"`)
	d.ParsedData = reLink.FindAllString(d.RawHtml, -1)
}
