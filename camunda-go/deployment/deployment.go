package main

import (
	"context"
	"fmt"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func main() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         "localhost:26500",
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	_, err = client.NewTopologyCommand().Send(ctx)
	if err != nil {
		panic(err)
	}

	//for _, broker := range topology.Brokers {
	//	fmt.Println("Broker", broker.Host, ":", broker.Port)
	//	for _, partition := range broker.Partitions {
	//		fmt.Println("  Partition", partition.PartitionId, ":", roleToString(partition.Role))
	//	}
	//}

	// After the client is created (add this to the end of your func main())
	response, err := client.NewDeployResourceCommand().AddResourceFile("order-process.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.String())

	// After the process is deployed.
	variables := make(map[string]interface{})
	variables["orderId"] = "31243"

	request, err := client.NewCreateInstanceCommand().BPMNProcessId("order-process").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}

	msg, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(msg.String())
}

func roleToString(role pb.Partition_PartitionBrokerRole) string {
	switch role {
	case pb.Partition_LEADER:
		return "Leader"
	case pb.Partition_FOLLOWER:
		return "Follower"
	default:
		return "Unknown"
	}
}
