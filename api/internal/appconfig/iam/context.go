package iam

import (
	"context"
)

var UserProfileKey = "host_profile"

// HostProfile represents host profile information
type HostProfile struct {
	ID    int64
	Email string
}

// SetUserProfileInContext sets the HostProfile in the given context
func SetUserProfileInContext(ctx context.Context, p HostProfile) context.Context {
	return context.WithValue(ctx, UserProfileKey, p)
}

// UserProfileFromContext gets the HostProfile from the given context
func UserProfileFromContext(ctx context.Context) HostProfile {
	if v, ok := ctx.Value(UserProfileKey).(HostProfile); ok {
		return v
	}
	return HostProfile{}
}
