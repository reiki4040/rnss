package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetEC2ListForStartSession(ctx context.Context, region string) ([]string, error) {
	instances, err := GetInstances(ctx, region)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0, len(instances))
	for _, ins := range instances {
		// filter running instance
		if ins.State.Name != types.InstanceStateNameRunning {
			continue
		}

		n := "No Name Tag"
		for _, t := range ins.Tags {
			if t.Key != nil && *t.Key == "Name" {
				n = *t.Value
				break
			}
		}

		list = append(list, fmt.Sprintf("%s\t%s\t%s\t%s",
			convertNilString(ins.InstanceId),
			n,
			convertNilString(ins.PrivateIpAddress),
			convertNilString(ins.PublicIpAddress)))
	}

	return list, nil
}

func GetInstances(ctx context.Context, region string) ([]*types.Instance, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	cli := ec2.NewFromConfig(cfg)

	stateQueryName := "instance-state-name"
	filters := []types.Filter{
		{
			Name:   &stateQueryName,
			Values: []string{"running"},
		},
	}
	resp, err := cli.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Reservations) == 0 {
		return []*types.Instance{}, nil
	}

	instances := make([]*types.Instance, 0)
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			instances = append(instances, &i)
		}
	}

	return instances, nil
}

func convertNilString(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}
