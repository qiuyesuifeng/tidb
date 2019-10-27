// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Mapped from a file path, which is in the format
// of /{root}/{ip}/{component}:{port}/{xxx}.log
type FileWrapper struct {
	Root     string
	Host     string
	Folder   string
	Filename string
}

// Open the file fw represent.
func (fw *FileWrapper) Open() (*os.File, error) {
	filePath := path.Join(fw.Root, fw.Host, fw.Folder, fw.Filename)
	return os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
}

// Return the component name and port it listening on.
func (fw *FileWrapper) ParseFolderName() (string, string, error) {
	//s := strings.Split(fw.Folder, "-")
	//if len(s) < 2 {
	//	return "", "", fmt.Errorf("unexpect folder name: %s", s)
	//}
	//return s[0], s[1], nil
	return "tidb", "", nil
}

func NewFileWrapper(root, host, folder, filename string) *FileWrapper {
	return &FileWrapper{
		Root:     root,
		Host:     host,
		Folder:   folder,
		Filename: filename,
	}
}

// Traversing a folder and parse it's struct, generating
// a list of file wrapper.
func ResolveDir(src string) ([]*FileWrapper, error) {
	var wrappers []*FileWrapper
	dir, err := ioutil.ReadDir(src) // {cluster_uuid}
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		filename := fi.Name()
		if !strings.Contains(filename, "tidb.log") {
			continue
		}
		fw := NewFileWrapper(src, "", "", filename)
		wrappers = append(wrappers, fw)
	}
	return wrappers, nil
}
