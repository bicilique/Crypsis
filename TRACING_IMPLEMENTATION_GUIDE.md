# 🔍 OPENTELEMETRY TRACING IMPLEMENTATION GUIDE

## ✅ What We've Implemented

We've added **comprehensive, clean, and DRY** OpenTelemetry tracing to the Crypsis backend. Here's what's been done:

### 1. **Automatic Database Tracing** ✨
- **GORM Plugin** automatically traces ALL database operations
- **Zero code changes needed** in repositories
- Tracks: SELECT, INSERT, UPDATE, DELETE, queries

### 2. **Tracing Helper Utilities** 🛠️
- Centralized span creation
- Type-specific helpers (DB, KMS, Storage, HTTP, Crypto)
- Consistent error handling
- Easy-to-use API

### 3. **Service Layer Tracing** 📦
- Added tracing to KMS operations
- Added tracing to Storage operations  
- Template for adding to other services

---

## 🎯 Architecture Overview

```
HTTP Request (Auto-traced by middleware)
    ↓
Service Layer (Manual spans for business logic)
    ↓
Repository Layer (Auto-traced by GORM plugin)
    ↓
External Services (Manual spans: KMS, Storage, HTTP)
```

### Data Flow:
```
Application Code
    ↓ Creates Spans
OpenTelemetry SDK
    ↓ OTLP Protocol (HTTP)
OpenTelemetry Collector
    ↓ Exports
Jaeger (Distributed Tracing UI)
```

---

## 📚 How to Use Tracing in Your Code

### Pattern 1: Simple Span (Service Layer)

```go
func (s *SomeService) DoSomething(ctx context.Context, id string) error {
    // Get the tracing helper
    tracer := helper.GetTracingHelper()
    
    // Start a span - always defer End()
    ctx, span := tracer.StartSpan(ctx, "DoSomething")
    defer span.End()
    
    // Add custom attributes
    helper.AddAttributes(span, map[string]interface{}{
        "resource.id": id,
        "operation": "process",
    })
    
    // Your business logic here
    result, err := s.processData(ctx, id)
    if err != nil {
        helper.RecordError(span, err)  // Record error in span
        return err
    }
    
    helper.RecordSuccess(span, "Processing completed")
    return nil
}
```

### Pattern 2: Database Operations (Automatic!)

```go
// NO CHANGES NEEDED! GORM plugin automatically traces this:
func (r *fileRepository) GetByID(ctx context.Context, id string) (*entity.Files, error) {
    var file entity.Files
    // This query is automatically traced with span name: "DB: SELECT files"
    if err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &file, nil
}
```

**Key Point**: Always pass `ctx` to GORM methods using `.WithContext(ctx)`

### Pattern 3: KMS Operations

```go
func (s *KmsService) EncryptData(ctx context.Context, keyID, data string) (string, error) {
    // Use the KMS-specific span helper
    tracer := helper.GetTracingHelper()
    ctx, span := tracer.StartKMSSpan(ctx, "EncryptData", keyID)
    defer span.End()
    
    // Generate request
    jsonBody, err := helper.GenerateEncryptTemplate(keyID, data)
    if err != nil {
        helper.RecordError(span, err)
        return "", err
    }
    
    // Send request
    response, err := s.sendRequest(ctx, jsonBody)
    if err != nil {
        helper.RecordError(span, err)
        return "", err
    }
    
    helper.RecordSuccess(span, "Data encrypted")
    return response, nil
}
```

### Pattern 4: Storage Operations

```go
func (s *MinioService) DownloadFile(ctx context.Context, bucketName, fileName string) ([]byte, error) {
    // Use the storage-specific span helper
    tracer := helper.GetTracingHelper()
    ctx, span := tracer.StartStorageSpan(ctx, "GetObject", bucketName, fileName)
    defer span.End()
    
    object, err := s.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    defer object.Close()
    
    // Read data
    buf := new(bytes.Buffer)
    if _, err := io.Copy(buf, object); err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    helper.RecordSuccess(span, "File downloaded")
    return buf.Bytes(), nil
}
```

### Pattern 5: HTTP Client Calls

