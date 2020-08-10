package main

import (
  "fmt"
  "path/filepath"
  "os"
  "log"
  "flag"
)

func main() {
  wordPtr := flag.String("d", ".", "directory")
  flag.Parse()
  
  err := filepath.Walk(*wordPtr,
    func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path, info.Size(), )
    return nil
  })
  if err != nil {
    log.Println(err)
  }
}
