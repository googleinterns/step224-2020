package hermes


import (
	"fmt"
	"strconv"
	"crypto/sha1"
)
type HermesFile struct {
	name string
	contents string
}
func (file *HermesFile) GenerateFileName(file_ID int, file_checksum string) {
	id_string :=""
	if file_ID < 10 {
		id_string = "0" + strconv.Itoa(file_ID)
	} else {
		id_string = strconv.Itoa(file_ID)
	}
	file_name := "Hermes_"+id_string+file_checksum
	file.name = file_name
}
func (file *HermesFile) GenerateFileContents(file_ID int) {
	file.contents = "jhfvjhdfjhfjjhjhdfvjvcvfjh";
}

func (file HermesFile) GenerateFileChecksum() string {
	
	file_contents := []byte(file.contents)
	hash := sha1.Sum(file_contents)
	checksum := fmt.Sprintf("_%x", hash)
	return checksum;
}
func GenerateHermesFile (id int) HermesFile {
	file := HermesFile{}
	file.GenerateFileContents(id)
	checksum := file.GenerateFileChecksum()
	file.GenerateFileName(id, checksum)
	return file
}
