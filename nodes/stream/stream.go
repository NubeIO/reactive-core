package stream

//
//type AnyBufferNode struct {
//	*BaseNode
//	mu            sync.Mutex
//	messageQueue  []interface{}
//	startTime     time.Time // Add startTime variable
//	totalMessages int
//}
//
//// NewAnyBufferNode creates a new AnyBufferNode with the given ID, name, EventBus, and Flow.
//func NewAnyBufferNode(nodeUUID, name string, bus *EventBus) *AnyBufferNode {
//	node := NewBaseNode("any-buffer", nodeUUID, name, bus)
//	node.NewInputPort("input", "input", portTypeAny)
//	node.NewOutputPort("output-1", "output-1", portTypeAny)
//	node.NewOutputPort("output-2", "output-2", portTypeFloat)
//	node.NewOutputPort("output-3", "output-3", portTypeAvgMessageCount)
//	node.AddToNodesMap(nodeUUID, node)
//	node.SetLoaded(true)
//	node.SetHotFix()
//
//	return &AnyBufferNode{
//		BaseNode:      node,
//		messageQueue:  make([]interface{}, 0),
//		totalMessages: 0,
//	}
//}
//
//func (n *AnyBufferNode) Start() {
//	for _, sub := range n.Connections {
//		go func(sub *Connection) {
//			inputChannel, exists := n.Bus[sub.TargetUUID]
//			if !exists {
//				fmt.Printf("Input channel for target input %s does not exist\n", sub.TargetUUID)
//				return
//			}
//
//			for {
//				select {
//				case msg, ok := <-inputChannel:
//					if !ok {
//						return
//					}
//
//					// Lock the mutex before accessing the buffer
//					n.mu.Lock()
//					n.messageQueue = append(n.messageQueue, msg.Port.Value)
//					n.totalMessages++
//
//					if n.totalMessages > 10 {
//						// Calculate the sum of valid numeric values
//						sum := calculateSum(n.messageQueue)
//
//						// Calculate the average message count
//						avgPerSec := float64(n.totalMessages) / time.Since(n.startTime).Seconds()
//						avgPerMin := avgPerSec * 60
//						avgPerHour := avgPerMin * 60
//						avgPerDay := avgPerHour * 24
//
//						// Send values on the outputs
//						sendValueOnOutput(n, "output-1", sum)
//						sendValueOnOutput(n, "output-2", float64(n.totalMessages))
//						sendValueOnOutput(n, "output-3", AvgMessageCount{
//							PerSec:  avgPerSec,
//							PerMin:  avgPerMin,
//							PerHour: avgPerHour,
//							PerDay:  avgPerDay,
//						})
//					}
//					// Unlock the mutex after updating the buffer
//					n.mu.Unlock()
//				}
//			}
//		}(sub)
//	}
//}
//
//// AvgMessageCount represents the average message counts per time interval.
//type AvgMessageCount struct {
//	PerSec  float64 `json:"perSec"`
//	PerMin  float64 `json:"perMin"`
//	PerHour float64 `json:"perHour"`
//	PerDay  float64 `json:"perDay"`
//}
//
//// CalculateAvgMessageCount calculates the average message counts based on the total messages and time duration.
//func CalculateAvgMessageCount(totalMessages int, duration time.Duration) AvgMessageCount {
//	perSec := float64(totalMessages) / duration.Seconds()
//	perMin := perSec * 60
//	perHour := perMin * 60
//	perDay := perHour * 24
//
//	return AvgMessageCount{
//		PerSec:  perSec,
//		PerMin:  perMin,
//		PerHour: perHour,
//		PerDay:  perDay,
//	}
//}
//
//func calculateSum(values []interface{}) float64 {
//	sum := 0.0
//	for _, value := range values {
//		// Check if the value is a valid numeric type (int or float)
//		switch value.(type) {
//		case int:
//			sum += float64(value.(int))
//		case float64:
//			sum += value.(float64)
//		}
//	}
//	return sum
//}
//
//func sendValueOnOutput(n *AnyBufferNode, outputID string, value interface{}) {
//	out := &Port{
//		ID:        outputID,
//		Name:      outputID,
//		Value:     value,
//		Direction: output,
//		DataType:  portTypeAny,
//	}
//	n.PublishMessage(out, true)
//}
