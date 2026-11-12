package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/term"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	// Number of pre-generated lines in the pool
	linePoolSize = 100
	// Write buffer size (256 KB)
	bufferSize = 256 * 1024
	// Test duration
	testDuration = 10 * time.Second
)

var (
	totalBytes uint64
	totalChars uint64
	// Available characters for random generation
	chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>? ")
	// ANSI color codes (8 basic colors + 8 bright colors)
	colors = []string{
		"\x1b[31m", "\x1b[32m", "\x1b[33m", "\x1b[34m", "\x1b[35m", "\x1b[36m", "\x1b[37m",
		"\x1b[91m", "\x1b[92m", "\x1b[93m", "\x1b[94m", "\x1b[95m", "\x1b[96m", "\x1b[97m",
	}
	reset = "\x1b[0m"
)

func main() {
	// Detect terminal width
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80 // Fallback if detection fails
	}

	fmt.Fprintf(os.Stderr, "Terminal: %d columns\n", width)
	fmt.Fprintf(os.Stderr, "Starting %v benchmark... (Ctrl+C to stop early)\n\n", testDuration)

	// Pre-generate a pool of random lines
	linePool := make([][]byte, linePoolSize)
	lineCharCounts := make([]int, linePoolSize)
	for i := range linePoolSize {
		linePool[i], lineCharCounts[i] = generateRandomLine(width)
	}

	startTime := time.Now()

	// Handle interrupt signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Buffered writer to optimize writes
	writer := bufio.NewWriterSize(os.Stdout, bufferSize)

	// Timer for test duration
	timer := time.NewTimer(testDuration)

	// Main write loop in a goroutine
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// Random selection of a line from the pool
				idx := rand.Intn(linePoolSize)
				line := linePool[idx]
				n, err := writer.Write(line)
				if err != nil {
					return
				}
				atomic.AddUint64(&totalBytes, uint64(n))
				atomic.AddUint64(&totalChars, uint64(lineCharCounts[idx]))

				// Periodic flush to avoid blocking the buffer indefinitely
				if atomic.LoadUint64(&totalBytes)%uint64(bufferSize*10) == 0 {
					writer.Flush()
				}
			}
		}
	}()

	// Wait for timer or stop signal
	select {
	case <-timer.C:
		fmt.Fprintln(os.Stderr, "\n\nTest duration reached...")
	case <-sigChan:
		fmt.Fprintln(os.Stderr, "\n\nStopping early...")
	}

	close(done)
	writer.Flush()

	// Display final results
	elapsed := time.Since(startTime).Seconds()
	finalBytes := atomic.LoadUint64(&totalBytes)
	finalChars := atomic.LoadUint64(&totalChars)

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Get terminal size again (in case it was resized)
	finalWidth, finalHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		finalWidth = width
		finalHeight = 0
	}

	blue := "\x1b[34m"
	resetColor := "\x1b[0m"

	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s========================================\n", blue)
	fmt.Fprintf(os.Stderr, "Final Results:\n")
	fmt.Fprintf(os.Stderr, "========================================\n")
	fmt.Fprintf(os.Stderr, "Hostname:     %s\n", hostname)
	fmt.Fprintf(os.Stderr, "Terminal:     %dx%d\n", finalWidth, finalHeight)
	fmt.Fprintf(os.Stderr, "Duration:     %.2f seconds\n", elapsed)
	fmt.Fprintf(os.Stderr, "Total chars:  %s (%.2f GChars)\n", formatNumber(finalChars), float64(finalChars)/1_000_000_000)
	fmt.Fprintf(os.Stderr, "Total bytes:  %s (%.2f MB)\n", formatNumber(finalBytes), float64(finalBytes)/1024/1024)
	fmt.Fprintf(os.Stderr, "Throughput:   %s chars/sec\n", formatNumber(uint64(float64(finalChars)/elapsed)))
	fmt.Fprintf(os.Stderr, "Throughput:   %.2f MB/s\n", float64(finalBytes)/elapsed/1024/1024)
	fmt.Fprintf(os.Stderr, "========================================%s\n", resetColor)
	time.Sleep(10 * time.Second)
}

// generateRandomLine creates a random line of the specified width with ANSI color codes
// Returns the line bytes and the count of visible characters (excluding ANSI codes)
func generateRandomLine(width int) ([]byte, int) {
	var line []byte
	currentLen := 0

	for currentLen < width {
		// Random color injection (about 10% chance)
		if rand.Float32() < 0.1 {
			color := colors[rand.Intn(len(colors))]
			line = append(line, []byte(color)...)
		}

		// Add a random character
		char := chars[rand.Intn(len(chars))]
		line = append(line, string(char)...)
		currentLen++

		// Random color reset (about 5% chance)
		if rand.Float32() < 0.05 {
			line = append(line, []byte(reset)...)
		}
	}

	// Final reset and newline
	line = append(line, []byte(reset)...)
	line = append(line, '\n')
	return line, width
}

// formatNumber formats a number with thousands separators
func formatNumber(n uint64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", n)
}
