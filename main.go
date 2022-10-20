package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// TODO: Update exit codes to match BSD suggested standard
// https://www.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html

func main() {
	if len(os.Args) <= 3 {
		fmt.Printf("Must provide key and value of metric to record")
		return
	}

	// TODO: This should be constrained against AWS types;
	awsRegionPtr := flag.String("region", "us-west-1", "AWS Region")
	metricNameSpacePtr := flag.String("namespace", "wicket-metrics", "Cloudwatch Metric Namespace")

	name := os.Args[2]
	rawValue := os.Args[3]

	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		os.Stderr.WriteString("Must provide valid numeric value")
		return
	}

	// TODO: Pusb AWS code behind to a storage interface
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		os.Stderr.WriteString("Failed to load config for AWS")
		panic(err)
	}

	client := cloudwatch.NewFromConfig(cfg, func(o *cloudwatch.Options) {
		o.Region = *awsRegionPtr
	})

	datum := types.MetricDatum{
		MetricName: &name,
		Value:      &value,
	}

	data := make([]types.MetricDatum, 1)
	data[0] = datum

	metricData, err := client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		MetricData: data,
		Namespace:  metricNameSpacePtr,
	})

	// A go-idiom for getting around unnused variables?
	_ = metricData

	if err != nil {
		os.Stderr.WriteString("Failed to store metric data")
		return
	}
}
