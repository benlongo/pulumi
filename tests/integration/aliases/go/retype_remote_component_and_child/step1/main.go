// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.
//go:build !all
// +build !all

package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Bucket struct {
	pulumi.ResourceState
}

type BucketComponent struct {
	pulumi.ResourceState
}

func NewBucket(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*Bucket, error) {
	bucketResource := &Bucket{}
	err := ctx.RegisterResource("wibble:index:Bucket", name, nil, bucketResource, opts...)
	if err != nil {
		return nil, err
	}
	return bucketResource, nil
}

func NewBucketComponent(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*BucketComponent, error) {
	bucketComponent := &BucketComponent{}
	err := ctx.RegisterRemoteComponentResource("wibble:aws/s3:Bucket", name, nil, bucketComponent, opts...)
	if err != nil {
		return nil, err
	}

	parentOpt := pulumi.Parent(bucketComponent)
	_, err = NewBucket(ctx, name+"-child", parentOpt)
	if err != nil {
		return nil, err
	}
	_, err = NewBucket(ctx, "otherchild", parentOpt)
	if err != nil {
		return nil, err
	}
	return bucketComponent, nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := NewBucketComponent(ctx, "main-bucket")
		if err != nil {
			return err
		}

		return nil
	})
}
