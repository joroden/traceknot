package flags

import (
	"fmt"
	"sort"

	"traceknot/internal/export"
)

type Kind string

const (
	KindCostOutlier     Kind = "cost_outlier"
	KindDurationOutlier Kind = "duration_outlier"
	KindError           Kind = "error"
	KindRepeatedCall    Kind = "repeated_call"
)

type Flag struct {
	Kind    Kind
	NodeIDs []string
	Reason  string
}

const (
	outlierZScore  = 3.5
	minOutlierPeer = 4
	maxOutliers    = 20
	minLoopCount   = 3
)

func Compute(session *export.Session) []Flag {
	var result []Flag
	result = append(result, costOutliers(session.Nodes)...)
	result = append(result, durationOutliers(session.Nodes)...)
	result = append(result, errors(session.Nodes)...)
	result = append(result, repeatedCalls(session.Nodes)...)
	return result
}

func errors(nodes []*export.Node) []Flag {
	var result []Flag
	for _, node := range nodes {
		if node.Kind != "tool_call" || node.Status == nil || *node.Status != "error" {
			continue
		}
		result = append(result, Flag{
			Kind:    KindError,
			NodeIDs: []string{node.NodeID},
			Reason:  "tool call failed (status=error)",
		})
	}
	return result
}

func repeatedCalls(nodes []*export.Node) []Flag {
	type group struct {
		toolName string
		nodeIDs  []string
	}
	groups := make(map[string]*group)
	var order []string
	for _, node := range nodes {
		if node.Kind != "tool_call" || node.ToolName == nil {
			continue
		}
		parent := ""
		if node.ParentNodeID != nil {
			parent = *node.ParentNodeID
		}
		args := ""
		if node.Detail != nil && node.Detail.ToolCall != nil && node.Detail.ToolCall.ArgumentsJSON != nil {
			args = *node.Detail.ToolCall.ArgumentsJSON
		}
		key := parent + "\x00" + *node.ToolName + "\x00" + args
		g, ok := groups[key]
		if !ok {
			g = &group{toolName: *node.ToolName}
			groups[key] = g
			order = append(order, key)
		}
		g.nodeIDs = append(g.nodeIDs, node.NodeID)
	}

	var result []Flag
	for _, key := range order {
		g := groups[key]
		if len(g.nodeIDs) < minLoopCount {
			continue
		}
		result = append(result, Flag{
			Kind:    KindRepeatedCall,
			NodeIDs: g.nodeIDs,
			Reason:  fmt.Sprintf("same %s call repeated %d× under the same parent with identical arguments — possible loop", g.toolName, len(g.nodeIDs)),
		})
	}
	sort.Slice(result, func(i, j int) bool { return len(result[i].NodeIDs) > len(result[j].NodeIDs) })
	return result
}

func costOutliers(nodes []*export.Node) []Flag {
	return metricOutliers(nodes, KindCostOutlier, "cost", func(node *export.Node) float64 {
		if node.Kind == "agent" {
			return node.AggCost
		}
		return node.Cost
	}, func(value float64) string {
		return fmt.Sprintf("$%.4f", value)
	})
}

func durationOutliers(nodes []*export.Node) []Flag {
	return metricOutliers(nodes, KindDurationOutlier, "duration", func(node *export.Node) float64 {
		if node.DurationMs == nil {
			return 0
		}
		return *node.DurationMs
	}, func(value float64) string {
		return fmt.Sprintf("%.1fs", value/1000)
	})
}

func metricOutliers(
	nodes []*export.Node,
	kind Kind,
	metricName string,
	value func(node *export.Node) float64,
	format func(float64) string,
) []Flag {
	byKind := make(map[string][]*export.Node)
	for _, node := range nodes {
		if value(node) > 0 {
			byKind[node.Kind] = append(byKind[node.Kind], node)
		}
	}

	var result []Flag
	for nodeKind, peers := range byKind {
		if len(peers) < minOutlierPeer {
			continue
		}
		values := make([]float64, len(peers))
		for i, node := range peers {
			values[i] = value(node)
		}
		med := median(values)
		mad := medianAbsoluteDeviation(values, med)

		type candidate struct {
			node  *export.Node
			value float64
			z     float64
		}
		var candidates []candidate
		for _, node := range peers {
			v := value(node)
			var z float64
			if mad > 0 {
				z = 0.6745 * (v - med) / mad
			} else if v > med {
				z = outlierZScore
			}
			if z >= outlierZScore {
				candidates = append(candidates, candidate{node: node, value: v, z: z})
			}
		}
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].value > candidates[j].value })
		if len(candidates) > maxOutliers {
			candidates = candidates[:maxOutliers]
		}
		for _, c := range candidates {
			result = append(result, Flag{
				Kind:    kind,
				NodeIDs: []string{c.node.NodeID},
				Reason: fmt.Sprintf("%s %s %s vs. this session's %s median of %s",
					nodeKind, metricName, format(c.value), nodeKind, format(med)),
			})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Reason < result[j].Reason })
	return result
}

func median(values []float64) float64 {
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

func medianAbsoluteDeviation(values []float64, med float64) float64 {
	deviations := make([]float64, len(values))
	for i, v := range values {
		deviations[i] = abs(v - med)
	}
	return median(deviations)
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
