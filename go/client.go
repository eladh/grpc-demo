package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"grpc-poc/insecure"
	"grpc-poc/registration"
)

//go:generate protoc -I ./registration --go_out=plugins=grpc:./registration ./registration/registration.proto

func NewClient(serverAddr string, encryption bool, compressed bool) (registration.RegistrationClient, *grpc.ClientConn, error) {
	ctx := context.Background()
	var dialOptions []grpc.DialOption
	if encryption {
		dialOptions = append(dialOptions,
			grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")))

		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(tokenAuth{token: "abc"}))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	if compressed {
		dialOptions = append(dialOptions,
			grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
		)
	}

	conn, err := grpc.DialContext(ctx, serverAddr, dialOptions...)

	if err != nil {
		return nil, nil, err
	}

	registrationClient := registration.NewRegistrationClient(conn)

	return registrationClient, conn, nil
}

type tokenAuth struct {
	token string
}

// Return value is mapped to request headers.
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}
