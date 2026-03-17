package auth

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
)

type AzureAuth struct {
	credential *azidentity.InteractiveBrowserCredential
	tenantID   string
	clientID   string
}

func NewAzureAuth(tenantID, clientID string) (*AzureAuth, error) {
	c, err := cache.New(&cache.Options{
		Name: "th",
	})
	if err != nil {
		return nil, fmt.Errorf("creating token cache: %w", err)
	}

	options := &azidentity.InteractiveBrowserCredentialOptions{
		ClientID: clientID,
		TenantID: tenantID,
		Cache:    c,
	}

	cred, err := azidentity.NewInteractiveBrowserCredential(options)
	if err != nil {
		return nil, fmt.Errorf("creating browser credential: %w", err)
	}

	return &AzureAuth{
		credential: cred,
		tenantID:   tenantID,
		clientID:   clientID,
	}, nil
}

func (a *AzureAuth) GetToken(ctx context.Context, scope string) (string, error) {
	token, err := a.credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		if isAuthenticationRequired(err) {
			if err := a.clearCache(); err != nil {
				return "", fmt.Errorf("clearing token cache: %w", err)
			}
			token, err = a.credential.GetToken(ctx, policy.TokenRequestOptions{
				Scopes: []string{scope},
			})
			if err != nil {
				return "", fmt.Errorf("re-authenticating: %w", err)
			}
			return token.Token, nil
		}
		return "", fmt.Errorf("getting token: %w", err)
	}

	return token.Token, nil
}

func (a *AzureAuth) clearCache() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("getting user cache dir: %w", err)
	}
	cachePath := filepath.Join(cacheDir, "th", "token-cache.json")
	os.Remove(cachePath)
	return nil
}

func isAuthenticationRequired(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "AADSTS") ||
		strings.Contains(errStr, "authentication required") ||
		(strings.Contains(errStr, "token") && strings.Contains(errStr, "invalid"))
}

func GetDefaultScope(endpoint string) string {
	return endpoint + "/.default"
}
