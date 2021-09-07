package gag

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type rngContextKey string

var rngContextSeedKey = rngContextKey("seed")

type RNGEvent struct {
	Output struct {
		Execution
		RNG int
	}
}

func (r *RNGEvent) Run(ctx *Context) error {
	// ctxSeed := ctx.State[rngContextSeedKey].(int64)
	// rand.Seed(ctxSeed)

	r.Output.RNG = rand.Intn(100)

	return r.Output.Execution(ctx)
}

type Adder struct {
	Input struct {
		Number1 int
		Number2 int
	}

	Output struct {
		Sum int
	}
}

func (a *Adder) Run(ctx *Context) error {
	a.Output.Sum = a.Input.Number1 + a.Input.Number2
	return nil
}

type Comparer struct {
	Input struct {
		Number1 int
		Number2 int
	}

	Output struct {
		Greater Execution
		Less    Execution
	}
}

func (c *Comparer) Run(ctx *Context) error {
	if c.Input.Number1 > c.Input.Number2 {
		return c.Output.Greater(ctx)
	}

	return c.Output.Less(ctx)
}

type Panicer struct {
	Input struct {
		Execution
	}

	Output struct{}
}

func (p *Panicer) Run(ctx *Context) error {
	panic("ahhh")
}

type Printer struct {
	Input struct {
		Value interface{}
	}
}

func (p *Printer) Run(ctx *Context) error {
	//fmt.Println(p.Input.Value)
	return nil
}

func TestExample(t *testing.T) {
	Register(Adder{}, Comparer{}, Panicer{})

	g := &Graph{
		Nodes: []Node{
			{Name: "Adder"},
			{Name: "Comparer"},
			{Name: "Panicer"},
		},
		Edges: []Edge{
			{Output: Vertex{Raw: 1}, Input: Vertex{ID: 0, Field: "Number1"}},
			{Output: Vertex{Raw: 2}, Input: Vertex{ID: 0, Field: "Number2"}},
			{Output: Vertex{ID: 0, Field: "Sum"}, Input: Vertex{ID: 1, Field: "Number1"}},
			{Output: Vertex{Raw: 4}, Input: Vertex{ID: 1, Field: "Number2"}},
			{Output: Vertex{ID: 1, Field: "Greater"}, Input: Vertex{ID: 2}},
		},
	}

	drawGraph(g)

	if err := g.Run(1, nil); err != nil && !errors.Is(err, ErrGraphDone) {
		fmt.Println(err)
		t.Fail()
	}

}

func BenchmarkExample(b *testing.B) {
	Register(Adder{}, Comparer{}, Panicer{})

	g := &Graph{
		Nodes: []Node{
			{Name: "Adder"},
			{Name: "Comparer"},
			{Name: "Panicer"},
		},
		Edges: []Edge{
			{Output: Vertex{Raw: 1}, Input: Vertex{ID: 0, Field: "Number1"}},
			{Output: Vertex{Raw: 2}, Input: Vertex{ID: 0, Field: "Number2"}},
			{Output: Vertex{ID: 0, Field: "Sum"}, Input: Vertex{ID: 1, Field: "Number1"}},
			{Output: Vertex{Raw: 4}, Input: Vertex{ID: 1, Field: "Number2"}},
			{Output: Vertex{ID: 1, Field: "Greater"}, Input: Vertex{ID: 2}},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := g.Run(1, nil); err != nil && !errors.Is(err, ErrGraphDone) {
			fmt.Println(err)
			b.Fail()
		}
	}
}

func drawGraph(gagg *Graph) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}

	nodes := map[int]*cgraph.Node{}

	for id, n := range gagg.Nodes {
		gn, err := graph.CreateNode(n.Name)
		if err != nil {
			panic(err)
		}
		nodes[id] = gn
	}

	for _, e := range gagg.Edges {
		output := nodes[e.Output.ID]

		if e.Output.Field == "" {
			output, _ = graph.CreateNode(fmt.Sprintf("%+v", e.Output.Raw))
		}

		edge, err := graph.CreateEdge(e.Output.Field+" -> "+e.Input.Field, output, nodes[e.Input.ID])
		label := e.Output.Field + " -> " + e.Input.Field

		if e.Output.Field == "" {
			label = "Raw -> " + e.Input.Field
		}

		edge.SetLabel(label)

		if err != nil {
			panic(err)
		}
	}

	// 1. write encoded PNG data to buffer
	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		log.Fatal(err)
	}

	g.SetLayout(graphviz.TWOPI)
	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, "./graph.png"); err != nil {
		log.Fatal(err)
	}
}
