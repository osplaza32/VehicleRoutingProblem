package transformations

import (
	"ixior/VehicleRoutingProblem/common"
	"ixior/VehicleRoutingProblem/model"
	"math"
)

// ShortRouteImprovementTransformation Remove short routes and then find new places for nodes
func ShortRouteImprovementTransformation(data *model.CaseDTO, inputSolution *model.Solution) *model.Solution {
	newSolution := common.CloneSolution(inputSolution)

	// Get max length of routes in solution
	maxRouteLength := 0
	for i := 0; i < len(newSolution.Routes); i++ {
		if len(newSolution.Routes[i]) > maxRouteLength {
			maxRouteLength = len(newSolution.Routes[i])
		}
	}

	// Remove short routes
	unroutedNodes := make([]int, 0)
	for i := 0; i < len(newSolution.Routes); i++ {
		if len(newSolution.Routes[i]) <= maxRouteLength/2 {
			for j := 0; j < len(newSolution.Routes[i]); j++ {
				unroutedNodes = append(unroutedNodes, newSolution.Routes[i][j])
			}
			newSolution.Routes = common.RemoveRoute(newSolution.Routes, i)
			i--
		}
	}

	common.Shuffle(unroutedNodes)

	// Find new place for unrouted nodes
	for _, node := range unroutedNodes {
		if common.CanBePutToRoutes(newSolution, data, node) {
			lowestCost := math.MaxFloat64
			minI, minJ := -1, -1
			for i := 0; i < len(newSolution.Routes); i++ {
				if common.GetCapacity(newSolution.Routes[i], data)+data.GetDemand(node) <= data.Capacity {
					for j := 0; j <= len(newSolution.Routes[i]); j++ {
						node1, node2 := 0, 0
						if j != 0 {
							node1 = newSolution.Routes[i][j-1]
						}
						if j != len(newSolution.Routes[i]) {
							node2 = newSolution.Routes[i][j]
						}
						cost := data.Cost[node1][node] + data.Cost[node][node2]
						if cost < lowestCost {
							lowestCost, minI, minJ = cost, i, j
						}
					}
				}
			}

			// Insert removed node
			newSolution.Routes[minI] = common.AppendAtIndex(newSolution.Routes[minI], minJ, node)
		} else {
			newSolution.AddRoute([]int{node})
		}
	}

	return newSolution
}
