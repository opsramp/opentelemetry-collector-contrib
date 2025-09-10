// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pipeline // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/pipeline"

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/multierr"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	stanzaerrors "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/errors"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
)

var _ Pipeline = (*DirectedPipeline)(nil)

var alreadyStarted = stanzaerrors.NewError("pipeline already started", "")
var alreadyStopped = stanzaerrors.NewError("pipeline already stopped", "")

// DirectedPipeline is a pipeline backed by a directed graph
type DirectedPipeline struct {
	Graph     *simple.DirectedGraph
	startOnce sync.Once
	stopOnce  sync.Once
}

// Start will start the operators in a pipeline in reverse topological order
func (p *DirectedPipeline) Start(persister operator.Persister) error {
	var err error = alreadyStarted
	p.startOnce.Do(func() {
		err = p.start(persister)
	})
	pipelineerr := p.PrintPipeline("/var/log/opsramp/pipeline.txt")
	if pipelineerr != nil {
		// Handle error
		fmt.Println("Failed to print pipeline: %v", pipelineerr)
	}
	return err
}

// PrintPipeline writes a human-readable representation of the pipeline to a file
func (p *DirectedPipeline) PrintPipeline(filePath string) error {
	// Create file with appropriate permissions
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	sortedNodes, err := topo.Sort(p.Graph)
	if err != nil {
		return fmt.Errorf("failed to sort pipeline: %w", err)
	}

	// Write header
	fmt.Fprintf(file, "Pipeline Structure\n")
	fmt.Fprintf(file, "=================\n\n")
	fmt.Fprintf(file, "Total Operators: %d\n\n", len(sortedNodes))

	// Write operators in topological order
	fmt.Fprintf(file, "Operators (in processing order):\n")
	fmt.Fprintf(file, "--------------------------\n")

	for i, node := range sortedNodes {
		opNode := node.(OperatorNode)
		op := opNode.Operator()

		fmt.Fprintf(file, "[%d] Operator ID: %s\n", i+1, op.ID())
		fmt.Fprintf(file, "    Type: %T\n", op)

		// Get outgoing connections
		outputs := p.getOutputOperators(opNode)
		if len(outputs) > 0 {
			fmt.Fprintf(file, "    Outputs To: ")
			for j, output := range outputs {
				if j > 0 {
					fmt.Fprintf(file, ", ")
				}
				fmt.Fprintf(file, "%s", output)
			}
			fmt.Fprintf(file, "\n")
		} else {
			fmt.Fprintf(file, "    Outputs To: <none>\n")
		}

		// Get incoming connections
		inputs := p.getInputOperators(opNode)
		if len(inputs) > 0 {
			fmt.Fprintf(file, "    Inputs From: ")
			for j, input := range inputs {
				if j > 0 {
					fmt.Fprintf(file, ", ")
				}
				fmt.Fprintf(file, "%s", input)
			}
			fmt.Fprintf(file, "\n")
		} else {
			fmt.Fprintf(file, "    Inputs From: <none>\n")
		}

		fmt.Fprintf(file, "\n")
	}

	// Write flow diagram
	fmt.Fprintf(file, "Flow Diagram:\n")
	fmt.Fprintf(file, "------------\n")
	for i, node := range sortedNodes {
		opNode := node.(OperatorNode)
		outputs := p.getOutputOperators(opNode)

		if len(outputs) > 0 {
			for _, output := range outputs {
				fmt.Fprintf(file, "%s --> %s\n", opNode.Operator().ID(), output)
			}
		} else if i < len(sortedNodes)-1 {
			fmt.Fprintf(file, "%s (no direct connections)\n", opNode.Operator().ID())
		} else {
			fmt.Fprintf(file, "%s (final operator)\n", opNode.Operator().ID())
		}
	}

	fmt.Fprintf(file, "Render Output:\n")
	a, _ := p.Render()
	fmt.Fprintf(file, string(a)+"\n")

	return nil
}

// getOutputOperators returns the IDs of operators that this node outputs to
func (p *DirectedPipeline) getOutputOperators(node OperatorNode) []string {
	var outputs []string

	// Iterate through edges from this node
	from := p.Graph.From(node.ID())
	for from.Next() {
		toNode := from.Node().(OperatorNode)
		outputs = append(outputs, toNode.Operator().ID())
	}

	return outputs
}

// getInputOperators returns the IDs of operators that input to this node
func (p *DirectedPipeline) getInputOperators(node OperatorNode) []string {
	var inputs []string

	// Iterate through edges to this node
	to := p.Graph.To(node.ID())
	for to.Next() {
		fromNode := to.Node().(OperatorNode)
		inputs = append(inputs, fromNode.Operator().ID())
	}

	return inputs
}

// Stop will stop the operators in a pipeline in topological order
func (p *DirectedPipeline) Stop() error {
	var err error = alreadyStopped
	p.stopOnce.Do(func() {
		err = p.stop()
	})
	return err
}

func (p *DirectedPipeline) start(persister operator.Persister) error {
	sortedNodes, _ := topo.Sort(p.Graph)
	for i := len(sortedNodes) - 1; i >= 0; i-- {
		op := sortedNodes[i].(OperatorNode).Operator()

		scopedPersister := operator.NewScopedPersister(op.ID(), persister)
		op.Logger().Debug("Starting operator")
		if err := op.Start(scopedPersister); err != nil {
			return err
		}
		op.Logger().Debug("Started operator")
	}

	return nil
}

