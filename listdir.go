package main

import (
  "fmt"
  "path/filepath"
  "os"
  "log"
)

func main() {

  err := filepath.Walk(".",
    func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path, info.Size())
    return nil
  })
  if err != nil {
    log.Println(err)
  }
}
