package ecspresso

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var _ecs *ecs.ECS

type ListOption struct {
	ConfigFile *bool
}

func (d *App) ListServices(opt ListOption) error {
	config := d.config
	fmt.Println("+++ List of services for cluster: ", config.Cluster)
	svc, err := ListServices(config.Cluster)

	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println("+++")
	for _, v := range svc {
		fmt.Println(v)
	}

	return nil
}

func ListServices(cluster string) ([]string, error) {
	svc := assertECS()
	params := &ecs.ListServicesInput{Cluster: aws.String(cluster)}
	result := make([]*string, 0)
	err := svc.ListServicesPages(params, func(services *ecs.ListServicesOutput, lastPage bool) bool {
		result = append(result, services.ServiceArns...)
		return !lastPage
	})
	if err != nil {
		return nil, err
	}
	out := make([]string, len(result))
	for i, s := range result {
		out[i] = path.Base(*s)
	}
	return out, nil
}

func assertECS() *ecs.ECS {
	if _ecs == nil {
		_ecs = ecs.New(session.New(getServiceConfiguration()))
	}
	return _ecs
}

func getServiceConfiguration() *aws.Config {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	return &aws.Config{Region: aws.String(region)}
}
