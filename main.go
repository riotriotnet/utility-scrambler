package main

import (
  "os"
  "log"
  "io/ioutil"
  "time"
  "path/filepath"
  "encoding/json"
)

type Job struct {
  Name             string              `json:"name"`
  TargetLocation   string              `json:"targetlocation"`
  ResultLocation   string              `json:"resultlocation"`
  ConfigLocation   string              `json:"configlocation"`
  ExtensionMapping map[string][]string `json:"extensionmapping"`
  Created          time.Time           `json:"created"`
  Updated          time.Time           `json:"updated"`
  Completed        time.Time           `json:"completed"`
}

func main() {
  job := NewJob("./config.json")
  job.Run()
}

func NewJob(jobconfig string) *Job {
  return &Job{
    ConfigLocation: jobconfig,
    Created: time.Now(),
    ExtensionMapping: make(map[string][]string),
  }
}

func (job *Job) Run() {
  job.LoadConfiguration()
  job.ProcessDirectory()
}

func (job *Job) LoadConfiguration() {
	b, err := ioutil.ReadFile( job.ConfigLocation )
  if err != nil {
      log.Print(err)
  }
  err = json.Unmarshal(b, &job)
  if err != nil {
    panic(err)
  }
}

func (job *Job) ProcessDirectory() {
  files, _ := ioutil.ReadDir( job.TargetLocation )
  log.Println("Files:", len(files))
  if len(files) > 0 {
    for _, f := range files {
      job.MoveFile(f.Name())
    }
  }
}

func (job *Job) MoveFile(filename string) {
  subFolder := job.GetMapping(filepath.Ext(filename))
  if subFolder != "" {
    ef := job.TargetLocation + filename
    sl := job.ResultLocation + subFolder + "/" + filename
    log.Println("[SOURCE]: ", ef, " -> [TARGET]: ", sl)
    err :=  os.Rename(ef, sl)
    if err != nil {
      log.Println(err)
      return
    }
  } else {
    log.Println("Missing Extension:", filepath.Ext(filename))
  }
}

func (job *Job) MapNewExtension(ext string, group string) {
  if _, ok := job.ExtensionMapping[group]; ok {
    if contains(job.ExtensionMapping[group], ext) {
      return
    }
    job.ExtensionMapping[group] = append(job.ExtensionMapping[group], ext)
  } else {
    if len(job.ExtensionMapping[group]) == 0 {
      job.ExtensionMapping[group] = []string{}
    }
    job.ExtensionMapping[group] = append(job.ExtensionMapping[group], ext)
  }
  log.Println(job.ExtensionMapping)
}

func (job *Job) GetMapping(ext string) string {
  response := ""
  for index, array := range job.ExtensionMapping {
    if contains(array, ext) {
      return index
    }
  }
  return response
}

func contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }
    _, ok := set[item]
    return ok
}

func StorageWriteJsonFile(country string, kind string, payload []byte) {
	// filename := country+"-"+kind+"-"+CreateUUID()+".json"
  //   err := ioutil.WriteFile(api_storagelocation+"/"+filename, payload, 0644)
	// if err != nil {
	// 	log.Println(err)
	// }
}
