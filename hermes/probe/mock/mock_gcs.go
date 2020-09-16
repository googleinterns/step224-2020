package mock

import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
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
	idx  int                                                                                                                                                 
	objects map[int]*storage.ObjectAttrs                                                                                                                        
	stiface.ObjectIterator                                                                                                                                      
}  

func NewFakeClient() stiface.Client {
	return &fakeClient{buckets: map[string]*fakeBucket{}}
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

func (b fakeBucketHandle) Object(name string) stiface.ObjectHandle {
	return fakeObjectHandle{c: b.c, bucketName: b.name, name: name}
}

func (c *fakeClient) Bucket(name string) stiface.BucketHandle {
	return fakeBucketHandle{c: c, name: name}
}

func (w *fakeWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
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
