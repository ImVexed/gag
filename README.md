# GAG - A Directed-Acyclic-Graph JIT in Go  

GAG is a library I created while developing https://isobot.io to experiment with different ways of implementing the core runtime.  

It intends to be a fast, highly parallel, DAG JIT, while still maintaining the balance between performance and usability.  

While the runtime is included in this library, a significantly more complex type system would need to be implemented ontop of GAG before it would likely be useful for any such similar use case.  

## Concepts
Similar to traditional DAGs, there are 4 fundamental primatives in GAG:  

- Graph
	- A collection of Nodes & Edges
- Vertex
	- A field on a node
- Edge
	- A connection between two vertex
- Node
	- A unit of work with inputs and outputs, analagous to a function

### **Execution**
Unlike traditional Flow-Based Programming, GAG includes the concept of "executing" or "running" a node. This fundamentally controls the flow of the graph's execution. Nodes can include an `Execution` type in their `Output` field which will be linked to the next node at runtime. 

### **Caching**
GAG also caches the `Output` fields of any node that has already been run. In the future when the JIT is improved this may only conditionally occur when it is optimal.

## Examples
A rather trival graph that simply adds 2 numbers together to compare the sum and panic depending on the result could be represented as such:

```go
&Graph{
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
```

Where the implementation of the `Adder` and `Comparer` nodes look like:

```go
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
```

The graph can also be visually represented like:  
![graph](/docs/graph.png)  
Or in isobot.io:
![iso-graph](/docs/iso-graph.png)  


# TODO
- Allow tainting specific outputs to invalidate downstream caches
- Additional node lifecycle methods

# WARNINGS
There is lots of unsafe non-standard reflection going on in this project. It is not nearly tested enough to be used in production, and there are likely many dragons hiding in the code itself.
