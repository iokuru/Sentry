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

var (
	ErrInviteNotFound = fmt.Errorf("invite code not found")
	ErrInviteExpired  = fmt.Errorf("invite code expired")
	ErrInviteUsed     = fmt.Errorf("invite code already used")
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
	// ttl expires 90 days after creation
	if inc.TTL == 0 {
		inc.TTL = now.Add(90 * 24 * time.Hour).Unix()
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

func (s *Store) CreateInvite(ctx context.Context, invite *model.Invite) error {
	now := time.Now().UTC().Unix()
	if invite.CreatedAt == 0 {
		invite.CreatedAt = now
	}
	if invite.ExpiresAt == 0 {
		invite.ExpiresAt = now + 86400
	}
	invite.TTL = invite.ExpiresAt
	invite.Used = false

	item := map[string]types.AttributeValue{
		"caregiver_id":      &types.AttributeValueMemberS{Value: "INVITE#" + invite.Code},
		"senior_id":         &types.AttributeValueMemberS{Value: "METADATA"},
		"code":              &types.AttributeValueMemberS{Value: invite.Code},
		"caregiver_id_real": &types.AttributeValueMemberS{Value: invite.CaregiverID},
		"senior_name":       &types.AttributeValueMemberS{Value: invite.SeniorName},
		"created_at":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", invite.CreatedAt)},
		"expires_at":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", invite.ExpiresAt)},
		"ttl":               &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", invite.TTL)},
		"used":              &types.AttributeValueMemberBOOL{Value: false},
	}

	_, err := s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.linksTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put invite %s: %w", invite.Code, err)
	}
	return nil
}

func (s *Store) GetInvite(ctx context.Context, code string) (*model.Invite, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.linksTable),
		Key: map[string]types.AttributeValue{
			"caregiver_id": &types.AttributeValueMemberS{Value: "INVITE#" + code},
			"senior_id":    &types.AttributeValueMemberS{Value: "METADATA"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get invite %s: %w", code, err)
	}
	if out.Item == nil {
		return nil, nil
	}
	var invite model.Invite
	if err := attributevalue.UnmarshalMap(out.Item, &invite); err != nil {
		return nil, fmt.Errorf("unmarshal invite %s: %w", code, err)
	}
	return &invite, nil
}

func (s *Store) RedeemInvite(ctx context.Context, code, seniorID string) (*model.CaregiverLink, error) {
	inv, err := s.GetInvite(ctx, code)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, ErrInviteNotFound
	}
	if inv.Used {
		return nil, ErrInviteUsed
	}
	if time.Now().UTC().Unix() > inv.ExpiresAt {
		return nil, ErrInviteExpired
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.linksTable),
		Key: map[string]types.AttributeValue{
			"caregiver_id": &types.AttributeValueMemberS{Value: "INVITE#" + code},
			"senior_id":    &types.AttributeValueMemberS{Value: "METADATA"},
		},
		UpdateExpression:    aws.String("SET #u = :true"),
		ConditionExpression: aws.String("#u = :false"),
		ExpressionAttributeNames: map[string]string{
			"#u": "used",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":true":  &types.AttributeValueMemberBOOL{Value: true},
			":false": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("consume invite: %w", err)
	}

	seniorName := inv.SeniorName
	if seniorName == "" {
		seniorName = "Senior"
	}
	link := &model.CaregiverLink{
		CaregiverID: inv.CaregiverID,
		SeniorID:    seniorID,
		SeniorName:  seniorName,
		LinkedAt:    time.Now().UTC(),
	}

	if err := s.CreateLink(ctx, link); err != nil {
		return nil, fmt.Errorf("create link on redeem: %w", err)
	}

	return link, nil
}

func (s *Store) CreateLink(ctx context.Context, link *model.CaregiverLink) error {
	if link.LinkedAt.IsZero() {
		link.LinkedAt = time.Now().UTC()
	}
	item, err := attributevalue.MarshalMap(link)
	if err != nil {
		return fmt.Errorf("marshal link: %w", err)
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.linksTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put link: %w", err)
	}
	return nil
}

func (s *Store) DeleteLink(ctx context.Context, caregiverID, seniorID string) error {
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.linksTable),
		Key: map[string]types.AttributeValue{
			"caregiver_id": &types.AttributeValueMemberS{Value: caregiverID},
			"senior_id":    &types.AttributeValueMemberS{Value: seniorID},
		},
	})
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}
	return nil
}

func (s *Store) GetLink(ctx context.Context, caregiverID, seniorID string) (*model.CaregiverLink, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.linksTable),
		Key: map[string]types.AttributeValue{
			"caregiver_id": &types.AttributeValueMemberS{Value: caregiverID},
			"senior_id":    &types.AttributeValueMemberS{Value: seniorID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get link: %w", err)
	}
	if out.Item == nil {
		return nil, nil
	}
	var link model.CaregiverLink
	if err := attributevalue.UnmarshalMap(out.Item, &link); err != nil {
		return nil, fmt.Errorf("unmarshal link: %w", err)
	}
	return &link, nil
}

func (s *Store) ListLinksForCaregiver(ctx context.Context, caregiverID string) ([]model.CaregiverLink, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.linksTable),
		KeyConditionExpression: aws.String("caregiver_id = :cid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cid": &types.AttributeValueMemberS{Value: caregiverID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query links for caregiver: %w", err)
	}
	var links []model.CaregiverLink
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &links); err != nil {
		return nil, fmt.Errorf("unmarshal links: %w", err)
	}
	return links, nil
}

func (s *Store) ListLinksForSenior(ctx context.Context, seniorID string) ([]model.CaregiverLink, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.linksTable),
		IndexName:              aws.String("senior_id-index"),
		KeyConditionExpression: aws.String("senior_id = :sid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sid": &types.AttributeValueMemberS{Value: seniorID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query links for senior: %w", err)
	}
	var links []model.CaregiverLink
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &links); err != nil {
		return nil, fmt.Errorf("unmarshal links: %w", err)
	}
	return links, nil
}

func (s *Store) SetTelegramChatID(ctx context.Context, caregiverID, chatID string) error {
	links, err := s.ListLinksForCaregiver(ctx, caregiverID)
	if err != nil {
		return err
	}
	for _, l := range links {
		_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(s.linksTable),
			Key: map[string]types.AttributeValue{
				"caregiver_id": &types.AttributeValueMemberS{Value: l.CaregiverID},
				"senior_id":    &types.AttributeValueMemberS{Value: l.SeniorID},
			},
			UpdateExpression: aws.String("SET telegram_chat_id = :cid"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":cid": &types.AttributeValueMemberS{Value: chatID},
			},
		})
		if err != nil {
			return fmt.Errorf("update telegram chat id: %w", err)
		}
	}
	return nil
}
