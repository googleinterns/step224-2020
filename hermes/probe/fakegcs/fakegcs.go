// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Authors: Alicja Kwiecinska, Evan Spendlove, GitHub: alicjakwie, espendlove
//
// package fakegcs contains all of the logic necessary to create a fake instance of GCS, supporting operations on buckets and object.
// TODO (#67) add doc strings
package fakegcs

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"google.golang.org/api/iterator"
)

type fakeBucket struct {
	attrs   *storage.BucketAttrs
	objects map[string][]byte
}

type fakeClient struct {
	stiface.Client
	buckets map[string]*fakeBucket
}

type fakeBucketHandle struct {
	stiface.BucketHandle
	c    *fakeClient
	name string
}

type fakeObjectHandle struct {
	stiface.ObjectHandle
	c          *fakeClient
	bucketName string
	name       string
}

type fakeWriter struct {
	stiface.Writer
	obj fakeObjectHandle
	buf bytes.Buffer
}

type fakeObjectIterator struct {
	idx     int
	objects map[int]*storage.ObjectAttrs
	stiface.ObjectIterator
}

type fakeReader struct {
	stiface.Reader
	r *bytes.Reader
}

func NewClient() *fakeClient {
	return &fakeClient{
		buckets: make(map[string]*fakeBucket),
	}
}

func (o fakeObjectHandle) NewWriter(context.Context) stiface.Writer {
	return &fakeWriter{obj: o}
}

func (b fakeBucketHandle) Create(_ context.Context, _ string, attrs *storage.BucketAttrs) error {
	if _, ok := b.c.buckets[b.name]; ok {
		return fmt.Errorf("bucket %q already exists", b.name)
	}
	if attrs == nil {
		attrs = &storage.BucketAttrs{}
	}
	attrs.Name = b.name
	b.c.buckets[b.name] = &fakeBucket{attrs: attrs, objects: map[string][]byte{}}
	return nil
}

func (o fakeObjectHandle) Delete(context.Context) error {
	bucket, ok := o.c.buckets[o.bucketName]
	if !ok {
		return fmt.Errorf("fakeObjectHandle.Delete(): bucket %q not found", o.bucketName)
	}
	delete(bucket.objects, o.name)
	return nil
}

func (b fakeBucketHandle) Object(name string) stiface.ObjectHandle {
	return fakeObjectHandle{c: b.c, bucketName: b.name, name: name}
}

func (c *fakeClient) Bucket(name string) stiface.BucketHandle {
	return fakeBucketHandle{c: c, name: name}
}

func (b fakeBucketHandle) Objects(_ context.Context, query *storage.Query) stiface.ObjectIterator {
	iterator := &fakeObjectIterator{
		idx:     0,
		objects: make(map[int]*storage.ObjectAttrs),
	}

	bkt, ok := b.c.buckets[b.name]
	if !ok {
		return nil
	}

	for obj := range bkt.objects {
		if strings.HasPrefix(obj, query.Prefix) {
			iterator.objects[iterator.idx] = &storage.ObjectAttrs{Name: obj}
		}
	}

	return iterator
}

func (i *fakeObjectIterator) Next() (*storage.ObjectAttrs, error) {
	if i.idx >= len(i.objects) {
		return nil, iterator.Done
	}
	obj := i.objects[i.idx]
	i.idx++
	return obj, nil
}

func (i *fakeObjectIterator) PageInfo() *iterator.PageInfo {
	return nil
}

func (w *fakeWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func (r fakeReader) Read(buf []byte) (int, error) {
	return r.r.Read(buf)
}

func (b fakeBucketHandle) Attrs(context.Context) (*storage.BucketAttrs, error) {
	bkt, ok := b.c.buckets[b.name]
	if !ok {
		return nil, fmt.Errorf("bucket %q does not exist", b.name)
	}
	return bkt.attrs, nil
}

func (w *fakeWriter) Close() error {
	bkt, ok := w.obj.c.buckets[w.obj.bucketName]
	if !ok {
		return fmt.Errorf("bucket %q not found", w.obj.bucketName)
	}
	bkt.objects[w.obj.name] = w.buf.Bytes()
	return nil
}

func (r fakeReader) Close() error {
	return nil
}

func (o fakeObjectHandle) NewReader(context.Context) (stiface.Reader, error) {
	bkt, ok := o.c.buckets[o.bucketName]
	if !ok {
		return nil, fmt.Errorf("bucket %q not found", o.bucketName)
	}
	contents, ok := bkt.objects[o.name]
	if !ok {
		return nil, fmt.Errorf("object %q not found in bucket %q", o.name, o.bucketName)
	}
	return fakeReader{r: bytes.NewReader(contents)}, nil
}
