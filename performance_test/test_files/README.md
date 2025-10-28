 # Test Files Directory

This directory contains pre-generated test files used by k6 performance tests.

## Files

- `test_1mb.txt` - 1MB test file (1,048,576 bytes)
- `test_3mb.txt` - 3MB test file (3,145,728 bytes)
- `test_5mb.txt` - 5MB test file (5,242,880 bytes)

## Purpose

These files are used by the k6 load tests to simulate real file uploads without the overhead of generating files during test execution. This approach:

- **Reduces memory usage**: No need to generate files in memory for each virtual user
- **Improves performance**: Files are loaded once at the start of the test
- **More realistic**: Uses actual file data instead of generated strings
- **Faster test startup**: No generation delay before tests begin

## Regenerating Files

If you need to regenerate the test files, run:

```bash
cd performance_test/test_files

# Generate 1MB file
dd if=/dev/urandom of=test_1mb.txt bs=1024 count=1024

# Generate 3MB file
dd if=/dev/urandom of=test_3mb.txt bs=1024 count=3072

# Generate 5MB file
dd if=/dev/urandom of=test_5mb.txt bs=1024 count=5120
```

## Adding More Files

To add additional test files for different sizes:

1. Create the file using `dd` or any other method
2. Update the k6 test scripts to reference the new files
3. Document the new file in this README
