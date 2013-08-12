package main

// “It is what you read when you don't have to that determines what you will
//  be when you can't help it.”
// ― Oscar Wilde

import (
	"fmt"
	"github.com/APTrust/bagins"
	"github.com/APTrust/bagins/bagutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func main() {

	src := "/Users/swt8w/Downloads/duracloudfiles/uva"
	dst := "/Users/swt8w/Downloads/duracloudfiles/bags"

	fb, err := NewFedoraBagger(src, "md5", dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	fb.MakeBags()
}

type FedoraBagger struct {
	src string // path of fedora datastream files to bag
	hsh string // name of hashtype to use for checksums
	dst string // path to write bags to
}

func NewFedoraBagger(src string, hashName string, dst string) (*FedoraBagger, error) {

	if _, err := os.Stat(filepath.Clean(src)); err != nil {
		return nil, fmt.Errorf("Error reading soruce directory: %v", err)
	}
	if _, err := os.Stat(filepath.Clean(dst)); err != nil {
		return nil, fmt.Errorf("Error reading destination directory: %v", err)
	}

	b := new(FedoraBagger)
	b.src = src
	b.hsh = hashName
	b.dst = dst

	return b, nil
}

// Create bags from the groups of Fedora Object datastreams in the source
// directory or return an error if none were found.
func (b *FedoraBagger) MakeBags() error {

	bagList := walkFedora("uva", b.src)
	if len(bagList) < 1 {
		return fmt.Errorf("Unable to find files to bag in %s", b.src)
	}

	for key, files := range bagList {
		cs, err := bagutil.NewCheckByName(b.hsh)
		if err != nil {
			return fmt.Errorf("Error creating checksum hash:", err)
		}

		bag, err := bagins.NewBag(b.dst, key, cs)
		if err != nil {
			return fmt.Errorf("Error creating bag:", err)
		}

		tfName := "bag-info.txt"
		bag.AddTagfile(tfName)
		tf, err := bag.TagFile(tfName)
		if err != nil {
			return fmt.Errorf("Error getting tagfile:", err)
		}
		tf.Data = bagInfoData()

		for _, file := range files {
			srcPath := filepath.Join(b.src, file)
			bag.AddFile(srcPath, url.QueryEscape(file))
		}

		bag.Close()
	}
	return nil
}

// Walks a directory and tries to gather information about serialized
// fedora datatream files in a directory.  It groups each file into a
// map with a key consisting of the destination bag name and a value
// of the filepath to be added to that bag.
func walkFedora(code string, src string) map[string][]string {

	bagFiles := make(map[string][]string)

	visit := func(pth string, info os.FileInfo, vErr error) error {
		// THINK this regex will work (^.*?)(\+|$)
		// Matches a fedora PID from a datastream filename
		re, err := regexp.Compile(`(^.*?)(\+|$)`) // Now I have 2 problems
		if err != nil {
			return err
		}

		matches := re.FindStringSubmatch(info.Name())
		if len(matches) < 2 {
			return fmt.Errorf("Bad match on filename: %s", info.Name())
		}

		key := code + "_" + url.QueryEscape(matches[1])
		bagFiles[key] = append(bagFiles[key], info.Name())

		return vErr
	}

	if err := filepath.Walk(src, visit); err != nil {
		fmt.Println(err)
	}

	return bagFiles
}

// Returns a string formatted by the current date.
func stringTimeNow() string {
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	t := time.Now()
	return t.UTC().Format(layout)
}

func bagInfoData() map[string]string {
	data := make(map[string]string)
	data["Source-Organization"] = "University of Virginia"
	data["Bagging-Date"] = stringTimeNow()

	return data
}
