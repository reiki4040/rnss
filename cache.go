package main

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

const (
	BaseDir  = ".rnss"
	CacheDir = "cache"
)

func NewEC2Cache(region string) (*EC2Cache, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	fName, err := GetCacheFileName(region)
	if err != nil {
		return nil, err
	}
	absoluteCacheDirPath := filepath.Join(u.HomeDir, BaseDir, CacheDir)
	if _, err := os.Stat(absoluteCacheDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(absoluteCacheDirPath, 0700)
		if err != nil {
			if !os.IsExist(err) {
				return nil, err
			}
		}
	}

	p := filepath.Join(absoluteCacheDirPath, fName)
	return &EC2Cache{
		cachePath: p,
		l:         sync.Mutex{},
	}, nil
}

type EC2Cache struct {
	cachePath string

	l sync.Mutex
}

// if no cache file then return os.NotExists error
func (c *EC2Cache) Get() ([]string, error) {
	if _, err := os.Stat(c.cachePath); err != nil {
		return nil, err
	}

	fd, err := os.Open(c.cachePath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var list []string
	s := bufio.NewScanner(fd)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		if len(l) > 0 {
			list = append(list, l)
		}
	}

	return list, nil
}

func (c *EC2Cache) Store(list []string) error {
	c.l.Lock()
	defer c.l.Unlock()

	fd, err := os.Create(c.cachePath)
	if err != nil {
		return err
	}
	defer fd.Close()
	w := bufio.NewWriter(fd)
	defer w.Flush()

	for _, l := range list {
		w.WriteString(l + "\n")
	}

	return nil
}

func GetCacheFileName(region string) (string, error) {
	// ap-notheast-1_ec2list.tsv
	return strings.Join([]string{region, "ec2list.tsv"}, "_"), nil
}
