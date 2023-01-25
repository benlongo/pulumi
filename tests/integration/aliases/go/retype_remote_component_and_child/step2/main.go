// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.
//go:build !all
// +build !all

package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type BucketV2 struct {
	pulumi.ResourceState
}

type BucketComponentV2 struct {
	pulumi.ResourceState
}

func NewBucketV2(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*BucketV2, error) {
	bucketResource := &BucketV2{}
	aliases := pulumi.Aliases([]pulumi.Alias{
		{
			Type: pulumi.String("wibble:index:Bucket"),
		},
	})
	opts = append(opts, aliases)
	err := ctx.RegisterResource("wibble:index:BucketV2", name, nil, bucketResource, opts...)
	if err != nil {
		return nil, err
	}
	return bucketResource, nil
}

func NewBucketComponentV2(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*BucketComponentV2, error) {
	bucketComponent := &BucketComponentV2{}
	aliases := pulumi.Aliases([]pulumi.Alias{
		{
			Type: pulumi.String("wibble:aws/s3:Bucket"),
		},
	})

	opts = append(opts, aliases)
	err := ctx.RegisterRemoteComponentResource("wibble:aws/s3:BucketV2", name, nil, bucketComponent, opts...)
	if err != nil {
		return nil, err
	}

	parentOpt := pulumi.Parent(bucketComponent)
	_, err = NewBucketV2(ctx, name+"-child", parentOpt)
	if err != nil {
		return nil, err
	}
	_, err = NewBucketV2(ctx, "otherchild", parentOpt)
	if err != nil {
		return nil, err
	}
	return bucketComponent, nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := NewBucketComponentV2(ctx, "main-bucket")
		if err != nil {
			return err
		}

		return nil
	})
}
