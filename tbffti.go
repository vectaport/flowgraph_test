package main

import (
	"math"
	"math/rand"

	"github.com/vectaport/fgbase"
	"github.com/vectaport/fgbase/imglab"
)

const infitesimal = 1.e-15

func tbiFire(n *fgbase.Node) error {
	x := n.Dsts[0]
	const sz = 128
	var vec = make([]complex128, sz, sz)
	rand.Seed(0x1515)

	delta := 3 * 2 * math.Pi / float64(sz)
	domain := float64(0)

	for i := range vec {
		vec[i] = complex(math.Sin(domain), 0.0)
		domain += delta
	}
	x.DstPut(vec)
	return nil
}

func tbi(x fgbase.Edge) fgbase.Node {
	node := fgbase.MakeNode("tbi", nil, []*fgbase.Edge{&x}, nil, tbiFire)
	return node
}

func tboFire(n *fgbase.Node) error {
	a := n.Srcs[0]
	b := n.Srcs[1]
	av := a.SrcGet().([]complex128)
	bv := b.SrcGet().([]complex128)
	if len(av) == len(bv) {
		for i := range av {
			if (real(av[i])-real(bv[i])) < -infitesimal || (real(av[i])-real(bv[i])) > infitesimal ||
				(imag(av[i])-imag(bv[i])) < -infitesimal || (imag(av[i])-imag(bv[i])) > infitesimal {
				n.Tracef("!SAME:  for %d delta is %v\n", i, av[i]-bv[i])
				n.Tracef("!SAME:  a = %v,  b = %v\n", av[i], bv[i])
				return nil
			}
		}
		n.Tracef("SAME all differences smaller than %v\n", infitesimal)
		return nil
	}
	n.Tracef("!SAME:  different sizes\n")
	return nil
}

func tbo(a, b fgbase.Edge) fgbase.Node {
	node := fgbase.MakeNode("tbo", []*fgbase.Edge{&a, &b}, nil, nil, tboFire)
	return node
}

func main() {

	fgbase.ConfigByFlag(nil)

	e, n := fgbase.MakeGraph(9, 7)

	e[7].Const(false)
	e[8].Const(true)

	n[0] = tbi(e[0])

	n[1] = fgbase.FuncFork(e[0], e[1], e[2])

	n[2] = imglab.FuncFFT(e[1], e[7], e[3])
	n[3] = fgbase.FuncPass(e[2], e[4])

	n[4] = imglab.FuncFFT(e[3], e[8], e[5])
	n[5] = fgbase.FuncPass(e[4], e[6])

	n[6] = tbo(e[5], e[6])

	fgbase.RunAll(n)

}
