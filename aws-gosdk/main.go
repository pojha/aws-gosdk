package main

import (
	"fmt"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	
    m.Get("/", func() string {
		nt := getLatest()
        return nt;
    })
    m.Run()
	//m.RunOnAddr(":8080")
	
}

func getLatest() string{
	nt := "_________________USA________________________________ \n "
	nt += list_instance("us-east-1")
	nt += "________________UK________________________________ \n "
	nt += list_instance("eu-west-2")
	nt += "________________AUS________________________________ \n "
	nt += list_instance("ap-southeast-2")

	// nt += "________________S3 Buckets_______________________________ \n "
	// nt += list_bucket("us-east-1")
	// nt += list_bucket("eu-west-2")
	// nt += list_bucket("ap-southeast-2")
	fmt.Println(nt)
	return nt
}

func list_instance(region string) string {
	fmt.Println(region);

	ec2svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			// {
			// 	Name:   aws.String("tag:app_name"),
			// 	Values: []*string{aws.String(""), aws.String("")},
			// },
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
			//  {
			// 	"Name": "availability-zone",
			// 	"Values": ["us-east-1"]
			// }
		},
	}
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}


	var nt string	
	for idx, res := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {			
			if res != nil {
				for _, t := range inst.Tags {
					if *t.Key == "Name" {
						nt +=  *inst.InstanceType + "  "+*inst.State.Name+"     " + *inst.PrivateIpAddress +"    "+*t.Value 
						nt += "\n" 
				//		break
					}
				}

			}						
		}		
	}
	fmt.Println(nt)
	return nt
}


func list_bucket(region string) string{
		
	s3svc := s3.New(session.New(), &aws.Config{Region: aws.String(region)})
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Failed to list buckets", err)		
	}
	var nt string
	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		nt +=  aws.StringValue(bucket.Name)		
		nt += "\n" 
	}

	fmt.Println(nt)
	return nt
}
