package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceOAuthIDPWriteModel struct {
	OAuthIDPWriteModel
}

func NewOAuthInstanceIDPWriteModel(instanceID, id string) *InstanceOAuthIDPWriteModel {
	return &InstanceOAuthIDPWriteModel{
		OAuthIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceOAuthIDPWriteModel) Reduce() error {
	return wm.OAuthIDPWriteModel.Reduce()
}

func (wm *InstanceOAuthIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OAuthIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *instance.OAuthIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPChangedEvent)
		}
	}
}

func (wm *InstanceOAuthIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OAuthIDPAddedEventType,
			instance.OAuthIDPChangedEventType,
		).
		Builder()
}

func (wm *InstanceOAuthIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options idp.Options,
) (*instance.OAuthIDPChangedEvent, error) {

	changes, err := wm.OAuthIDPWriteModel.NewChanges(
		name,
		clientID,
		clientSecretString,
		secretCrypto,
		authorizationEndpoint,
		tokenEndpoint,
		userEndpoint,
		scopes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewOAuthIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type InstanceOIDCIDPWriteModel struct {
	OIDCIDPWriteModel
}

func NewOIDCInstanceIDPWriteModel(instanceID, id string) *InstanceOIDCIDPWriteModel {
	return &InstanceOIDCIDPWriteModel{
		OIDCIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceOIDCIDPWriteModel) Reduce() error {
	return wm.OIDCIDPWriteModel.Reduce()
}

func (wm *InstanceOIDCIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OIDCIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *instance.OIDCIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.RemovedEvent)
		case *instance.IDPConfigAddedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPOIDCConfigAddedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *instance.IDPOIDCConfigChangedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.OIDCIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceOIDCIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OIDCIDPAddedEventType,
			instance.OIDCIDPChangedEventType,
			instance.IDPRemovedEventType,
			instance.IDPConfigAddedEventType,
			instance.IDPOIDCConfigAddedEventType,
			instance.IDPOIDCConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		Builder()
}

func (wm *InstanceOIDCIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	issuer,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.OIDCIDPChangedEvent, error) {

	changes, err := wm.OIDCIDPWriteModel.NewChanges(
		name,
		issuer,
		clientID,
		clientSecretString,
		secretCrypto,
		scopes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewOIDCIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type InstanceJWTIDPWriteModel struct {
	JWTIDPWriteModel
}

func NewJWTInstanceIDPWriteModel(instanceID, id string) *InstanceJWTIDPWriteModel {
	return &InstanceJWTIDPWriteModel{
		JWTIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceJWTIDPWriteModel) Reduce() error {
	return wm.JWTIDPWriteModel.Reduce()
}

func (wm *InstanceJWTIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.JWTIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *instance.JWTIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.RemovedEvent)
		case *instance.IDPConfigAddedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPJWTConfigAddedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *instance.IDPJWTConfigChangedEvent:
			if wm.ID != e.IDPConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			if wm.ID != e.ConfigID {
				continue
			}
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceJWTIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.JWTIDPAddedEventType,
			instance.JWTIDPChangedEventType,
			instance.IDPRemovedEventType,
			instance.IDPConfigAddedEventType,
			instance.IDPJWTConfigAddedEventType,
			instance.IDPJWTConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		Builder()
}

func (wm *InstanceJWTIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) (*instance.JWTIDPChangedEvent, error) {

	changes, err := wm.JWTIDPWriteModel.NewChanges(
		name,
		issuer,
		jwtEndpoint,
		keysEndpoint,
		headerName,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewJWTIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type InstanceGoogleIDPWriteModel struct {
	GoogleIDPWriteModel
}

func NewGoogleInstanceIDPWriteModel(instanceID, id string) *InstanceGoogleIDPWriteModel {
	return &InstanceGoogleIDPWriteModel{
		GoogleIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceGoogleIDPWriteModel) Reduce() error {
	return wm.GoogleIDPWriteModel.Reduce()
}

func (wm *InstanceGoogleIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GoogleIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.GoogleIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		}
	}
}

func (wm *InstanceGoogleIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GoogleIDPAddedEventType,
			instance.GoogleIDPChangedEventType,
		).
		Builder()
}

func (wm *InstanceGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewGoogleIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}
