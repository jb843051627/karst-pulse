package main

const serviceName = "karst-pulse"

func serviceInfo() map[string]string {
	return map[string]string{"name": serviceName, "runtime": "go1.22", "storage": "modernc.org/sqlite"}
}
