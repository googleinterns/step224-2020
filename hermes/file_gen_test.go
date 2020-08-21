package hermes

import(
	//"fmt"
	"testing"
	"math/rand"
	"strconv"
	"time"
	//"crypto/sha1"
)


func TestChecksum(t *testing.T){
	var global_contents string = "jhfvjhdfjhfjjhjhdfvjvcvfjh"
	var result_expected string = "_68f3caf439065824dcf75651c202e9f7c28ebf07"
	file := HermesFile{}
	file.contents = global_contents;
	result := generate_file_checksum(file)
	if result != result_expected {
		t.Errorf("generate_file_checksum(\"jhfvjhdfjhfjjhjhdfvjvcvfjh\") failed expected %v got %v", result_expected, result)
	}
}

func TestFileName(t *testing.T){
	rand.Seed(time.Now().UnixNano())
	var file_ID int = rand.Intn(40)+10;
	var fake_checksum string ="_abba"
	var result_expected string =  "Hermes_"+strconv.Itoa(file_ID)+"_abba"
	result := generate_file_name(file_ID, fake_checksum)
	if result != result_expected {
		t.Errorf("generate_file_name(%v, \"abba\") failed expected %v got %v", file_ID, result_expected, result)
	}
	file_ID = rand.Intn(10);
	result_expected = "Hermes_0"+strconv.Itoa(file_ID)+"_abba"
	result = generate_file_name(file_ID, fake_checksum)
	if result != result_expected {
		t.Errorf("generate_file_name(%v, \"abba\") failed expected %v got %v", file_ID, result_expected, result)
	}
}
