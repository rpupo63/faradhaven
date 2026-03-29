package api

import (
	"context"
	"errors"
)

type keyType string

const (
	userIDKey         keyType = "userID"
	isAdminKey        keyType = "isAdmin"
	organizationIDKey keyType = "organizationID"
	userKey           keyType = "user"
)

// ctxWithUserID adds a user ID to the context
func ctxWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// ctxWithIsAdmin adds admin status to the context
func ctxWithIsAdmin(ctx context.Context, isAdmin bool) context.Context {
	return context.WithValue(ctx, isAdminKey, isAdmin)
}

// ctxWithOrganizationID adds an organization ID to the context
func ctxWithOrganizationID(ctx context.Context, organizationID string) context.Context {
	return context.WithValue(ctx, organizationIDKey, organizationID)
}

// ctxGetUserID retrieves a user ID from the context
func ctxGetUserID(ctx context.Context) (string, error) {
	return ctxGetStringValue(ctx, userIDKey)
}

// ctxGetIsAdmin retrieves admin status from the context
func ctxGetIsAdmin(ctx context.Context) (bool, error) {
	if ctxValue := ctx.Value(isAdminKey); ctxValue == nil {
		return false, errors.New("isAdmin not found in context")
	} else if valueAsBool, ok := ctxValue.(bool); !ok {
		return false, errors.New("isAdmin is not of type `bool`")
	} else {
		return valueAsBool, nil
	}
}

// ctxGetStringValue is a helper function to retrieve string values from the context by key
func ctxGetStringValue(ctx context.Context, key keyType) (string, error) {
	if ctxValue := ctx.Value(key); ctxValue == nil {
		return "", errors.New("key not found in context")
	} else if valueAsString, ok := ctxValue.(string); !ok {
		return "", errors.New("value is not of type `string`")
	} else {
		return valueAsString, nil
	}
}
