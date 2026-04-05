package aof

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/KrishPatel10/GoRedis.git/internal/resp"
)

type AOF struct {
	file *os.File
	mu   sync.Mutex
}

func NewAOF(path string) (*AOF, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	aof := &AOF{
		file: file,
	}

	go aof.startSync(time.Minute)

	return aof, nil
}

func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.file.Sync()
	return a.file.Close()
}

func (a *AOF) Write(cmd []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, err := a.file.Write(cmd)

	return err
}

func (a *AOF) Read(callback func(value resp.Value)) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.file.Seek(0, 0)

	parser := resp.NewParser(a.file)

	for {
		val, err := parser.Parse()

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		callback(val)
	}

	return nil
}

func (a *AOF) startSync(interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		a.sync()
	}
}

func (a *AOF) sync() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	err := a.file.Sync()

	if err != nil {
		fmt.Printf("Error happened during syncing file: \n %s", err)
	}

	return nil
}