```go
func (s *HydraService) IntrospectToken(ctx context.Context, token string) (*TokenInfo, error) {
    tracer := helper.GetTracingHelper()
    ctx, span := tracer.StartHTTPClientSpan(ctx, "POST", s.hydraURL+"/oauth2/introspect")
    defer span.End()
    
    req, err := http.NewRequestWithContext(ctx, "POST", url, body)
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    resp, err := s.client.Do(req)
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        err := fmt.Errorf("unexpected status: %d", resp.StatusCode)
        helper.RecordError(span, err)
        return nil, err
    }
    
    helper.RecordSuccess(span, "Token introspected")
    return parseResponse(resp.Body)
}
```

### Pattern 6: Cryptographic Operations

```go
func (s *CryptoService) EncryptFile(ctx context.Context, plaintext []byte, algorithm string) ([]byte, error) {
    tracer := helper.GetTracingHelper()
    ctx, span := tracer.StartCryptoSpan(ctx, "Encrypt", algorithm)
    defer span.End()
    
    helper.AddAttributes(span, map[string]interface{}{
        "plaintext.size": len(plaintext),
    })
    
    ciphertext, err := s.aead.Encrypt(plaintext, nil)
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    helper.AddAttributes(span, map[string]interface{}{
        "ciphertext.size": len(ciphertext),
    })
    
    helper.RecordSuccess(span, "Encryption completed")
    return ciphertext, nil
}
```

---

## 🎨 Best Practices

### ✅ DO:

1. **Always pass context**
   ```go
   ctx, span := tracer.StartSpan(ctx, "operation")
   defer span.End()
   // Pass ctx to all downstream calls
   result, err := s.repository.GetData(ctx, id)
   ```

2. **Always defer span.End()**
   ```go
   ctx, span := tracer.StartSpan(ctx, "operation")
   defer span.End()  // ← IMPORTANT!
   ```

3. **Record errors**
   ```go
   if err != nil {
       helper.RecordError(span, err)
       return err
   }
   ```

4. **Add meaningful attributes**
   ```go
   helper.AddAttributes(span, map[string]interface{}{
       "user.id": userID,
       "file.size": fileSize,
       "operation.type": "upload",
   })
   ```

5. **Use descriptive span names**
   ```go
   // ✅ Good
   ctx, span := tracer.StartSpan(ctx, "EncryptFileWithKMS")
   
   // ❌ Bad
   ctx, span := tracer.StartSpan(ctx, "Process")
   ```

### ❌ DON'T:

1. **Don't forget to pass context**
   ```go
   // ❌ Bad - context not passed
   r.db.First(&file, "id = ?", id)
   
   // ✅ Good
   r.db.WithContext(ctx).First(&file, "id = ?", id)
   ```

2. **Don't create spans without ending them**
   ```go
   // ❌ Bad - memory leak!
   ctx, span := tracer.StartSpan(ctx, "operation")
   // forgot defer span.End()
   ```

3. **Don't ignore errors in spans**
   ```go
   // ❌ Bad
   if err != nil {
       return err  // Error not recorded in span
   }
   
   // ✅ Good
   if err != nil {
       helper.RecordError(span, err)
       return err
   }
   ```

---

## 📊 What You'll See in Jaeger

After implementing tracing, you'll see traces like this:

```
POST /api/files (HTTP Request - auto-traced)
├── FileService.UploadFile (service span)
│   ├── Crypto: Encrypt (crypto span)
│   ├── KMS: EncryptKey (KMS span)
│   ├── Storage: PutObject (storage span)
│   ├── DB: INSERT files (auto-traced by GORM)
│   └── DB: INSERT file_logs (auto-traced by GORM)
└── HTTP Response
```

Each span shows:
- **Duration** - How long it took
- **Attributes** - Custom data (IDs, sizes, etc.)
- **Events** - Logs/errors
- **Status** - OK or Error

---

## 🚀 Quick Start Checklist

To add tracing to a new service method:

1. [ ] Get the tracing helper: `tracer := helper.GetTracingHelper()`
2. [ ] Start span: `ctx, span := tracer.StartXXXSpan(ctx, "operation", ...)`
3. [ ] Defer end: `defer span.End()`
4. [ ] Add attributes: `helper.AddAttributes(span, {...})`
5. [ ] Record errors: `helper.RecordError(span, err)`
6. [ ] Record success: `helper.RecordSuccess(span, "message")`
7. [ ] Pass context to all calls: `someFunc(ctx, ...)`

