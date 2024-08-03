package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "runtime"
    "runtime/pprof"
    "syscall"
    "time"

    "github.com/fsnotify/fsnotify"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/spf13/viper"
    "html/template"
)

// MemoryStats holds custom memory statistics
type MemoryStats struct {
    Alloc         uint64
    TotalAlloc    uint64
    Sys           uint64
    Lookups       uint64
    Mallocs       uint64
    Frees         uint64
    HeapAlloc     uint64
    HeapSys       uint64
    HeapIdle      uint64
    HeapInuse     uint64
    HeapReleased  uint64
    HeapObjects   uint64
    StackInuse    uint64
    StackSys      uint64
    MSpanInuse    uint64
    MSpanSys      uint64
    MCacheInuse   uint64
    MCacheSys     uint64
    BuckHashSys   uint64
    GCSys         uint64
    OtherSys      uint64
    NextGC        uint64
    LastGC        uint64
    PauseTotalNs  uint64
    NumGC         uint32
    NumForcedGC   uint32
    GCCPUFraction float64
}

var (
    memAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "go_memory_alloc",
        Help: "Bytes of allocated heap objects.",
    })
    templates = template.Must(template.ParseFiles("dashboard.html"))
)

func init() {
    prometheus.MustRegister(memAlloc)
}

func initConfig() {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config file changed: %s", e.Name)
    })
}

func validateConfig() {
    if viper.GetDuration("interval") <= 0 {
        log.Fatal("Invalid interval value")
    }
    if viper.GetDuration("duration") <= 0 {
        log.Fatal("Invalid duration value")
    }
    if viper.GetString("profile") == "" {
        log.Fatal("Profile type not specified")
    }
    if viper.GetString("profileFile") == "" {
        log.Fatal("Profile file not specified")
    }
    if viper.GetString("logFile") == "" {
        log.Fatal("Log file not specified")
    }
    if viper.GetString("username") == "" || viper.GetString("password") == "" {
        log.Fatal("Username or password not specified")
    }
}

func CaptureMemoryStats() MemoryStats {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return MemoryStats{
        Alloc:         memStats.Alloc,
        TotalAlloc:    memStats.TotalAlloc,
        Sys:           memStats.Sys,
        Lookups:       memStats.Lookups,
        Mallocs:       memStats.Mallocs,
        Frees:         memStats.Frees,
        HeapAlloc:     memStats.HeapAlloc,
        HeapSys:       memStats.HeapSys,
        HeapIdle:      memStats.HeapIdle,
        HeapInuse:     memStats.HeapInuse,
        HeapReleased:  memStats.HeapReleased,
        HeapObjects:   memStats.HeapObjects,
        StackInuse:    memStats.StackInuse,
        StackSys:      memStats.StackSys,
        MSpanInuse:    memStats.MSpanInuse,
        MSpanSys:      memStats.MSpanSys,
        MCacheInuse:   memStats.MCacheInuse,
        MCacheSys:     memStats.MCacheSys,
        BuckHashSys:   memStats.BuckHashSys,
        GCSys:         memStats.GCSys,
        OtherSys:      memStats.OtherSys,
        NextGC:        memStats.NextGC,
        LastGC:        memStats.LastGC,
        PauseTotalNs:  memStats.PauseTotalNs,
        NumGC:         memStats.NumGC,
        NumForcedGC:   memStats.NumForcedGC,
        GCCPUFraction: memStats.GCCPUFraction,
    }
}

func DumpProfile(profile string, filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("could not create profile: %v", err)
    }
    defer f.Close()

    switch profile {
    case "heap":
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            return fmt.Errorf("could not write heap profile: %v", err)
        }
    case "goroutine":
        p := pprof.Lookup("goroutine")
        if p == nil {
            return fmt.Errorf("could not find goroutine profile")
        }
        if err := p.WriteTo(f, 0); err != nil {
            return fmt.Errorf("could not write goroutine profile: %v", err)
        }
    case "cpu":
        if err := pprof.StartCPUProfile(f); err != nil {
            return fmt.Errorf("could not start CPU profile: %v", err)
        }
        time.Sleep(time.Second * 10) // adjust duration as needed
        pprof.StopCPUProfile()
    case "threadcreate":
        p := pprof.Lookup("threadcreate")
        if p == nil {
            return fmt.Errorf("could not find threadcreate profile")
        }
        if err := p.WriteTo(f, 0); err != nil {
            return fmt.Errorf("could not write threadcreate profile: %v", err)
        }
    case "block":
        p := pprof.Lookup("block")
        if p == nil {
            return fmt.Errorf("could not find block profile")
        }
        if err := p.WriteTo(f, 0); err != nil {
            return fmt.Errorf("could not write block profile: %v", err)
        }
    default:
        return fmt.Errorf("unknown profile type: %v", profile)
    }
    return nil
}

func recordMetrics() {
    go func() {
        for {
            memStats := CaptureMemoryStats()
            memAlloc.Set(float64(memStats.Alloc))
            time.Sleep(10 * time.Second)
        }
    }()
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    if err := templates.ExecuteTemplate(w, "dashboard.html", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
    memStats := CaptureMemoryStats()
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(memStats); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, pass, ok := r.BasicAuth()
        if !ok || user != viper.GetString("username") || pass != viper.GetString("password") {
            w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    }
}

func main() {
    initConfig()
    validateConfig()

    interval := viper.GetDuration("interval")
    profile := viper.GetString("profile")
    profileFile := viper.GetString("profileFile")
    logFile := viper.GetString("logFile")

    logF, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("could not open log file: %v", err)
    }
    defer logF.Close()
    log.SetOutput(logF)
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    ticker := time.NewTicker(interval)
    done := make(chan bool)
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    http.Handle("/metrics", promhttp.Handler())
    http.HandleFunc("/dashboard", basicAuth(dashboardHandler))
    http.HandleFunc("/stats", basicAuth(statsHandler))

    go func() {
        log.Println("Starting HTTP server on :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatalf("HTTP server failed: %v", err)
        }
    }()

    recordMetrics()

    go func() {
        for {
            select {
            case <-done:
                return
            case <-ticker.C:
                memStats := CaptureMemoryStats()
                log.Printf("Memory Stats: %+v\n", memStats)

                if err := DumpProfile(profile, profileFile); err != nil {
                    log.Printf("Error dumping %s profile: %v\n", profile, err)
                }
            }
        }
    }()

    go func() {
        <-sigs
        log.Println("Received shutdown signal")
        ticker.Stop()
        done <- true
    }()

    // Wait for a shutdown signal
    select {}
}
