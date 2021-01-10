package main

import (
  "path/filepath"
  "os"
  log "github.com/sirupsen/logrus"
  "encoding/csv"
  "strconv"
  "regexp"
  "strings"
  "github.com/gookit/config"
  "github.com/gookit/config/yaml"
)

type ConfType struct {
    Path                    string
    Type                    string
}

type ConfigInfo struct {
    Input                   []string
    Output                  string
    Ignore                  Ignore
    ConfType                []ConfType
}

type Ignore struct {
  Dirs   IgnoreItem
  Files  IgnoreItem
}

type IgnoreItem struct {
  Name      []string
  Regex     []string
  RegexCpl  []* regexp.Regexp
}

var (
  Config                 ConfigInfo
  OutputWriter           * csv.Writer
)

func main() {
  Config = readConfig()
  setLoglevel()
  log.WithFields(log.Fields{
    "config":    Config,
  }).Info("Init OK")
  
  file, err := os.Create(Config.Output)
  checkError("Cannot create output file", err)
  defer file.Close()

  OutputWriter = csv.NewWriter(file)
  defer OutputWriter.Flush()
  
  Config.Ignore.Dirs.RegexCpl = make([]* regexp.Regexp, len(Config.Ignore.Dirs.Regex))
  for i, s := range Config.Ignore.Dirs.Regex {
    log.Debugln("Add Dir Filter Element", s)
    Config.Ignore.Dirs.RegexCpl[i] = regexp.MustCompile("(?i)" + s)
  }
  
  Config.Ignore.Files.RegexCpl = make([]* regexp.Regexp, len(Config.Ignore.Files.Regex))
  for i, s := range Config.Ignore.Files.Regex {
    log.Debugln("Add File Filter Element", s)
    Config.Ignore.Files.RegexCpl[i] = regexp.MustCompile("(?i)" + s)
  }
  
  for _, s := range Config.Input {
    log.Debugln("Walk ", s)
    err = filepath.Walk(s, dealFile)
    checkError("Walk error : ", err)
  }
}

func readConfig() ConfigInfo {

  config.WithOptions(config.ParseEnv)
  config.AddDriver(yaml.Driver)

  err := config.LoadFiles("conf.yaml")
  checkError("Failed load conf.yaml", err)
  conf := ConfigInfo{};
  err = config.BindStruct("", &conf)
  checkError("Failed parse conf.yaml", err)
  return conf;
}

func setLoglevel() {
  //log.SetReportCaller(true)
  log.SetFormatter(&log.TextFormatter{ForceColors: true})
  //log.SetOutput(colorable.NewColorableStdout())
  log.SetOutput(os.Stdout)
  switch config.String("loglevel") {
    case "debug":
      log.SetLevel(log.DebugLevel)
    case "info":
      log.SetLevel(log.InfoLevel)
    case "warn":
      log.SetLevel(log.WarnLevel)
    case "error":
      log.SetLevel(log.ErrorLevel)
    default:
      log.SetLevel(log.InfoLevel)
      log.Warnf("Home: invalid log level supplied: '%s'", config.String("loglevel"))
  }
}


func checkError(message string, err error) {
  if err != nil {
    log.Errorf(message, err)
    panic(err)
  }
}

func needIgnoreFile(path string, name string, isDir bool) bool {

  log.Debug("Process File ", path, name, isDir)

  if !isDir {
    for _, igr := range Config.Ignore.Files.Name {
      
      if strings.EqualFold(igr,path) || strings.EqualFold(igr,name) {
        log.WithFields(log.Fields{"path": path, "name" : name, "isDir": isDir }).Info("Ignore File")
        return true 
      }
    }
    for _, re := range Config.Ignore.Files.RegexCpl {
      
      if re.MatchString(path) {
        log.WithFields(log.Fields{"path": path, "name" : name, "isDir": isDir }).Info("Ignore File")
        return true 
      }
    }
  }
  return false
}

func needIgnoreDir(path string, name string, isDir bool) bool {

  if isDir {
    for _, igr := range Config.Ignore.Dirs.Name {
      
      if strings.EqualFold(igr,path) || strings.EqualFold(igr,name) {
        log.WithFields(log.Fields{"path": path, "name" : name, "isDir": isDir }).Info("Ignore Dir ")
        return true 
      }
    }
    for _, re := range Config.Ignore.Dirs.RegexCpl {
      
      if re.MatchString(path) {
        log.WithFields(log.Fields{"path": path, "name" : name, "isDir": isDir }).Info("Ignore Dir ")
        return true 
      }
    }
  }
  //log.Debug("passed ")
  return false
}

func dealFile(fullpath string, info os.FileInfo, err error) error {

  checkError("Cannot get file : ", err)
  path := filepath.Dir(fullpath)
  if needIgnoreDir(fullpath, info.Name(), info.IsDir()) {
    log.Debug("Ignore ", fullpath)
    return filepath.SkipDir
  }
  if !info.IsDir() && !needIgnoreFile(fullpath, info.Name(), info.IsDir()) {
    log.Debug("Add ", info.Name())
    err = OutputWriter.Write([]string{path, info.Name(), filepath.Ext(info.Name()), getType(fullpath), strconv.FormatInt(info.Size() ,10), info.ModTime().Format("2006-01-02 15:04:05")})
    checkError("Cannot write file : ", err)
  }
  return nil
}


func getType(fullpath string) string {
  for _, v := range Config.ConfType {
    if strings.HasPrefix(fullpath, v.Path) {
      return v.Type
    }
  }
  return ""
}

func dumpMap(space string, m map[string]string) {
  for k, v := range m {
    log.Debug(space, k, ":", v)
  }
}

func dumpArray(space string, a []interface{}) {
  log.Debug(space, "[")
  for _, v := range a {
    log.Debug(space+"\t", v, ",")
  }
  log.Debug(space,"]")
}
