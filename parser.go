package main

import (
	"strings"

	"github.com/alecthomas/participle"
)

type Prime struct {
	L *Lit  `( @@ |`
	J *Junc `"(" @@ ")" )`
}

type Lit struct {
	Neg string `@"!"?`
	Sym string `@Ident`
}

type Junc struct {
	L  *Prime `@@`
	Op string `@("&" | "|" | "->" | "<-" | "<->")`
	R  *Prime `@@`
}

func (p *Prime) String() string {
	if p.L != nil {
		return p.L.String()
	}
	return p.J.String()
}
func (p *Lit) String() string {
	return p.Neg + p.Sym
}
func (p *Junc) String() string {
	return "(" + p.L.String() + p.Op + p.R.String() + ")"
}

func (p *Prime) Form(symbols *[]string) Formula {
	if p.L != nil {
		return p.L.Form(symbols)
	}
	return p.J.Form(symbols)
}
func (p *Lit) Form(symbols *[]string) Formula {
	i := -1
	for k, v := range *symbols {
		if v == p.Sym {
			i = k
			break
		}
	}
	if i == -1 {
		i = len(*symbols)
		*symbols = append(*symbols, p.Sym)
	}

	if p.Neg == "!" {
		return Not{S(i)}
	}
	return S(i)
}
func (p *Junc) Form(symbols *[]string) Formula {
	switch p.Op {
	case "&":
		return And{p.L.Form(symbols), p.R.Form(symbols)}
	case "|":
		return Or{p.L.Form(symbols), p.R.Form(symbols)}
	case "->":
		return Imp{p.L.Form(symbols), p.R.Form(symbols)}
	case "<-":
		return Imp{p.R.Form(symbols), p.L.Form(symbols)}
	case "<->":
		return BiImp{p.L.Form(symbols), p.R.Form(symbols)}
	}
	return nil
}

func Parse(s string) (Prime, error) {
	s = strings.TrimSpace(s)
	if s[0] != '(' {
		s = "(" + s
	}
	if s[len(s)-1] != ')' {
		s += ")"
	}
	parser := participle.MustBuild(&Prime{})
	//fmt.Println(parser)
	var e Prime
	err := parser.ParseString(s, &e)
	return e, err
}
