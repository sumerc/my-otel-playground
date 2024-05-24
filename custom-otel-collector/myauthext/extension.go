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

func (ba *myAuth) authenticate(ctx context.Context, headers map[string][]string) (context.Context, error) {
	fmt.Println("authenticate")
	fmt.Println(headers)
	cl := client.FromContext(ctx)
	//cl.Auth = authData
	return client.NewContext(ctx, cl), nil
}
