package parser

import (
	"go_remote_control/collector/data"
	"regexp"
)

type Parse interface {
	GetLink(d *data.Data)
	GetImageLink(d *data.Data)
}

type Parser struct {
}

func (p *Parser) GetImageLink(d *data.Data) {
	reImageLink := regexp.MustCompile(`src="(.*?)"`)
	d.ParsedData = reImageLink.FindAllString(d.RawHtml, -1)
}

func (p *Parser) GetLink(d *data.Data) {
	reLink := regexp.MustCompile(`href="(.*?)"`)
	d.ParsedData = reLink.FindAllString(d.RawHtml, -1)
}
