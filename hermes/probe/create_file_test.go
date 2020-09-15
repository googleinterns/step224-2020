package probe
 
import (
   "fmt"
   "testing"
   "bytes"
   "context"
 
   "cloud.google.com/go/storage"
   probepb "github.com/googleinterns/step224-2020/config/proto"
   "github.com/googleapis/google-cloud-go-testing/storage/stiface"
   journalpb "github.com/googleinterns/step224-2020/hermes/proto"
   "github.com/googleinterns/step224-2020/hermes/probe/metrics"
)
 
func TestFileName(t *testing.T) {
   testUnits := [6]RandomHermesFile{
       {51, 12}, // Invalid ID 
       {0, 50}, // Invalid ID
       {3, 100}, // Expected Hermes_03_checksum
       {12, 100}, // Expected Hermes_12_checksum
       {3, 0}, // Invalid Size
       {3, 1001}, // Invalid Size
   }
   // Case: testUnits[0] RandomHermesFile{51, 12}
   want, wantErr := "", "The file ID provided 51 wasn't in the required range [1,50]"
   got, err := testUnits[0].FileName()
   gotErr:=err.Error() // error in string form
   if want!=got {
       t.Errorf("{51, 12}.FileName() failed expected %v got %v", want, got)
   }
   if wantErr!=gotErr {
       t.Errorf("{51, 12}.FileName() gave wrong error expected %v got %v", wantErr, gotErr)
   }
   // Case: testUnits[1] RandomHermesFile{0, 50}
   want, wantErr = "", fmt.Sprintf("The file ID provided %v wasn't in the required range [1,50]", testUnits[1].ID)
   got, err = testUnits[1].FileName()
   gotErr=err.Error() // error in string form
   if want!=got {
       t.Errorf("{0, 50}.FileName() failed expected %v got %v", want, got)
   }
   if wantErr!=gotErr {
       t.Errorf("{0, 50}.FileName() gave wrong error expected %v got %v", wantErr, gotErr)
   }
   // Case: testUnits[2] RandomHermesFile{3, 100}
   wantPrefix := "Hermes_03_"
   got, err = testUnits[2].FileName()
   if err!=nil {
       t.Errorf("Unexpected error: %v", err)
   }
   gotPrefix := got[0:10]
   if gotPrefix != wantPrefix {
       t.Errorf("{3,100}.FileName() failed expected prefix %v got prefix %v", wantPrefix, gotPrefix)
   }
   // Case: testUnits[3] RandomHermesFile{12, 100}
   wantPrefix = "Hermes_12_"
   got, err = testUnits[3].FileName()
   if err!=nil {
       t.Errorf("Unexpected error: %v", err)
   }
   gotPrefix = got[0:10]
   if gotPrefix != wantPrefix {
       t.Errorf("{12,100}.FileName() failed expected prefix %v got prefix %v", wantPrefix, gotPrefix)
   }
   // Case: testUnits[4] RandomHermesFile{3, 0}
   want, wantErr = "", "The file size provided 0 is not a positive number as required"
   got, err = testUnits[4].FileName()
   gotErr=err.Error() // error in string form
   if want!=got {
       t.Errorf("{3, 0}.FileName() failed expected %v got %v", want, got)
   }
   if wantErr!=gotErr {
       t.Errorf("{3, 0}.FileName() gave wrong error expected %v got %v", wantErr, gotErr)
   }
   // Case: testUnits[5] RandomHermesFile{3, 1001}
   want, wantErr = "", "The file size provided 1001 bytes exceeded the limit 1000 bytes"
   got, err = testUnits[5].FileName()
   gotErr=err.Error() // error in string form
   if want!=got {
       t.Errorf("{3, 1001}.FileName() failed expected %v got %v", want, got)
   }
   if wantErr!=gotErr {
       t.Errorf("{3, 1001}.FileName() gave wrong error expected %v got %v", wantErr, gotErr)
   }
}
 
func TestChecksum (t *testing.T) {
   f := RandomHermesFile{11, 100}
   fTwo := RandomHermesFile{13, 1000}
   fCopy := RandomHermesFile{11, 100}
   checksum, err := f.CheckSum()
   if(err!=nil){
       t.Error(err)
   }
   checksumTwo, err := fTwo.CheckSum()
   if(err!=nil){
       t.Error(err)
   }
   checksumCopy, err := fCopy.CheckSum()
   if(err!=nil){
       t.Error(err)
   }
   if fmt.Sprintf("%x",checksum) != fmt.Sprintf("%x",checksumCopy) { //comparing checksums converted to strings as a slice can only be compared to nil
       t.Errorf("Checksum returned different values for the same RandomHermesFiles {%v, %v}", f.ID, f.Size)
   }
   if fmt.Sprintf("%x",checksum) == fmt.Sprintf("%x",checksumTwo) {
       t.Errorf("Checksum returned the same value for two different RandomHermesFiles {%v, %v} and {%v, %v}", f.ID, f.Size, fTwo.ID, fTwo.Size)
   }
 
}
 
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
 
func newFakeClient() stiface.Client {
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
 
func TestCreateFile(t *testing.T){
   ctx := context.Background()
   bucketName := "test_bucket_probe0"
   client := newFakeClient()
   fbh := client.Bucket(bucketName) // fakeBucketHandle
   err := fbh.Create(ctx, bucketName, nil) // creates the bucket with name "test_bucket_probe0"
   var fileID int32 = 6
   fileSize := 50
   target := &Target {
       &probepb.Target{
           Name:                   "hermes",
           TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
           TotalSpaceAllocatedMib: int64(1000),
           BucketName:             "test_bucket_probe0",
       },
       &journalpb.StateJournal{
		Filenames :  make(map[int32]string),
	   },
       &metrics.Metrics{},
   }
   err = CreateFile(ctx, target, fileID, fileSize, &client, nil)
   if(err!=nil){
       t.Error(err)
   }
}