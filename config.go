package main

type StampyConfig struct {
	Buckets int `yaml:"buckets"`
	Ip      string `yaml:"ip"`
	Port    int `yaml:"port"`
}

