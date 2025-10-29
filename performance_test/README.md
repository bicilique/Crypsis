# 📊 Crypsis Observability Stack

Complete observability solution for monitoring, tracing, and analyzing Crypsis application performance.

## 🎯 What's Included

- **✅ OpenTelemetry**: Distributed tracing and metrics collection
- **✅ Prometheus**: Time-series metrics database
- **✅ Jaeger**: Distributed tracing visualization
- **✅ Grafana**: Unified dashboards and visualization
- **✅ Pre-configured Dashboards**: Ready-to-use performance dashboards
- **✅ Automatic Instrumentation**: Built-in middleware for HTTP tracing

## 🚀 Quick Start (30 seconds)

### Option 1: Use the Startup Script (Recommended)

```bash
./start-observability.sh
```

This script will:
1. Start all observability services
2. Verify they're healthy
3. Display access URLs and credentials
4. Show you next steps

### Option 2: Manual Start

```bash
docker-compose up -d
```

Then open Grafana at http://localhost:3000 (admin/admin)

## 📱 Access Your Dashboards

| Service | URL | Purpose |
|---------|-----|---------|
| **🎨 Grafana** | http://localhost:3000 | Main visualization platform |
| **📈 Prometheus** | http://localhost:9090 | Metrics query engine |
| **🔍 Jaeger** | http://localhost:16686 | Trace visualization |

**Grafana Login**:
- Username: `admin`
- Password: `admin`

## 📊 Available Dashboards

### 1. Crypsis Application Performance
**What it shows**:
- ⚡ Request rate and throughput
- ⏱️ Response time (p50, p90, p95, p99)
- ❌ Error rate
- 💾 Memory usage
- 🔐 Cryptographic operation performance
- 📁 File processing metrics

### 2. OpenTelemetry Collector Metrics
**What it shows**:
- 📥 Traces and metrics received
- 📤 Export status
- 🚨 Errors and failures
- 💻 Collector resource usage

## 🔍 How to Use

### Monitor Real-Time Performance

1. Open Grafana: http://localhost:3000
2. Go to **Dashboards** → **Browse**
3. Select **Crypsis Application Performance**
4. Watch real-time metrics as your app handles requests

### Investigate Slow Requests

1. In Grafana, check the **Response Time Latency** panel
2. Note the time period with high latency
3. Open Jaeger: http://localhost:16686
4. Select service: `crypsis-backend`
5. Search for traces in that time period
6. Click on slow traces to see detailed breakdown
7. Identify which operation is causing the delay

### Find Memory Leaks

1. In Grafana, check the **Memory Usage** panel
2. Look for continuous upward trend (not saw-tooth pattern)
3. Check **Active Goroutines** panel
4. If both are increasing → likely goroutine leak
5. Use traces to identify which endpoints are involved

## 🧪 Performance Testing

### Run Load Test

```bash
cd performance_test
k6 run scripts/k6_load_test.js
```

**While test is running**:
1. Open Grafana dashboard
2. Watch metrics in real-time
3. Take screenshots at peak load
4. Note any errors or anomalies

**After test**:
- Review k6 results in `performance_test/results/`
- Export Grafana dashboard screenshots
- Analyze slow traces in Jaeger

## 📚 Documentation

- **[Complete Guide](OBSERVABILITY_COMPLETE_GUIDE.md)**: Detailed architecture and usage
- **[Quick Reference](OBSERVABILITY_QUICK_REFERENCE.md)**: Commands and troubleshooting
- **[Architecture](OBSERVABILITY_ARCHITECTURE.md)**: System design and components (if exists)

## 🛠️ Common Tasks

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f otel-collector
docker-compose logs -f grafana
docker-compose logs -f prometheus
docker-compose logs -f jaeger
```

### Restart Observability Stack

```bash
docker-compose restart otel-collector prometheus jaeger grafana
```

### Stop Everything

```bash
docker-compose down
```

### Clean Start (Remove All Data)

```bash
docker-compose down -v
docker-compose up -d
```

## ✅ Health Check

### Quick Verification

```bash
# Check all containers are running
docker-compose ps

