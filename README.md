# Terminal Benchmark

A high-performance Go application to benchmark terminal output speed by displaying pseudo-random colored text.

## Overview

This tool measures how many characters per second your terminal emulator can handle by writing full-width lines of random text with ANSI color codes. It's useful for comparing terminal performance across different emulators (Alacritty, Kitty, Gnome Terminal, etc.) and configurations.

## Features

- **High-throughput output**: Uses buffered I/O (256KB buffer) to maximize write performance
- **Pre-generated line pool**: 100 random lines cached to minimize CPU overhead during the benchmark
- **ANSI colors**: Random color injection to stress-test terminal rendering
- **Real-time metrics**: Live display of current and average throughput
- **Comprehensive results**: Final report with hostname, terminal size, and performance statistics
- **Graceful shutdown**: Stop early with Ctrl+C while preserving statistics

## Installation

### Prerequisites

- Go 1.21 or later

### Build

```bash
git clone <repository-url>
cd terminal-bench
go build -o terminal-bench main.go
```

## Usage

Run the benchmark for 10 seconds (default duration):

```bash
./terminal-bench
```

Stop early with `Ctrl+C` to see results immediately.

### Example Output

```
Terminal: 196 columns
Starting 10s benchmark... (Ctrl+C to stop early)

Current: 43,174,899 chars/s | Total: 431,748,996 chars | Avg: 43,154,400 chars/s

Test duration reached...

========================================
Final Results:
========================================
Hostname:     i5
Terminal:     196x57
Duration:     10.00 seconds
Total chars:  431,748,996 (0.43 GChars)
Total bytes:  745,547,127 (711.01 MB)
Throughput:   43,154,400 chars/sec
Throughput:   71.07 MB/s
========================================
```

### Measuring Raw Performance

To measure the Go application's performance without terminal rendering overhead:

```bash
./terminal-bench > /dev/null
```

This redirects output to `/dev/null`, showing the maximum throughput your system can generate.

## Configuration

You can modify these constants in `main.go`:

- `testDuration`: Benchmark duration (default: 10 seconds)
- `linePoolSize`: Number of pre-generated lines (default: 100)
- `bufferSize`: Write buffer size (default: 256 KB)
- `reportInterval`: Metrics update frequency (default: 1 second)

## Performance Tips

### Factors Affecting Speed

1. **Terminal emulator**: GPU-accelerated terminals (Alacritty, Kitty) are generally faster
2. **Font size**: Smaller fonts = less pixels to render = higher throughput
3. **Window size**: More visible lines can affect scrolling performance
4. **Colors**: ANSI escape sequences add rendering overhead
5. **Hardware**: GPU, CPU, and display refresh rate all matter

### Recommended Terminals for Testing

- **Alacritty**: GPU-accelerated, very fast
- **Kitty**: GPU-accelerated with excellent features
- **WezTerm**: GPU-accelerated, cross-platform
- **Gnome Terminal**: CPU-based, widely available
- **iTerm2** (macOS): Popular with good performance

## Technical Details

### Architecture

- **Buffered I/O**: Uses `bufio.Writer` with 256KB buffer to minimize syscalls
- **Line-based output**: Generates full terminal-width lines followed by newlines for optimal scrolling
- **Atomic counters**: Thread-safe metrics using `sync/atomic`
- **Color codes**: Random injection of ANSI 8-color and bright-color codes
- **Character set**: Alphanumeric + symbols for realistic text rendering

### Output Format

- Progress metrics displayed on `stderr` (real-time updates)
- Benchmark text output on `stdout` (can be redirected)
- Final results formatted with thousand separators for readability
- Blue-colored summary for visibility

## Use Cases

- **Compare terminal emulators**: Which is fastest on your system?
- **Test font rendering**: Does font size affect performance?
- **Benchmark hardware**: How does your GPU/CPU handle terminal output?
- **CI/CD testing**: Measure terminal performance in automated environments
- **Educational**: Learn about terminal I/O optimization

## Contributing

Contributions are welcome! Feel free to:

- Report bugs or performance issues
- Suggest new features or metrics
- Submit pull requests with improvements

## License

MIT License - feel free to use and modify as needed.

## Acknowledgments

Built with Go's standard library and:
- `golang.org/x/term` for terminal size detection
- `golang.org/x/text` for number formatting
