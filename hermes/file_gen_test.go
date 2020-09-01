package hermes

import(
	"testing"
	"strconv"
)


func TestChecksum(t *testing.T){
	global_contents := "jhfvjhdfjhfjjhjhdfvjvcvfjh"
	want := "_68f3caf439065824dcf75651c202e9f7c28ebf07" //expected checksum result
	file := HermesFile{}
	file.contents = global_contents;
	got := file.generateFileChecksum()
	if want != got {
		t.Errorf("generateFileChecksum() failed expected %v got %v", want, got)
	}
}

func TestFileName(t *testing.T){
	file := HermesFile{}
	file_ID := 23;
	fake_checksum :="_abba"
	want :=  "Hermes_"+strconv.Itoa(file_ID)+"_abba" //expected file name result
	file.generateFileName(file_ID, fake_checksum)
	got := file.name
	if got != want {
		t.Errorf("generateFileName(%v, \"abba\") failed expected %v got %v", file_ID, want, got)
	}
	file_ID = 4;
	want = "Hermes_0"+strconv.Itoa(file_ID)+"_abba" //expected file name result
	file.generateFileName(file_ID, fake_checksum)
	got = file.name
	if got != want {
		t.Errorf("generateFileName(%v, \"abba\") failed expected %v got %v", file_ID, want, got)
	}
}
