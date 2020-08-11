package main

import (
  "path/filepath"
  "os"
  "log"
  "flag"
  "encoding/csv"
  "strconv"
)

func main() {
  input := flag.String("d", ".", "directory")
  output := flag.String("o", ".", "directory")
  flag.Parse()

  file, err := os.Create(filepath.Join(*output, "result.csv"))
  checkError("Cannot create file", err)
  defer file.Close()

  writer := csv.NewWriter(file)
  defer writer.Flush()

  
  err = filepath.Walk(*input,
    func(path string, info os.FileInfo, err error) error {
      checkError("Cannot get file : ", err)
      err = writer.Write([]string{path, info.Name(), strconv.FormatInt(info.Size() ,10), info.ModTime().Format("2006-01-02 15:04:05"), strconv.FormatBool(info.IsDir()) })
      checkError("Cannot write file : ", err)
      return nil
  })
  checkError("Walk error : ", err)
}


func checkError(message string, err error) {
  if err != nil {
    log.Fatal(message, err)
  }
}