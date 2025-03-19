
## Memory Dump Analyzer

A robust memory dump analyzer for Go applications with advanced profiling, real-time monitoring, and anomaly detection. This tool is designed to help developers and system administrators monitor and analyze memory usage in Go applications, providing valuable insights and facilitating performance tuning.

### Features

- **Dynamic Configuration Reloading**: Automatically reloads configuration changes without restarting the application.
- **Web-Based UI**: Provides a user-friendly dashboard for real-time monitoring of memory statistics.
- **Advanced Profiling Options**: Supports multiple profiling types, including heap, goroutine, CPU, thread creation, and block profiles.
- **Prometheus Integration**: Exposes metrics for Prometheus, enabling detailed monitoring and alerting.
- **Basic Authentication**: Secure access to the dashboard and profiling data.
- **Comprehensive Logging**: Logs detailed memory statistics and profiling information.
- **Automated Reporting**: Generates detailed reports periodically or after significant events.

### Installation

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/yourusername/memdumpanalyzer.git
   cd memdumpanalyzer
   ```

2. **Create Configuration File**:
   Create a `config.yaml` file in the project directory with the following content:
   ```yaml
   interval: 10s
   duration: 1m
   profile: heap
   profileFile: memprofile.prof
   logFile: analyzer.log
   ```

3. **Build the Application**:
   ```sh
   go build -o memdumpanalyzer
   ```

### Usage

1. **Run the Application**:
   ```sh
   ./memdumpanalyzer
   ```

2. **Access the Dashboard**:
   Open your browser and go to `http://localhost:8080/dashboard`. Use `admin` as the username and `password` as the password for basic authentication.

3. **Monitor Metrics**:
   Open your browser and go to `http://localhost:2112/metrics` to view Prometheus metrics.

### Configuration

The application uses a configuration file (`config.yaml`) for setting various parameters. Below is an example configuration file:

```yaml
interval: 10s             # Interval between memory stats collection
duration: 1m              # Total duration to run the analyzer
profile: heap             # Type of profile to capture (heap, goroutine, cpu, threadcreate, block)
profileFile: memprofile.prof  # File to write memory profile
logFile: analyzer.log     # File to write logs
```

### Dashboard

The web-based dashboard provides a real-time view of memory statistics, including metrics such as allocated memory, total allocations, heap usage, garbage collection details, and more. The dashboard updates every 5 seconds to reflect the latest memory statistics.

### Prometheus Integration

Metrics are exposed at `http://localhost:2112/metrics`, allowing for integration with Prometheus. This enables detailed monitoring, alerting, and visualization using tools like Grafana.

### Security

Basic authentication is implemented to secure access to the dashboard and profiling data. The default username is `admin` and the default password is `password`. These can be modified in the source code or by implementing a more secure authentication mechanism.

### Explanation of Metrics
- Alloc: Bytes of allocated heap objects (live objects).
- TotalAlloc: Cumulative bytes allocated for heap objects.
- sys: Total bytes of memory obtained from the OS.
- Lookups: Number of pointer lookups performed by the runtime.
- Mallocs: Total number of heap object allocations.
- Frees: Total number of heap object deallocations.
- HeapAlloc: Bytes of allocated heap objects (same as Alloc).
- HeapSys: Total bytes of heap memory obtained from the OS.
- HeapIdle: Bytes of heap memory that are not currently in use.
- HeapInuse: Bytes of heap memory that are in use.
- HeapReleased: Bytes of heap memory returned to the OS.
- HeapObjects: Number of allocated heap objects.
- StackInuse: Bytes of stack memory in use.
- StackSys: Total bytes of stack memory obtained from the OS.
- MSpanInuse: Bytes of memory in use by mspan structures.
- MSpanSys: Total bytes of memory obtained from the OS for mspan - structures.
- MCacheInuse: Bytes of memory in use by mcache structures.
- MCacheSys: Total bytes of memory obtained from the OS for mcache structures.
- BuckHashSys: Bytes of memory in use by the runtime's hash tables.
- GCSys: Bytes of memory in use by the garbage collector.
- OtherSys: Bytes of memory in use by other system allocations.
- NextGC: Target heap size for the next garbage collection.
- LastGC: Time of the last garbage collection in nanoseconds since the epoch.
- PauseTotalNs: Cumulative nanoseconds in GC stop-the-world pauses since the program started.
- NumGC: Number of completed garbage collection cycles.
- NumForcedGC: Number of GC cycles forced by the application.
- GCCPUFraction: Fraction of CPU time spent in garbage collection.