---

## 🔧 Available Helper Methods

### Tracing Helper Methods:

| Method | Use Case | Example |
|--------|----------|---------|
| `StartSpan()` | Generic operations | `StartSpan(ctx, "ProcessData")` |
| `StartDBSpan()` | Database queries | `StartDBSpan(ctx, "SELECT", "files")` |
| `StartKMSSpan()` | KMS operations | `StartKMSSpan(ctx, "Encrypt", keyID)` |
| `StartStorageSpan()` | Storage ops | `StartStorageSpan(ctx, "Put", bucket, object)` |
| `StartHTTPClientSpan()` | HTTP calls | `StartHTTPClientSpan(ctx, "GET", url)` |
| `StartCryptoSpan()` | Crypto ops | `StartCryptoSpan(ctx, "Encrypt", "AES-256")` |

### Utility Methods:

| Method | Purpose |
|--------|---------|
| `RecordError(span, err)` | Mark span as error and record error details |
| `RecordSuccess(span, msg)` | Mark span as successful with message |
| `AddAttributes(span, map)` | Add multiple custom attributes at once |

---

## 📝 Example: Full Service Method

Here's a complete example showing all best practices:

```go
func (s *FileService) EncryptAndUpload(ctx context.Context, req EncryptRequest) (*FileResponse, error) {
    // 1. Start main span
    tracer := helper.GetTracingHelper()
    ctx, span := tracer.StartSpan(ctx, "FileService.EncryptAndUpload")
    defer span.End()
    
    // 2. Add request attributes
    helper.AddAttributes(span, map[string]interface{}{
        "file.id":   req.FileID,
        "file.size": req.Size,
        "app.id":    req.AppID,
    })
    
    // 3. Validate input (no span needed for quick validation)
    if req.FileID == "" {
        err := model.ErrInvalidInput
        helper.RecordError(span, err)
        return nil, err
    }
    
    // 4. Database operation (auto-traced by GORM plugin)
    file, err := s.fileRepository.GetByID(ctx, req.FileID)
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    // 5. Crypto operation (creates child span)
    ciphertext, err := s.cryptoService.Encrypt(ctx, file.Data)
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    // 6. KMS operation (creates child span)
    encryptedKey, err := s.kmsService.EncryptKey(ctx, s.keyConfig.UID, file.DEK)
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    // 7. Storage operation (creates child span)
    uploadResp, err := s.storageService.UploadFile(ctx, s.bucketName, file.ID, ciphertext, len(ciphertext))
    if err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    // 8. Database update (auto-traced)
    file.EncryptedDEK = encryptedKey
    file.StorageLocation = uploadResp.Location
    if err := s.fileRepository.Update(ctx, file); err != nil {
        helper.RecordError(span, err)
        return nil, err
    }
    
    // 9. Success!
    helper.RecordSuccess(span, "File encrypted and uploaded successfully")
    
    return &FileResponse{
        ID:       file.ID,
        Location: uploadResp.Location,
    }, nil
}
```

---

## 🎯 Summary

### What's Automatic:
- ✅ **All HTTP requests** (middleware)
- ✅ **All database queries** (GORM plugin)
- ✅ **Go runtime metrics** (memory, GC, goroutines)

### What Needs Manual Spans:
- 🔧 **Service layer business logic**
- 🔧 **KMS operations**
- 🔧 **Storage operations**
- 🔧 **HTTP client calls**
- 🔧 **Cryptographic operations**

### Benefits:
- 🎯 **Complete visibility** into request flow
- 🐛 **Easy debugging** - see exactly where time is spent
- 📊 **Performance insights** - identify bottlenecks
- 🔍 **Error tracking** - see full error context
- 📈 **Distributed tracing** - track requests across services

---

## 🆘 Need Help?

1. **Check Jaeger**: http://localhost:16686
2. **Run a request** and look for your service name
3. **Click on a trace** to see the span hierarchy
4. **Add more spans** where you need visibility

Happy tracing! 🎉
