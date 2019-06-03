package puppet

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"syscall"
)

type Metadata struct {
	Checksum struct {
		Type  string     `json:"type"`
		Value string     `json:"value"`
	}                    `json:"checksum"`
	Destination  *string `json:"destination"`
	Group        uint32  `json:"group"`
	Links        string  `json:"links"`
	Mode         int     `json:"mode"`
	Owner        uint32  `json:"owner"`
	Path         string  `json:"path"`
	RelativePath *string  `json:"relative_path"`
	Type         string  `json:"type"`
}

func GetFileMetadata(path string) (Metadata, error) {
	if path == "plugins" {
		path = "/etc/puppetlabs/code/environments/production/modules/rogue/lib/puppet/functions/rogue/go_rogue.rb"
	}
	var metadata Metadata

	fileInfo, err := os.Stat(path)

	if err != nil {
		return metadata, err
	}

	populateCheckSum(fileInfo, &metadata, path)
	populatePathData(fileInfo, &metadata, path)
	populateOwnership(fileInfo, &metadata)
	metadata.Mode = 420
	metadata.Path = "/etc/puppetlabs/code/environments/production/modules/rogue/lib/puppet/functions/rogue/go_rogue.rb"

	fmt.Printf("%+v", metadata)
	return metadata, nil
}

func calculateMd5CheckSum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]

	return hex.EncodeToString(hashInBytes), nil
}

func populateCheckSum(fileInfo os.FileInfo, metadata *Metadata, path string) error {
	if fileInfo.IsDir() {

	} else {
		md5, err := calculateMd5CheckSum(path)
		if err != nil {
			return err
		}

		metadata.Checksum.Type = "md5"
		metadata.Checksum.Value = "{md5}" + md5
	}

	return nil
}

func populatePathData(fileInfo os.FileInfo, metadata *Metadata, path string) error {
	metadata.Links = "manage"

	fileInfo, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if fileInfo.Mode() & os.ModeSymlink == os.ModeSymlink {
		var link string
		link, err = os.Readlink(path)
		metadata.Destination = &link
		metadata.Type = "link"

		if err != nil {
			return err
		}
	} else {
		metadata.Destination = nil
		if fileInfo.IsDir() {
			metadata.Type = "directory"
		} else {
			metadata.Type = "file"
		}
	}

	return nil
}

func populateOwnership(fileInfo os.FileInfo, metadata *Metadata) {
	fmt.Println(fileInfo.Sys())
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		fmt.Println("OK")
		metadata.Group = stat.Gid
		metadata.Owner = stat.Uid
	} else {
		fmt.Println("NOT OK")
	}
}