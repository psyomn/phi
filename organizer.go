/*
Package phi is a set of tools to get incoming pictures from devices to
your server, and manage said space on those peripheral devices.

Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package phi

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
)

var supportedFormats = [...]string{
	"jpg$", "jpeg$",
	"png$", "gif$",
	"mpeg$", "mp4$",
	"mkv$", "avi$",
	"webp$",
}

func getUserName() string {
	currentUser, err := user.Current()
	userStr := ""
	if err != nil {
		log.Println("warning: username could not be found; using 'default' as username")
		userStr = "default"
	} else {
		userStr = currentUser.Username
	}
	return userStr
}

func isSupportedFormat(name string) bool {
	matchedAll := false
	for _, format := range supportedFormats {
		matched, err := regexp.MatchString(format, name)
		if err != nil {
			log.Println("error with regexp: ", err)
		} else {
			matchedAll = matched || matchedAll
		}
	}
	return matchedAll
}

// SortByModTime will detect files with designated formats and place
// them chronologically in an output directory, in the form of:
//   $OUTDIR/username/yyyy/mm/
func SortByModTime(dirPath, outDir string) {
	userStr := getUserName()
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(dirPath, ": ", err)
			return err
		}

		if !isSupportedFormat(path) {
			return nil
		}

		if info.Mode().IsRegular() {
			year := info.ModTime().Year()
			month := info.ModTime().Month()
			baseDir := fmt.Sprintf("%s/%s/%d/%02d", outDir, userStr, year, month)

			err := os.MkdirAll(baseDir, 0700)
			if err != nil {
				log.Println("problem creating dir: ", baseDir, ": ", err)
			}

			_, fileName := filepath.Split(path)
			toMovePath := filepath.Join(baseDir, fileName)
			log.Println("moving: ", path, "to: ", toMovePath)

			// TODO will have to deal with this at some point
			//
			// Safety: until we have something to deal
			// with duplicates, need to think about how to
			// deal with this sort of thing.
			if _, err := os.Stat(toMovePath); !os.IsNotExist(err) {
				log.Println("avoiding overwritting for now")
				return nil
			}

			err = os.Rename(path, toMovePath)
			if err != nil {
				log.Println("problem moving file: ", err)
			}
		}

		return nil
	})
}
