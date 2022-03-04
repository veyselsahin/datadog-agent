package flowaggregator

import (
	"sync"

	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/DataDog/datadog-agent/pkg/netflow/common"
)

// flowAccumulator is used to accumulate aggregated flows
type flowAccumulator struct {
	flows map[string]*common.Flow
	mu    sync.Mutex
}

func newFlowAccumulator() *flowAccumulator {
	return &flowAccumulator{
		flows: make(map[string]*common.Flow),
	}
}

func (f *flowAccumulator) flush() []*common.Flow {
	f.mu.Lock()
	defer f.mu.Unlock()

	var flows []*common.Flow // init with optimal size
	for _, flow := range f.flows {
		flows = append(flows, flow)
	}

	// clear flows
	f.flows = make(map[string]*common.Flow)

	return flows
}

func (f *flowAccumulator) add(flowToAdd *common.Flow) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// TODO: handle port direction (see network-http-logger)
	// TODO: ignore ephemeral ports

	aggFlow, ok := f.flows[flowToAdd.AggregationHash()]
	log.Tracef("New Flow (digest=%s): %+v", flowToAdd.AggregationHash(), flowToAdd)
	if !ok {
		f.flows[flowToAdd.AggregationHash()] = flowToAdd
	} else {
		newAggFlow := *aggFlow
		newAggFlow.Bytes += flowToAdd.Bytes
		newAggFlow.Packets += flowToAdd.Packets
		newAggFlow.ReceivedTimestamp = minUint64(newAggFlow.ReceivedTimestamp, flowToAdd.ReceivedTimestamp)
		newAggFlow.StartTimestamp = minUint64(newAggFlow.StartTimestamp, flowToAdd.StartTimestamp)
		newAggFlow.EndTimestamp = maxUint64(newAggFlow.EndTimestamp, flowToAdd.EndTimestamp)

		// TODO: Cumulate TCPFlags (Cumulative of all the TCP flags seen for this flow)

		log.Tracef("Existing Aggregated Flow (digest=%s): %+v", flowToAdd.AggregationHash(), aggFlow)
		log.Tracef("New Aggregated Flow (digest=%s): %+v", flowToAdd.AggregationHash(), newAggFlow)
		f.flows[flowToAdd.AggregationHash()] = &newAggFlow
	}
}

func minUint64(a uint64, b uint64) uint64 {
	// TODO: TESTME
	if a < b {
		return a
	}
	return b
}

func maxUint64(a uint64, b uint64) uint64 {
	// TODO: TESTME
	if a > b {
		return a
	}
	return b
}
