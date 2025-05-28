package main

import (
	"log"
	"fmt"
	"net/url"
	"golang.org/x/net/context"

	"github.com/apple/pkl-go/pkl"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {
	client, err := pkl.NewExternalReaderClient(pkl.WithExternalClientResourceReader(awsSecretReader{}))
	if err != nil {
		log.Fatalln(err)
	}
	if err := client.Run(); err != nil {
		log.Fatalln(err)
	}
}

type awsSecretReader struct{}

var _ pkl.ResourceReader = &awsSecretReader{}

func (r awsSecretReader) Scheme() string {
	return "awssecret"
}

func (r awsSecretReader) HasHierarchicalUris() bool {
	return false
}

func (r awsSecretReader) IsGlobbable() bool {
	return false
}

func (r awsSecretReader) ListElements(baseURI url.URL) ([]pkl.PathElement, error) {
	return nil, nil // Not implemented as we're not globbable or heirarchical
}

func (r awsSecretReader) Read(uri url.URL) ([]byte, error) {
	if !arn.IsARN(uri.Opaque) {
		return nil, fmt.Errorf("invalid AWS ARN: %s", uri.Opaque)
	}
	arn, err := arn.Parse(uri.Opaque)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ARN: %w", err)
	}
	if arn.Service != "secretsmanager" {
		return nil, fmt.Errorf("ARN must be for the secretsmanager service: %s", uri.Opaque)
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := secretsmanager.NewFromConfig(cfg)
	secret, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &uri.Opaque,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret value: %w", err)
	}
	if secret.SecretString != nil {
		return []byte(*secret.SecretString), nil
	}
	if secret.SecretBinary != nil {
		return secret.SecretBinary, nil
	}
	return nil, fmt.Errorf("secret value is empty for ARN: %s", uri.Opaque)
}