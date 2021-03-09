/*-----------------------------------------------------------------------------------
  --  HttpTarget                                                                       --
  --  Copyright (C) 2021  HttpTarget's Contributors                                    --
  --                                                                               --
  --  This program is free software: you can redistribute it and/or modify         --
  --  it under the terms of the GNU Affero General Public License as published     --
  --  by the Free Software Foundation, either version 3 of the License, or         --
  --  (at your option) any later version.                                          --
  --                                                                               --
  --  This program is distributed in the hope that it will be useful,              --
  --  but WITHOUT ANY WARRANTY; without even the implied warranty of               --
  --  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the                --
  --  GNU Affero General Public License for more details.                          --
  --                                                                               --
  --  You should have received a copy of the GNU Affero General Public License     --
  --  along with this program.  If not, see <https:   -- www.gnu.org/licenses/>.   --
  -----------------------------------------------------------------------------------*/

package static

import (
	"embed"
	"fmt"
	"github.com/hyperjumptech/httptarget/static/mime"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var (
	errFileNotFound = fmt.Errorf("file not found")
)

//go:embed api
var fs embed.FS

type FileData struct {
	Bytes       []byte
	ContentType string
}

func IsDir(path string) bool {
	for _, s := range GetPathTree("static") {
		if s == "[DIR]"+path {
			return true
		}
	}
	return false
}

func GetPathTree(path string) []string {
	logrus.Infof("Into %s", path)
	var entries []os.DirEntry
	var err error
	if strings.HasPrefix(path, "./") {
		entries, err = fs.ReadDir(path[2:])
	} else {
		entries, err = fs.ReadDir(path)
	}
	ret := make([]string, 0)
	if err != nil {
		return ret
	}
	logrus.Infof("Path %s %d etries", path, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			ret = append(ret, "[DIR]"+path+"/"+e.Name())
			ret = append(ret, GetPathTree(path+"/"+e.Name())...)
		} else {
			ret = append(ret, path+"/"+e.Name())
		}
	}
	return ret
}

func GetFile(path string) (*FileData, error) {
	bytes, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	mimeType, err := mime.MimeForFileName(path)
	if err != nil {
		return &FileData{
			Bytes:       bytes,
			ContentType: http.DetectContentType(bytes),
		}, nil
	}
	return &FileData{
		Bytes:       bytes,
		ContentType: mimeType,
	}, nil
}
