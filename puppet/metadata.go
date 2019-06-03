package metadata

import (
	"fmt"
	"os"
	"syscall"
)

type Metadata struct {
	Checksum struct {
		Type  string    `json:"type"`
		Value string    `json:"value"`
	}                   `json:"checksum"`
	Destination  string `json:"destination"`
	Group        uint32 `json:"group"`
	Links        string `json:"links"`
	Mode         int    `json:"mode"`
	Owner        uint32 `json:"owner"`
	Path         string `json:"path"`
	RelativePath string `json:"relative_path"`
	Type         string `json:"file"`
}

func GetFileMetadata(path string) (Metadata, error) {
	var metadata Metadata

	fileInfo, err := os.Stat(path)

	if err != nil {
		return metadata, err
	}

	populateLinkData(&metadata, path)
	populateOwnership(fileInfo, &metadata)

	if fileInfo.IsDir() {
	}

	fmt.Println("%+v", metadata)
	return metadata, nil
}

//fileinfo, _ := os.Stat(file)
//fileinfo.Sys()
//stat, ok := fileinfo.Sys().(*syscall.Stat_t)
//if !ok {
//	fmt.Printf("Not a syscall.Stat_t")
//	return
//}

func populateCommonMetadata(metadata *Metadata, fileInfo os.FileInfo) {
	metadata.Links = "manage"


}

func populateLinkData(metadata *Metadata, path string) error {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if fileInfo.Mode() & os.ModeSymlink == os.ModeSymlink {
		metadata.Links, err = os.Readlink(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func populateOwnership(fileInfo os.FileInfo, metadata *Metadata) {
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		metadata.Group = stat.Gid
		metadata.Owner = stat.Uid
	}
}