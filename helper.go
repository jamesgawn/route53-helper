package route53helper

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func LoadConfig() (aws.Config, error) {
	awsConfig, awsConfigError := config.LoadDefaultConfig(context.TODO())
	return awsConfig, awsConfigError
}

func FindZone(ctx context.Context, client *route53.Client, zoneName *string) (*types.HostedZone, error) {

	listHostedZonesInput := &route53.ListHostedZonesInput{}
	zoneList, zoneListError := client.ListHostedZones(ctx, listHostedZonesInput)

	if zoneListError != nil {
		return nil, zoneListError
	}

	for _, zone := range zoneList.HostedZones {
		if *zone.Name == *zoneName {
			return &zone, nil
		}
	}

	return &types.HostedZone{}, fmt.Errorf("unable to find zone: %s", *zoneName)
}

func UpdateRecord(ctx context.Context, client *route53.Client, zone *types.HostedZone, domain *string, ip *string) error {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeAction("UPSERT"),
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: domain,
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: ip,
							},
						},
						TTL: aws.Int64(300),
					},
				},
			},
			Comment: aws.String("Automated update from route53 helper"),
		},
		HostedZoneId: zone.Id,
	}
	_, err := client.ChangeResourceRecordSets(ctx, params)
	return err
}
