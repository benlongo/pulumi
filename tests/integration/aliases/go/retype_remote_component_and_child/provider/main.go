// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.
//go:build !all
// +build !all

package main

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	pbempty "github.com/golang/protobuf/ptypes/empty"
)

type Bucket struct {
	pulumi.ResourceState
}

type BucketV2 struct {
	pulumi.ResourceState
}

type BucketArgs struct{}
type BucketV2Args struct{}

func NewBucket(ctx *pulumi.Context, name string, args *BucketArgs,
	opts ...pulumi.ResourceOption) (*Bucket, error) {

	component := &Bucket{}
	err := ctx.RegisterComponentResource(providerName+":aws/s3:Bucket", name, component, opts...)
	if err != nil {
		return nil, err
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{}); err != nil {
		return nil, err
	}

	return component, nil
}

const providerName = "wibble"
const version = "0.0.1"

var currentID int

func main() {
	err := provider.Main(providerName, func(host *provider.HostClient) (pulumirpc.ResourceProviderServer, error) {
		return makeProvider(host, providerName, version)
	})
	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}

type Provider struct {
	pulumirpc.UnimplementedResourceProviderServer

	host    *provider.HostClient
	name    string
	version string

	expectResourceArg bool
}

func makeProvider(host *provider.HostClient, name, version string) (pulumirpc.ResourceProviderServer, error) {
	return &Provider{
		host:    host,
		name:    name,
		version: version,
	}, nil
}

func (p *Provider) Create(ctx context.Context,
	req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	id := currentID
	currentID++

	return &pulumirpc.CreateResponse{
		Id: fmt.Sprintf("%v", id),
	}, nil
}

func (p *Provider) Construct(ctx context.Context,
	req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {
	return pulumiprovider.Construct(ctx, req, p.host.EngineConn(), func(ctx *pulumi.Context, typ, name string,
		inputs pulumiprovider.ConstructInputs, options pulumi.ResourceOption) (*pulumiprovider.ConstructResult, error) {

		allowedTypes := []string{
			providerName + ":aws/s3:Bucket",
			providerName + ":aws/s3:BucketV2",
			providerName + ":index:Bucket",
			providerName + ":index:BucketV2",
		}

		isAllowd := false
		for _, allowdType := range allowedTypes {
			if allowdType == typ {
				isAllowd = true
				break
			}
		}

		if !isAllowd {
			return nil, fmt.Errorf("unknown resource type %s", typ)
		}

		args := &BucketArgs{}
		if err := inputs.CopyTo(args); err != nil {
			return nil, fmt.Errorf("setting args: %w", err)
		}

		component, err := NewBucket(ctx, name, args, options)
		if err != nil {
			return nil, fmt.Errorf("creating component: %w", err)
		}

		return pulumiprovider.NewConstructResult(component)
	})
}

func (p *Provider) CheckConfig(ctx context.Context,
	req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

func (p *Provider) DiffConfig(ctx context.Context,
	req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return &pulumirpc.DiffResponse{}, nil
}

func (p *Provider) Configure(ctx context.Context,
	req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	if _, ok := req.GetArgs().Fields["expectResourceArg"]; ok {
		p.expectResourceArg = true
	}
	return &pulumirpc.ConfigureResponse{
		AcceptSecrets:   true,
		SupportsPreview: true,
		AcceptResources: true,
	}, nil
}

func (p *Provider) Invoke(ctx context.Context,
	req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	return nil, fmt.Errorf("Unknown Invoke token '%s'", req.GetTok())
}

func (p *Provider) StreamInvoke(req *pulumirpc.InvokeRequest,
	server pulumirpc.ResourceProvider_StreamInvokeServer) error {
	return fmt.Errorf("Unknown StreamInvoke token '%s'", req.GetTok())
}

func (p *Provider) Call(ctx context.Context,
	req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, fmt.Errorf("Unknown Call token '%s'", req.GetTok())
}

func (p *Provider) Check(ctx context.Context,
	req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

func (p *Provider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return &pulumirpc.DiffResponse{}, nil
}

func (p *Provider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	return &pulumirpc.ReadResponse{
		Id:         req.GetId(),
		Properties: req.GetProperties(),
	}, nil
}

func (p *Provider) Update(ctx context.Context,
	req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	return &pulumirpc.UpdateResponse{
		Properties: req.GetNews(),
	}, nil
}

func (p *Provider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}

func (p *Provider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: p.version,
	}, nil
}

func (p *Provider) Attach(ctx context.Context, req *pulumirpc.PluginAttach) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}

func (p *Provider) GetSchema(ctx context.Context,
	req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	return &pulumirpc.GetSchemaResponse{}, nil
}

func (p *Provider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}

func (p *Provider) GetMapping(context.Context, *pulumirpc.GetMappingRequest) (*pulumirpc.GetMappingResponse, error) {
	return &pulumirpc.GetMappingResponse{}, nil
}
