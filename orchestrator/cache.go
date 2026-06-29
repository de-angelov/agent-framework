package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type FileCache struct {
	mu    sync.Mutex
	items map[string]cachedFile
}

type cachedFile struct {
	modTime time.Time
	size    int64
	value   string
}

func (c *FileCache) Read(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[path]; ok &&
		item.modTime.Equal(info.ModTime()) &&
		item.size == info.Size() {
		return item.value, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	value := string(data)
	c.items[path] = cachedFile{
		modTime: info.ModTime(),
		size:    info.Size(),
		value:   value,
	}

	return value, nil
}

func mustRead(path string) string {
	data, err := fileCache.Read(path)
	if err != nil {
		return fmt.Sprintf("[failed to read %s: %v]", path, err)
	}
	return data
}
