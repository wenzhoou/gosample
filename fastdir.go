package main

import (
  "github.com/karrick/godirwalk"
  "os"
  "log"
  "flag"
  "encoding/csv"
  "strconv"
  "path/filepath"
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
  godirwalk.Walk(*input, &godirwalk.Options{
    Unsorted: true,
    Callback: func(osPathname string, de *godirwalk.Dirent) error {
      checkError("Cannot get file : ", err)
      info, err := os.Stat(osPathname)
      checkError("Cannot get fileInfo : ", err)
      err = writer.Write([]string{osPathname, de.Name(), strconv.FormatInt(info.Size() ,10), info.ModTime().Format("2006-01-02 15:04:05"), strconv.FormatBool(info.IsDir())})
      checkError("Cannot write file : ", err)
      return nil
    },
    ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
      checkError("Walk error : ", err)
      return godirwalk.SkipNode
    },
  })
}

func checkError(message string, err error) {
  if err != nil {
    log.Fatal(message, err)
  }
}