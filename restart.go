package ecspresso

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type RestartOption struct {
	ConfigFile *bool
}

func (d *App) RestartService(opt RestartOption) error {
	config := d.config
	fmt.Println("+++ List of services for cluster: ", config.Cluster)
	servicename := os.Getenv("SERVICE_NAME")
	err := ForceNewDeployment(config.Cluster, servicename)

	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println("+++")
	fmt.Println("+++ Restarting ", servicename)

	return nil
}

func ForceNewDeployment(clusterName string, serviceName string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	ecsSvc := ecs.New(sess, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	serviceParams := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(serviceName),
		},
		Cluster: aws.String(clusterName),
	}
	result, err := ecsSvc.DescribeServices(serviceParams)
	if err != nil {
		return err
	}
	if len(result.Services) == 0 {
		return fmt.Errorf("Could not find service %s in cluster %s", serviceName, clusterName)
	}

	// Update Service
	for i := range result.Services {
		service := result.Services[i]
		if *service.ServiceName == serviceName {
			newServiceParams := &ecs.UpdateServiceInput{
				Service:            service.ServiceName,
				Cluster:            aws.String(clusterName),
				ForceNewDeployment: aws.Bool(true),
			}
			_, err := ecsSvc.UpdateService(newServiceParams)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
