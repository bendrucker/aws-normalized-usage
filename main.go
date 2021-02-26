package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func main() {
	svc := flag.String("service", "", "AWS service")
	flag.Parse()

	switch *svc {
	case "rds":
		rdsUsage()
	case "ec2":
		ec2Usage()
	default:
		panic("unknown service")
	}
}

func ec2Usage() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := ec2.NewFromConfig(cfg)

	output, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		panic(err)
	}

	baseSizes := map[string]string{
		"t": "nano",
		"m": "large",
	}

	sizeToUnits := map[string]float64{
		"nano":     0.25,
		"micro":    0.5,
		"small":    1,
		"medium":   2,
		"large":    4,
		"xlarge":   8,
		"2xlarge":  12,
		"4xlarge":  24,
		"8xlarge":  48,
		"12xlarge": 72,
		"16xlarge": 96,
		"24xlarge": 144,
		"32xlarge": 192,
	}

	byFamily := map[string]float64{}

	for _, res := range output.Reservations {
		for _, instance := range res.Instances {
			parts := strings.Split(string(instance.InstanceType), ".")
			family, size := parts[0], parts[1]

			byFamily[family] += sizeToUnits[size]
		}
	}

	for family, units := range byFamily {
		baseSize, ok := baseSizes[family[0:1]]
		if !ok {
			continue
		}

		baseUnits := sizeToUnits[baseSize]
		fmt.Printf("%s: %d * %s.%s\n", family, int(units/baseUnits), family, baseSize)
	}
}

func rdsUsage() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := rds.NewFromConfig(cfg)

	output, err := client.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})

	if err != nil {
		panic(err)
	}

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
