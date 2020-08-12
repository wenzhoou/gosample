package main

import (
  "path/filepath"
  "os"
  "log"
  "encoding/csv"
  "strconv"
  "encoding/json"
  "regexp"
  "strings"
)

type Configuration struct {
    Input                   string
    Output                  string
    IgnoreDirsName          []string
    IgnoreDirsPattern       []string
    IgnoreDirsPatternCpl    []* regexp.Regexp
    IgnoreFilesName         []string
    IgnoreFilesPattern      []string
    IgnoreFilesPatternCpl   []* regexp.Regexp
}

var (
  Config                 Configuration;
  OutputWriter           * csv.Writer;
)

func main() {
  Config = readConfig()
  
  file, err := os.Create(Config.Output)
  checkError("Cannot create file", err)
  defer file.Close()

  OutputWriter = csv.NewWriter(file)
  defer OutputWriter.Flush()
  
  Config.IgnoreDirsPatternCpl = make([]* regexp.Regexp, len(Config.IgnoreDirsPattern))
  for i, s := range Config.IgnoreDirsPattern {
    log.Println("Add Filter Element", s)
    Config.IgnoreDirsPatternCpl[i] = regexp.MustCompile("(?i)" + s)
  }
  
  Config.IgnoreFilesPatternCpl = make([]* regexp.Regexp, len(Config.IgnoreFilesPattern))
  for i, s := range Config.IgnoreFilesPattern {
    log.Println("Add Filter Element", s)
    Config.IgnoreFilesPatternCpl[i] = regexp.MustCompile("(?i)" + s)
  }
    
  err = filepath.Walk(Config.Input, dealFile)
  checkError("Walk error : ", err)
}

func readConfig() Configuration {
  file, _ := os.Open("conf.json")
  defer file.Close()
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  err := decoder.Decode(&configuration)
  checkError("Read Config error : ", err)
  return configuration;
}

func checkError(message string, err error) {
  if err != nil {
    log.Fatal(message, err)
  }
}

func needIgnore(path string, name string, isDir bool) bool {

  log.Println("Process File ", path, name, isDir)

  if !isDir {
    for _, igr := range Config.IgnoreFilesName {
      
      if strings.EqualFold(igr,path) || strings.EqualFold(igr,name) {
        return true 
      }
    }
    for _, re := range Config.IgnoreFilesPatternCpl {
      
      if re.MatchString(path) {
        return true 
      }
    }
  }
  
  if isDir {
    for _, igr := range Config.IgnoreDirsName {
      
      if strings.EqualFold(igr,path) || strings.EqualFold(igr,name) {
        return true 
      }
    }
    for _, re := range Config.IgnoreDirsPatternCpl {
      
      if re.MatchString(path) {
        return true 
      }
    }
  }
  log.Println("passed ")
  return false
}

func dealFile(fullpath string, info os.FileInfo, err error) error {

  checkError("Cannot get file : ", err)
  path := filepath.Dir(fullpath)
  
  if needIgnore(fullpath, info.Name(), info.IsDir()) {
    return filepath.SkipDir
  }
  if !info.IsDir() {
    err = OutputWriter.Write([]string{path, info.Name(), strconv.FormatInt(info.Size() ,10), info.ModTime().Format("2006-01-02 15:04:05")})
    checkError("Cannot write file : ", err)
  }
  return nil
}