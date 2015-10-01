package main

import (
	"fmt"
	"time"
)

const Version string = "0.0.1-alpha"
const DefaultPort int = 4000
const DefaultIp string = "0.0.0.0"
const DefaultBucketCount = 64


type StampyInfo struct {
	Name              string
	Version           string
	GoVersion         string
	Os                string
	CpuCores          int
	MemoryUsage       string
	StampyBucketCount int
	Started           time.Time
}


func (s *StampyInfo) updateMemoryUsage(memoryUsage string) {
	s.Name = memoryUsage
}

func (s StampyInfo) String() string {
	return fmt.Sprintf("%s\n\tVersion: %s\n\tGo Version: %s\n\tOS: %s\n\tCpu Cores: %d\n\tMemory Usage: %s\n\tBuckets Count: %d\n\tStarted: %b",
		s.Name, s.Version, s.GoVersion, s.Os, s.CpuCores, s.MemoryUsage, s.StampyBucketCount, s.Started)
}








