package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/iokuru/sentry/internal/model"
)

type Store struct {
	client         *dynamodb.Client
	incidentsTable string
	linksTable     string
}

func NewStore(client *dynamodb.Client, incidentsTable, linksTable string) *Store {
	return &Store{
		client:         client,
		incidentsTable: incidentsTable,
		linksTable:     linksTable,
	}
}

func (s *Store) SaveIncident(ctx context.Context, inc *model.Incident) error {
	now := time.Now().UTC()
	if inc.CreatedAt == "" {
		inc.CreatedAt = now.Format(time.RFC3339)
	}
	inc.UpdatedAt = now.Format(time.RFC3339)
	// ttl expires 30 days after creation
	if inc.TTL == 0 {
		inc.TTL = now.Add(30 * 24 * time.Hour).Unix()
	}

	item, err := attributevalue.MarshalMap(inc)
	if err != nil {
		return fmt.Errorf("marshal incident: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.incidentsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put incident %s: %w", inc.ID, err)
	}

	return nil
}

func (s *Store) GetIncident(ctx context.Context, id string) (*model.Incident, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.incidentsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get incident %s: %w", id, err)
	}
	if out.Item == nil {
		return nil, nil
	}

	var inc model.Incident
	if err := attributevalue.UnmarshalMap(out.Item, &inc); err != nil {
		return nil, fmt.Errorf("unmarshal incident %s: %w", id, err)
	}

	return &inc, nil
}

func (s *Store) ListIncidents(ctx context.Context, limit int32) ([]model.Incident, error) {
	out, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.incidentsTable),
		Limit:     aws.Int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("scan incidents: %w", err)
	}

	var incidents []model.Incident
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &incidents); err != nil {
		return nil, fmt.Errorf("unmarshal incidents: %w", err)
	}

	// sort newest first
	sort.Slice(incidents, func(i, j int) bool {
		return incidents[i].CreatedAt > incidents[j].CreatedAt
	})

	return incidents, nil
}

func (s *Store) UpdateIncidentStatus(ctx context.Context, id string, status model.IncidentStatus) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.incidentsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET #st = :status, #ua = :updated_at"),
		ExpressionAttributeNames: map[string]string{
			"#st": "status",
			"#ua": "updated_at",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":     &types.AttributeValueMemberS{Value: string(status)},
			":updated_at": &types.AttributeValueMemberS{Value: now},
		},
	})
	if err != nil {
		return fmt.Errorf("update incident status %s: %w", id, err)
	}
	return nil
}

func (s *Store) ListUnresolvedIncidents(ctx context.Context) ([]model.Incident, error) {
	incidents, err := s.ListIncidents(ctx, 100)
	if err != nil {
		return nil, err
	}
	var unresolved []model.Incident
	for _, inc := range incidents {
		if inc.Status != model.StatusResolved {
			unresolved = append(unresolved, inc)
		}
	}
	return unresolved, nil
}
