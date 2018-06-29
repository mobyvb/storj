// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package objects

import (
	"context"
	"fmt"
	"io"
	"path"
	"time"

	"storj.io/storj/pkg/ranger"
)

//Objects structure
type Objects struct {
	//segStore    segments.SegmentStore
	//streamStore streams.StreamStore
}

//Meta structure
type Meta struct {
	Modified   time.Time
	Expiration time.Time
	Data       []byte
}

// func NewObjects(store streams.StreamStore) ObjectStore {
// 	panic("TODO")
// }

//PutObject interface method
func (o *Objects) PutObject(ctx context.Context, objpath string, data io.Reader, metadata []byte, expiration time.Time) (err error) {
	defer mon.Task()(&ctx)(&err)
	panic("TODO")
}

//GetObject interface method
func (o *Objects) GetObject(ctx context.Context, objpath string) (r ranger.Ranger, m Meta, err error) {
	defer mon.Task()(&ctx)(&err)
	getObjectPath := path.Join("object", objpath)

	/* TODO@ASK clean up the below line */
	fmt.Println(getObjectPath)
	panic("TODO")
}

//DeleteObject interface method
func (o *Objects) DeleteObject(ctx context.Context, objpath string) (err error) {
	defer mon.Task()(&ctx)(&err)
	panic("TODO")
}

//ListObjects interface method
func (o *Objects) ListObjects(ctx context.Context, startingPath, endingPath string) (objpaths []string, truncated bool, err error) {
	defer mon.Task()(&ctx)(&err)
	panic("TODO")
}

//SetXAttr interface method
func (o *Objects) SetXAttr(ctx context.Context, objpath, xattr string, data io.Reader, metadata []byte) (err error) {
	defer mon.Task()(&ctx)(&err)
	panic("TODO")
}

//GetXAttr interface method
func (o *Objects) GetXAttr(ctx context.Context, objpath, xattr string) (r ranger.Ranger, m Meta, err error) {
	defer mon.Task()(&ctx)(&err)
	getXAttrPath := path.Join("xattr", objpath, xattr)

	/* TODO@ASK clean up the below line */
	fmt.Println(getXAttrPath)
	panic("TODO")
}

//DeleteXAttr interface method
func (o *Objects) DeleteXAttr(ctx context.Context, path, xattr string) (err error) {
	defer mon.Task()(&ctx)(&err)
	panic("TODO")
}

//ListXAttrs interface method
func (o *Objects) ListXAttrs(ctx context.Context, path, startingXAttr, endingXAttr string) (xattrs []string, truncated bool, err error) {
	defer mon.Task()(&ctx)(&err)
	panic("TODO")
}