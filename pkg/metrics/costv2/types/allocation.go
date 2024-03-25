package costv2

import (
	"time"
)

type CostType string

const (
	AllocationPretaxAmount      CostType = "allocation_pretax_amount"
	AllocationPretaxGrossAmount CostType = "allocation_pretax_gross_amount"
	CostEstimated               CostType = "cost_estimated"
)

type Allocation struct {
	Name       string                `json:"name"`
	Properties *AllocationProperties `json:"properties,omitempty"`
	//Window               *Window                `json:"window"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	//CPUCoreHours          float64   `json:"cpuCoreHours"`
	CPUCoreRequestAverage float64 `json:"cpuCoreRequestAverage"`
	CPUCoreUsageAverage   float64 `json:"cpuCoreUsageAverage"`
	//GPUHours               float64               `json:"gpuHours"`
	//RAMByteHours           float64 `json:"ramByteHours"`
	RAMBytesRequestAverage float64 `json:"ramByteRequestAverage"`
	RAMBytesUsageAverage   float64 `json:"ramByteUsageAverage"`
	Cost                   float64 `json:"cost"`
	CostRatio              float64 `json:"costRatio"`
	CustomCost             float64 `json:"customCost"`
}

type AllocationProperties struct {
	Controller     string            `json:"controller,omitempty"`
	ControllerKind string            `json:"controllerKind,omitempty"`
	Namespace      string            `json:"namespace,omitempty"`
	Pod            string            `json:"pod,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
}

//	type AllocationSet struct {
//		Allocations map[string]*Allocation `json:"allocations"`
//		Window      Window                 `json:"window"`
//		Type        string
//	}
type AllocationSet map[string]*Allocation

// NewAllocationSet instantiates a new AllocationSet
func NewAllocationSet() *AllocationSet {
	as := AllocationSet(make(map[string]*Allocation))
	return &as
}

// IsEmpty returns true if the AllocationSet is nil, or if it contains
// zero allocations.
func (as *AllocationSet) IsEmpty() bool {
	if as == nil || len(*as) == 0 {
		return true
	}

	return false
}

// Set uses the given Allocation to overwrite the existing entry in the
// AllocationSet under the Allocation's name.
func (as *AllocationSet) Set(alloc *Allocation) error {
	if as.IsEmpty() {
		*as = make(map[string]*Allocation)
	}

	(*as)[alloc.Name] = alloc

	return nil
}

func (as *AllocationSet) AggregateBy(aggregateBy []string) error {
	return nil
}

type AllocationSetRange struct {
	Allocations []*AllocationSet `json:"data"`
}

// NewAllocationSetRange instantiates a new range composed of the given
// AllocationSets in the order provided.
func NewAllocationSetRange(allocs ...*AllocationSet) *AllocationSetRange {
	return &AllocationSetRange{
		Allocations: allocs,
	}
}

// AggregateBy aggregates each AllocationSet in the range by the given
// properties and options.
func (asr *AllocationSetRange) AggregateBy(aggregateBy []string) error {
	return nil
}

// Append appends the given AllocationSet to the end of the range. It does not
// validate whether or not that violates window continuity.
func (asr *AllocationSetRange) Append(that *AllocationSet) {
	asr.Allocations = append(asr.Allocations, that)
}