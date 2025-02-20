package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/koki-develop/go-fzf"
)

func main() {
	db := DB{file: "~/.cache/sesh.db"}
	f, err := db.open()
	if err != nil {
		fmt.Printf("Error: %v\n", f)
	}
	defer f.Close()
	if len(os.Args) == 1 {
		item, err := db.find(f)
		if err != nil {
			fmt.Printf("Err: %v\n", err)
		}
		fmt.Println("[DIR]: ", item)
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add":
			if len(os.Args) == 3 {
				err := db.add(f, os.Args[2])
				if err != nil {
					fmt.Printf("Err: %v\n", err)
				}
			} else {
				fmt.Println("Err: Invalid amount of args to call add")
			}
		case "rm":
			if len(os.Args) == 3 {
				err := db.rm(f, os.Args[2])
				if err != nil {
					fmt.Printf("Err: %v\n", err)
				}
			} else {
				fmt.Println("Err: Invalid amount of args to call add")
			}
		default:
			fmt.Printf("hello")
		}
	}
}

type DB struct {
	file string
}

func (db *DB) open() (*os.File, error) {
	expandedPath := db.file
	if expandedPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %v", err)
		}
		expandedPath = homeDir + expandedPath[1:]
	}

	f, err := os.OpenFile(expandedPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		fmt.Println("Initializing new database file...")
		_, err := f.Write([]byte{0, 0, 0, 0})
		if err != nil {
			f.Close()
			return nil, err
		}
	}
	return f, nil
}

func (db *DB) get(f *os.File) (map[string]bool, error) {
	_, err := f.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return make(map[string]bool), nil
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	res := make(map[string]bool)
	var entries int32
	if err := binary.Read(buf, binary.NativeEndian, &entries); err != nil {
		if err.Error() == "EOF" {
			return make(map[string]bool), nil
		}
		return nil, err
	}
	for i := 0; i < int(entries); i++ {
		var keyLen int32
		if err := binary.Read(buf, binary.NativeEndian, &keyLen); err != nil {
			return nil, err
		}
		keyBytes := make([]byte, keyLen)
		if err := binary.Read(buf, binary.NativeEndian, &keyBytes); err != nil {
			return nil, err
		}
		key := string(keyBytes)
		var boolByte byte
		if err := binary.Read(buf, binary.NativeEndian, &boolByte); err != nil {
			return nil, err
		}
		val := boolByte == 1
		res[key] = val
	}
	return res, nil
}

func (db *DB) add(f *os.File, newPath string) error {
	info, err := os.Stat(newPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("File is not a directory")
	}
	absPath, _ := filepath.Abs(newPath)
	existingEntries, err := db.get(f)
	if err != nil {
		return err
	}
	if _, exists := existingEntries[absPath]; exists {
		fmt.Println("Key already exists:", absPath)
		return nil
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	var numEntries int32
	if err := binary.Read(f, binary.NativeEndian, &numEntries); err != nil {
		return err
	}
	numEntries++
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	if err := binary.Write(f, binary.NativeEndian, numEntries); err != nil {
		return err
	}
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	keyBytes := []byte(absPath)
	keyLen := int32(len(keyBytes))
	if err := binary.Write(&buf, binary.NativeEndian, keyLen); err != nil {
		return err
	}
	if err := binary.Write(&buf, binary.NativeEndian, keyBytes); err != nil {
		return err
	}
	if err := binary.Write(&buf, binary.NativeEndian, byte(1)); err != nil {
		return err
	}
	_, err = f.Write(buf.Bytes())
	return err
}

func (db *DB) rm(f *os.File, keyToRemove string) error {
	existingEntries, err := db.get(f)
	if err != nil {
		return err
	}
	if _, exists := existingEntries[keyToRemove]; !exists {
		fmt.Println("Key not found:", keyToRemove)
		return nil
	}
	delete(existingEntries, keyToRemove)
	tempFile, err := os.CreateTemp("", "db_temp")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	var buf bytes.Buffer
	numEntries := int32(len(existingEntries))
	if err := binary.Write(&buf, binary.NativeEndian, numEntries); err != nil {
		return err
	}
	for key, value := range existingEntries {
		keyBytes := []byte(key)
		keyLen := int32(len(keyBytes))
		if err := binary.Write(&buf, binary.NativeEndian, keyLen); err != nil {
			return err
		}
		if err := binary.Write(&buf, binary.NativeEndian, keyBytes); err != nil {
			return err
		}
		var boolByte byte
		if value {
			boolByte = 1
		}
		if err := binary.Write(&buf, binary.NativeEndian, boolByte); err != nil {
			return err
		}
	}
	_, err = tempFile.Write(buf.Bytes())
	if err != nil {
		return err
	}
	f.Close()
	tempFile.Close()
	return os.Rename(tempFile.Name(), db.file)
}

func (db *DB) find(f *os.File) (string, error) {
	entries, err := db.get(f)
	if err != nil {
		return "", err
	}
	items := []string{}
	for e := range entries {
		items = append(items, e)
	}
	finder, err := fzf.New()
	if err != nil {
		return "", err
	}
	idxs, err := finder.Find(items, func(i int) string { return items[i] })
	if err != nil {
		return "", err
	}
	if len(idxs) == 0 {
		fmt.Println("Noting selected")
		os.Exit(0)
	}
	return items[idxs[0]], nil
}
