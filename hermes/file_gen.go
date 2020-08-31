package hermes

import (
	"fmt"
	"crypto/sha1"
)

type HermesFile struct {
	name string
	contents string
}
func (file *HermesFile) generateFileName(file_ID int, file_checksum string) {
	file_name :=  fmt.Sprintf("Hermes_%02d%v", file_ID, file_checksum)
	file.name = file_name
}
func (file *HermesFile) generateFileContents(file_ID int) {
	file.contents = "jhfvjhdfjhfjjhjhdfvjvcvfjh";
}
func (file HermesFile) generateFileChecksum() string {
	
	file_contents := []byte(file.contents)
	hash := sha1.Sum(file_contents)
	checksum := fmt.Sprintf("%s%x","_", hash)
	return checksum;
}
func GenerateFile (id int) HermesFile {
	file := HermesFile{}
	file.generateFileContents(id)
	checksum := file.generateFileChecksum()
	file.generateFileName(id, checksum)
	return file
}

