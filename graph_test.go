package gag

import (
	"context"
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	type args struct {
		nodes []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.nodes...)
		})
	}
}

func TestExecution_Next(t *testing.T) {
	type fields struct {
		output interface{}
		edge   Edge
		g      *Graph
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Execution{
				output: tt.fields.output,
				edge:   tt.fields.edge,
				g:      tt.fields.g,
			}
			if err := e.Next(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Execution.Next() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGraph_Run(t *testing.T) {
	type fields struct {
		Nodes []Node
		Edges []Edge
	}
	type args struct {
		nodeId int
		ctx    context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Graph{
				Nodes: tt.fields.Nodes,
				Edges: tt.fields.Edges,
			}
			if err := g.Run(tt.args.nodeId, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Graph.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGraph_fetchNodeOutput(t *testing.T) {
	type fields struct {
		Nodes []Node
		Edges []Edge
	}
	type args struct {
		ctx    context.Context
		nodeId int
		field  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Graph{
				Nodes: tt.fields.Nodes,
				Edges: tt.fields.Edges,
			}
			got, err := g.fetchNodeOutput(tt.args.ctx, tt.args.nodeId, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Graph.fetchNodeOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Graph.fetchNodeOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraph_findOutputEdge(t *testing.T) {
	type fields struct {
		Nodes []Node
		Edges []Edge
	}
	type args struct {
		id    int
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Edge
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Graph{
				Nodes: tt.fields.Nodes,
				Edges: tt.fields.Edges,
			}
			got, got1 := g.findOutputEdge(tt.args.id, tt.args.field)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Graph.findOutputEdge() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Graph.findOutputEdge() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGraph_findInputEdge(t *testing.T) {
	type fields struct {
		Nodes []Node
		Edges []Edge
	}
	type args struct {
		id    int
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Edge
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Graph{
				Nodes: tt.fields.Nodes,
				Edges: tt.fields.Edges,
			}
			got, got1 := g.findInputEdge(tt.args.id, tt.args.field)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Graph.findInputEdge() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Graph.findInputEdge() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGraph_buildNode(t *testing.T) {
	type fields struct {
		Nodes []Node
		Edges []Edge
	}
	type args struct {
		ctx    context.Context
		nodeId int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Runner
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Graph{
				Nodes: tt.fields.Nodes,
				Edges: tt.fields.Edges,
			}
			got, err := g.buildNode(tt.args.ctx, tt.args.nodeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Graph.buildNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Graph.buildNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
