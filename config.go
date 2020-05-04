package main

type OutletConfig struct {
	Name  string
	Type  string
	Group int
	Id    int
}

type MonitorConfig struct {
	Name string
	Type string
}

type FermenterConfig struct {
	Name         string
	Monitor      string
	HeaterOutlet string
	CoolerOutlet string
}
