package print

// 多行进度条

import (
	"fmt"
)

type print struct {
	Tag   string `json:"tag"`
	Index int    `json:"index"`
}

func (this *print) Print(pos int, out string) {

	npos := pos-this.Index
	if npos > 0 {
		fmt.Printf("\033[%dA", npos)
	}

	fmt.Printf("\033[K %s \n", out)

	if npos > 0 {
		fmt.Printf("\033[%dB", npos)
	}
}

type Print struct {
	keys        map[string]int
	progresses  []*print
	CurrPostion int //当前显示最后一次刷新在第几行
}

func NewPrint() *Print {
	return &Print{keys: make(map[string]int)}
}

func (p *Print) Add(tag, out string) (index int) {
	p1 := &print{Index: len(p.progresses), Tag: tag}
	p.progresses = append(p.progresses, p1)
	p1.Print(p.CurrPostion, out)
	p.CurrPostion = p1.Index + 1
	p.keys[tag] = p1.Index
	return p1.Index
}

func (p *Print) Print(index int, out string) {
	p1 := p.progresses[index]
	p1.Print(p.CurrPostion, out)
	//p.CurrPostion = p1.Index + 1
}

func (p *Print) Print2(tag, out string) {
	p1 := p.progresses[p.keys[tag]]
	p1.Print(p.CurrPostion, out)
	//p.CurrPostion = p1.Index + 1
}