# Should show all services as "Up"
```

### Detailed Verification

1. **Prometheus Targets**: http://localhost:9090/targets
   - ✅ All targets should be **UP** (green)

2. **Grafana Dashboards**: http://localhost:3000
   - ✅ Should see 2 dashboards

3. **OTEL Collector**:
   ```bash
   docker-compose logs otel-collector | tail -20
   ```
   - ✅ Should see traces/metrics being received

4. **Application Telemetry**:
   ```bash
   docker-compose logs crypsis-backend | grep -i "opentelemetry"
   ```
   - ✅ Should see "OpenTelemetry initialized successfully"

## 🎯 Key Metrics to Monitor

| Metric | Good | Warning | Critical |
|--------|------|---------|----------|
| Response Time (p95) | < 2s | 2-5s | > 10s |
| Error Rate | < 0.1% | 0.1-1% | > 5% |
| Memory Usage | Stable | Growing slowly | Continuously growing |
| Goroutines | Stable | Fluctuating | Continuously increasing |

## 🔧 Troubleshooting

### No Data in Grafana?

```bash
# 1. Restart observability stack
docker-compose restart otel-collector prometheus grafana

# 2. Check Prometheus targets
# Open: http://localhost:9090/targets
# All should be UP

# 3. Verify app is sending telemetry
docker-compose logs crypsis-backend | grep -i "otel"
```

### OTEL Collector Not Receiving Data?

```bash
# Check collector logs
docker-compose logs otel-collector | tail -50

# Verify app has OTEL enabled
docker-compose exec crypsis-backend env | grep OTEL
```

### Dashboard Shows "No Data"?

1. Check time range (top-right in Grafana)
2. Verify Prometheus data source is configured
3. Generate some traffic to the application
4. Wait 10-15 seconds for metrics to appear

## 📊 Architecture Overview

```
Application (Go) → OpenTelemetry SDK → OTLP (4318/4317)
                                            ↓
                                   OTEL Collector
                                    ↙          ↘
                               Jaeger      Prometheus
                                    ↘          ↙
                                     Grafana
```

**Flow**:
1. App generates traces and metrics (OpenTelemetry SDK)
2. Sends to OTEL Collector via OTLP protocol
3. Collector processes and exports to Jaeger (traces) and Prometheus (metrics)
4. Grafana queries both for unified visualization

## 🎓 Learn More

### Prometheus Query Examples

```promql
# Request rate
rate(http_requests_total[5m])

# 95th percentile response time
histogram_quantile(0.95, sum(rate(http_request_duration_milliseconds_bucket[5m])) by (le))

# Error rate
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

### Jaeger Search Tips

- **Service**: Select `crypsis-backend`
- **Operation**: Choose specific endpoint (e.g., POST /api/files/encrypt)
- **Tags**: Filter by `error=true` to find failed requests
- **Sort by**: Duration (longest first) to find slow requests

## 🔐 Security Notes

- Grafana uses default credentials (change in production!)
- Prometheus has no authentication (internal Docker network only)
- Jaeger has no authentication (internal Docker network only)
- OTEL Collector uses insecure gRPC (internal Docker network only)

## 📞 Support

For issues:
1. Check logs: `docker-compose logs [service]`
2. Review **Quick Reference**: `OBSERVABILITY_QUICK_REFERENCE.md`
3. Read **Complete Guide**: `OBSERVABILITY_COMPLETE_GUIDE.md`

## 🎉 Success Indicators

You'll know everything is working when:

✅ All Docker containers are running  
✅ Prometheus targets are UP  
✅ Grafana dashboards show data  
✅ Jaeger shows traces  
✅ OTEL Collector logs show received spans/metrics  

---

**Ready to start?** Run `./start-observability.sh` and open http://localhost:3000 🚀
