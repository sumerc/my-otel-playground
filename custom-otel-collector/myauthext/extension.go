package myauthext

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/auth"
)

type myAuth struct {
}

func newMyAuthExtension(cfg *Config) (auth.Server, error) {
	ba := myAuth{}
	return auth.NewServer(
		auth.WithServerStart(ba.serverStart),
		auth.WithServerAuthenticate(ba.authenticate),
	), nil
}

func (ba *myAuth) serverStart(_ context.Context, _ component.Host) error {
	fmt.Println("serverStart")
	return nil
}

type myAuthData struct {
	attributes map[string]any
}

func NewMyAuthData() *myAuthData {
	return &myAuthData{
		attributes: make(map[string]any),
	}
}

func (ad *myAuthData) GetAttribute(name string) any {
	return ad.attributes[name]
}

func (ad *myAuthData) GetAttributeNames() []string {
	names := make([]string, 0, len(ad.attributes))
	for name := range ad.attributes {
		names = append(names, name)
	}
	return names
}

func (ad *myAuthData) SetAttribute(name string, value any) {
	ad.attributes[name] = value
}

func (ba *myAuth) authenticate(ctx context.Context, headers map[string][]string) (context.Context, error) {
	fmt.Println("authenticate4")
	fmt.Println(headers)
	// cl := client.FromContext(ctx)
	// // //cl.Auth = authData
	// ctx = client.NewContext(ctx, cl)

	// ctx = context.WithValue(ctx, "auth.tenant_id", "tenant-id-value1")
	// ctx = context.WithValue(ctx, "tenant_id", "tenant-id-value2")

	// method1 of sharing data with the other components in the pipeline
	// see: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/39a7773175cdea009125c59e16acc031eefe2d0c/internal/coreinternal/attraction/attraction.go#L345
	// function getAttributeValueFromContext
	cli := client.FromContext(ctx)
	mdd := map[string][]string{"tenant_id": {"tenant-id-value1"}}
	md := client.NewMetadata(mdd)
	cli.Metadata = md

	// method2: setting AuthData
	ad := NewMyAuthData()
	ad.SetAttribute("tenant_id", "tenant-id-value2")
	cli.Auth = ad

	// Set the tenantID in the client
	//fmt.Println("cli.metadata", cli.Metadata)

	// Put the updated client back into the context
	ctx = client.NewContext(ctx, cli)

	return ctx, nil
}
