package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := rds.NewFromConfig(cfg)

	output, err := client.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})

	baseSizes := map[string]string{
		"t": "micro",
		"m": "large",
	}

	sizeToUnits := map[string]float64{
		"micro":    0.5,
		"small":    1,
		"large":    4,
		"xlarge":   8,
		"2xlarge":  16,
		"4xlarge":  32,
		"8xlarge":  64,
		"10xlarge": 80,
		"12xlarge": 96,
		"16xlarge": 128,
		"24xlarge": 192,
		"32xlarge": 256,
	}

	byFamily := map[string]float64{}

	for _, instance := range output.DBInstances {
		parts := strings.Split(aws.ToString(instance.DBInstanceClass), ".")
		family, size := parts[1], parts[2]

		multiplier := 1
		if instance.MultiAZ {
			multiplier = 2
		}

		byFamily[family] += sizeToUnits[size] * float64(multiplier)
	}

	for family, units := range byFamily {
		baseSize := baseSizes[family[0:1]]
		baseUnits := sizeToUnits[baseSize]
		fmt.Printf("%s: %d * %s.%s\n", family, int(units/baseUnits), family, baseSize)
	}
}
