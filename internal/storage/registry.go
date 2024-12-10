package storage

func GetValidKeys() map[string]string {
	gauge := "gauge"
	counter := "counter"

	gaugelist := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
		"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups",
		"MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc",
		"RandomValue"}
	counterList := []string{"PollCount"}

	m := map[string]string{}

	for _, key := range gaugelist {
		m[key] = gauge
	}
	for _, key := range counterList {
		m[key] = counter
	}
	return m
}
