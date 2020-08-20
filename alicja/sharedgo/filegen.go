package main
import (
	"fmt"
	// "math/rand"
	"strconv"
	"crypto/sha1"
	"encoding/hex"
)
type HermesFile struct {
	name string
	contents string
}
func generate_file_name (file_ID int, file_checksum string) string{
	var id_string string = strconv.Itoa(file_ID)
	var file_name string = "Hermes_"+id_string+file_checksum
	return file_name 
}
func generate_file_contents (file_ID int) string{
	return "jhfvjhdfjhfjjhjhdfvjvcvfjh";
}
func generate_file_checksum (file HermesFile) string{
	
	file_contents := []byte(file.contents)
	hash := sha1.Sum(file_contents)
	// checksum  :=hex.EncodeToString(hash) as an alternative
	checksum := fmt.Sprintf("%s%v","_", hash)
	return checksum;
}
func generate_file (id int) HermesFile {
	file := HermesFile{}
	file.contents = generate_file_contents(id)
	var checksum string = generate_file_checksum(file)
	file.name = generate_file_name(id, checksum)
	return file
}
func main(){
	var f HermesFile = generate_file(20)
	fmt.Println(f.name)
}