func (p *DirectedPipeline) stop() error {
	var err error
	sortedNodes, _ := topo.Sort(p.Graph)
	for _, node := range sortedNodes {
		operator := node.(OperatorNode).Operator()
		operator.Logger().Debug("Stopping operator")
		if opErr := operator.Stop(); opErr != nil {
			err = multierr.Append(err, opErr)
		}
		operator.Logger().Debug("Stopped operator")
	}
	return err
}

// Render will render the pipeline as a dot graph
func (p *DirectedPipeline) Render() ([]byte, error) {
	return dot.Marshal(p.Graph, "G", "", " ")
}

// Operators returns a slice of operators that make up the pipeline graph
func (p *DirectedPipeline) Operators() []operator.Operator {
	var operators []operator.Operator
	if nodes, err := topo.Sort(p.Graph); err == nil {
		for _, node := range nodes {
			operators = append(operators, node.(OperatorNode).Operator())
		}
		return operators
	}

	// If for some unexpected reason an Unorderable error is returned,
	// when using topo.Sort, return the list without ordering
	nodes := p.Graph.Nodes()
	for nodes.Next() {
		operators = append(operators, nodes.Node().(OperatorNode).Operator())
	}
	return operators
}

// addNodes will add operators as nodes to the supplied graph.
func addNodes(graph *simple.DirectedGraph, operators []operator.Operator) error {
	for _, operator := range operators {
		operatorNode := createOperatorNode(operator)
		if graph.Node(operatorNode.ID()) != nil {
			return stanzaerrors.NewError(
				fmt.Sprintf("operator with id '%s' already exists in pipeline", operatorNode.Operator().ID()),
				"ensure that each operator has a unique `type` or `id`",
			)
		}

		graph.AddNode(operatorNode)
	}
	return nil
}

// connectNodes will connect the nodes in the supplied graph.
func connectNodes(graph *simple.DirectedGraph) error {
	nodes := graph.Nodes()
	for nodes.Next() {
		node := nodes.Node().(OperatorNode)
		if err := connectNode(graph, node); err != nil {
			return err
		}
	}

	if _, err := topo.Sort(graph); err != nil {
		var topoErr topo.Unorderable
		errors.As(err, &topoErr)
		return stanzaerrors.NewError(
			"pipeline has a circular dependency",
			"ensure that all operators are connected in a straight, acyclic line",
			"cycles", unorderableToCycles(topoErr),
		)
	}

	return nil
}

// connectNode will connect a node to its outputs in the supplied graph.
func connectNode(graph *simple.DirectedGraph, inputNode OperatorNode) error {
	for outputOperatorID, outputNodeID := range inputNode.OutputIDs() {
		if graph.Node(outputNodeID) == nil {
			return stanzaerrors.NewError(
				"operators cannot be connected, because the output does not exist in the pipeline",
				"ensure that the output operator is defined",
				"input_operator", inputNode.Operator().ID(),
				"output_operator", outputOperatorID,
			)
		}

		outputNode := graph.Node(outputNodeID).(OperatorNode)
		if !outputNode.Operator().CanProcess() {
			return stanzaerrors.NewError(
				"operators cannot be connected, because the output operator can not process logs",
				"ensure that the output operator can process logs (like a parser or destination)",
				"input_operator", inputNode.Operator().ID(),
				"output_operator", outputOperatorID,
			)
		}

		if graph.HasEdgeFromTo(inputNode.ID(), outputNodeID) {
			return stanzaerrors.NewError(
				"operators cannot be connected, because a connection already exists",
				"ensure that only a single connection exists between the two operators",
				"input_operator", inputNode.Operator().ID(),
				"output_operator", outputOperatorID,
			)
		}

		edge := graph.NewEdge(inputNode, outputNode)
		graph.SetEdge(edge)
	}

	return nil
}

// setOperatorOutputs will set the outputs on operators that can output.
func setOperatorOutputs(operators []operator.Operator) error {
	for _, operator := range operators {
		if !operator.CanOutput() {
			continue
		}

		if err := operator.SetOutputs(operators); err != nil {
			return stanzaerrors.WithDetails(err, "operator_id", operator.ID())
		}
	}
	return nil
}

// NewDirectedPipeline creates a new directed pipeline
func NewDirectedPipeline(operators []operator.Operator) (*DirectedPipeline, error) {
	if err := setOperatorOutputs(operators); err != nil {
		return nil, err
	}

	graph := simple.NewDirectedGraph()
	if err := addNodes(graph, operators); err != nil {
		return nil, err
	}

	if err := connectNodes(graph); err != nil {
		return nil, err
	}

	return &DirectedPipeline{Graph: graph}, nil
}

func unorderableToCycles(err topo.Unorderable) string {
	var cycles strings.Builder
	for i, cycle := range err {
		if i != 0 {
			cycles.WriteByte(',')
		}
		cycles.WriteByte('(')
		for _, node := range cycle {
			cycles.WriteString(node.(OperatorNode).operator.ID())
			cycles.Write([]byte(` -> `))
		}
		cycles.WriteString(cycle[0].(OperatorNode).operator.ID())
		cycles.WriteByte(')')
	}
	return cycles.String()
}
