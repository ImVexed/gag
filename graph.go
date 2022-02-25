package gag

import (
	"context"
	"errors"
	"strconv"

	"github.com/modern-go/reflect2"
)

type Runner interface {
	Run(ctx context.Context) error
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	Name  string
	State interface{}
}

type Edge struct {
	Output, Input Vertex
}

type Vertex struct {
	ID    int
	Field string
	Raw   interface{}
}

type Execution struct {
	output interface{}
	edge   Edge
	g      *Graph
}

var ErrGraphDone = errors.New("graph done")

type nodeType struct {
	t       reflect2.StructType
	output  reflect2.StructField
	input   reflect2.StructField
	outputs map[string]reflect2.StructField
	inputs  map[string]reflect2.StructField
}

var nodeTypes = map[string]*nodeType{}

func Register(nodes ...interface{}) {
	for _, n := range nodes {
		t := reflect2.TypeOf(n).(reflect2.StructType)

		inputs := map[string]reflect2.StructField{}
		outputs := map[string]reflect2.StructField{}

		inputField := t.FieldByName("Input").Type().(reflect2.StructType)
		outputField := t.FieldByName("Output").Type().(reflect2.StructType)

		for i := 0; i < inputField.NumField(); i++ {
			field := inputField.Field(i)
			inputs[field.Name()] = field
		}

		for i := 0; i < outputField.NumField(); i++ {
			field := outputField.Field(i)
			outputs[field.Name()] = field
		}

		nodeTypes[t.Type1().Name()] = &nodeType{
			t,
			t.FieldByName("Output"),
			t.FieldByName("Input"),
			outputs,
			inputs,
		}
	}
}

type nodeState map[interface{}]interface{}
type nodeStateContextKey int

var nodeStateKey = nodeStateContextKey(0)

func (e Execution) Next(ctx context.Context) error {
	if e.edge.Input.Field == "" {
		return ErrGraphDone
	}

	n := e.g.Nodes[e.edge.Output.ID]

	p := nodeTypes[n.Name]

	state := ctx.Value(nodeStateKey).(nodeState)

	// Snapshot the output field of the node to cache before we head to the next node
	if _, ok := state[outputContextKey(e.edge.Output.ID)]; !ok {
		state[outputContextKey(e.edge.Output.ID)] = p.output.Get(e.output)
	}

	// JIT the next node to be ran
	r, err := e.g.buildNode(ctx, e.edge.Input.ID)

	if err != nil {
		return err
	}

	return r.Run(ctx)
}

func (g *Graph) Run(nodeId int, ctx context.Context) error {
	ctx = context.WithValue(ctx, nodeStateKey, nodeState{})

	r, err := g.buildNode(ctx, nodeId)

	if err != nil {
		return err
	}

	return r.Run(ctx)
}

type outputContextKey int

func (g *Graph) fetchNodeOutput(ctx context.Context, nodeId int, field string) (interface{}, error) {
	nt := nodeTypes[g.Nodes[nodeId].Name]

	state := ctx.Value(nodeStateKey).(nodeState)

	// If the node has already been ran and it's output cached, use the cached value
	if v, ok := state[outputContextKey(nodeId)]; ok {
		return nt.outputs[field].Get(v), nil
	}

	// If the output of the node wasn't cached, then the node hasn't been built or ran, so we need to build it
	r, err := g.buildNode(ctx, nodeId)

	if err != nil {
		return nil, err
	}

	if err := r.Run(ctx); err != nil {
		return nil, err
	}

	output := nt.output.Get(r)
	// Cache the Output field
	state[outputContextKey(nodeId)] = output

	return nt.outputs[field].Get(output), nil
}

func (g *Graph) findOutputEdge(id int, field string) (Edge, bool) {
	for _, edge := range g.Edges {
		if edge.Output.ID == id && edge.Output.Field == field {
			return edge, true
		}
	}

	return Edge{}, false
}

func (g *Graph) findInputEdge(id int, field string) (Edge, bool) {
	for _, edge := range g.Edges {
		if edge.Input.ID == id && edge.Input.Field == field {
			return edge, true
		}
	}

	return Edge{}, false
}

func (g *Graph) buildNode(ctx context.Context, nodeId int) (Runner, error) {
	if len(g.Nodes) < nodeId {
		return nil, errors.New("out of bounds")
	}

	n := g.Nodes[nodeId]

	p, ok := nodeTypes[n.Name]

	if !ok {
		return nil, errors.New("Node " + n.Name + " has not been registerd")
	}

	// Pull/Create the node from it's pool
	v := p.t.New()

	// Register the execution outputs
	for name, field := range p.outputs {
		// Assert that the output field is not private, and is an Executor
		if field.Type().RType() != reflect2.RTypeOf(Execution{}) {
			continue
		}

		e, _ := g.findOutputEdge(nodeId, name)

		exe := field.Get(p.output.Get(v)).(*Execution)
		exe.edge = e
		exe.g = g
		exe.output = v
	}

	// Populate all of the inputs
	for name, field := range p.inputs {
		// Assert the field is public and isn't an executor
		if field.Type().RType() == reflect2.RTypeOf(Execution{}) {
			continue
		}

		e, ok := g.findInputEdge(nodeId, name)

		// Assert we either have raw input, or have
		if !ok {
			return nil, errors.New("Input " + name + " for node " + n.Name + " at index " + strconv.Itoa(nodeId) + " has no form of input")
		}

		// If there is a raw output for this input, assert it's assignable and set it
		if e.Output.Raw != nil {
			if !field.Type().AssignableTo(reflect2.TypeOf(e.Output.Raw)) {
				return nil, errors.New("Input " + name + " for node " + n.Name + " at index " + strconv.Itoa(nodeId) + " can't be processed because the raw output is not assignable to the input")
			}
			field.UnsafeSet(reflect2.PtrOf(p.input.Get(v)), reflect2.PtrOf(e.Output.Raw))
			continue
		}

		// Otherwise, get the relevant value from the node output our edge is linked to
		otf, err := g.fetchNodeOutput(ctx, e.Output.ID, e.Output.Field)

		if err != nil {
			return nil, errors.New("Input " + name + " for node " + n.Name + " at index " + strconv.Itoa(nodeId) + " can't be processed because " + err.Error())
		}

		// Sanity check the value is assignable
		if !field.Type().AssignableTo(reflect2.TypeOfPtr(otf).Elem()) {
			return nil, errors.New("Input " + name + " for node " + n.Name + " at index " + strconv.Itoa(nodeId) + " can't be processed because output " + e.Output.Field + " for node " + g.Nodes[e.Output.ID].Name + " at index " + strconv.Itoa(e.Output.ID) + " is not assignable to the input")
		}

		field.UnsafeSet(reflect2.PtrOf(p.input.Get(v)), reflect2.PtrOf(otf))
	}

	return v.(Runner), nil
}
