package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ─────────────────────────── ANSI Color Palette ───────────────────────────

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
)

// ─────────────────────────── Helpers ──────────────────────────────────────

func colorize(color, text string) string {
	return color + text + Reset
}

func printBanner() {
	fmt.Println()
	fmt.Println(colorize(Bold+Cyan, "  ███████╗ █████╗ ███████╗████████╗███████╗██████╗ "))
	fmt.Println(colorize(Bold+Cyan, "  ██╔════╝██╔══██╗██╔════╝╚══██╔══╝██╔════╝██╔══██╗"))
	fmt.Println(colorize(Bold+Cyan, "  █████╗  ███████║███████╗   ██║   █████╗  ██████╔╝"))
	fmt.Println(colorize(Bold+Cyan, "  ██╔══╝  ██╔══██║╚════██║   ██║   ██╔══╝  ██╔══██╗"))
	fmt.Println(colorize(Bold+Cyan, "  ██║     ██║  ██║███████║   ██║   ███████╗██║  ██║"))
	fmt.Println(colorize(Bold+Cyan, "  ╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝"))
	fmt.Println()
	fmt.Println(colorize(Dim+White, "  ⚡ High-Speed Video & Audio Downloader for Termux"))
	fmt.Println(colorize(Dim+White, "  ─────────────────────────────────────────────────"))
	fmt.Println()
}

func printDivider() {
	fmt.Println(colorize(Dim, "  ──────────────────────────────────────────────────────"))
}

func printSuccess(msg string) {
	fmt.Println(colorize(Bold+Green, "  ✓ "+msg))
}

func printError(msg string) {
	fmt.Println(colorize(Bold+Red, "  ✗ Error: "+msg))
}

func printInfo(msg string) {
	fmt.Println(colorize(Cyan, "  ➜ "+msg))
}

func printWarning(msg string) {
	fmt.Println(colorize(Yellow, "  ⚠ "+msg))
}

func printPrompt(label string) {
	fmt.Print(colorize(Bold+Magenta, "  ❯ "+label+": "))
}

func readInput(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// ─────────────────────────── Validation ──────────────────────────────────

func isValidURL(u string) bool {
	urlPattern := regexp.MustCompile(`^(https?://)[^\s/$.?#].[^\s]*$`)
	return urlPattern.MatchString(u)
}

func checkDependency(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func checkDependencies() bool {
	deps := []string{"yt-dlp", "ffmpeg"}
	allOk := true
	for _, dep := range deps {
		if checkDependency(dep) {
			printSuccess(dep + " found")
		} else {
			printError(dep + " not found — install it with: pip install " + dep)
			allOk = false
		}
	}
	return allOk
}

// ─────────────────────────── Spinner ──────────────────────────────────────

func runSpinner(msg string, done chan bool) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r  %s %s\n", colorize(Green, "✓"), msg)
			return
		default:
			fmt.Printf("\r  %s %s", colorize(Cyan, frames[i%len(frames)]), msg)
			time.Sleep(80 * time.Millisecond)
			i++
		}
	}
}

// ─────────────────────────── Format Listing ───────────────────────────────

func listFormats(url string) (string, error) {
	done := make(chan bool)
	go runSpinner("Fetching available formats…", done)

	cmd := exec.Command("yt-dlp", "--list-formats", "--no-warnings", "--no-playlist", url)
	out, err := cmd.CombinedOutput()
	done <- true
	time.Sleep(100 * time.Millisecond)

	if err != nil {
		return "", fmt.Errorf("yt-dlp failed: %s", strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func printFormats(raw string) {
	lines := strings.Split(raw, "\n")
	fmt.Println()
	printDivider()
	fmt.Println(colorize(Bold+Yellow, "  📋 Available Formats"))
	printDivider()

	headerPrinted := false
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		// Print the header line in bold white
		if strings.HasPrefix(line, "ID") || strings.HasPrefix(line, "[info]") {
			if !headerPrinted && strings.HasPrefix(line, "ID") {
				fmt.Println(colorize(Bold+White, "  "+line))
				headerPrinted = true
			}
			continue
		}
		// Color video lines blue, audio lines green
		lower := strings.ToLower(line)
		if strings.Contains(lower, "audio only") || strings.Contains(lower, "m4a") || strings.Contains(lower, "opus") || strings.Contains(lower, "mp3") {
			fmt.Println(colorize(Green, "  "+line))
		} else if strings.Contains(lower, "video only") || strings.Contains(lower, "mp4") || strings.Contains(lower, "webm") {
			fmt.Println(colorize(Blue, "  "+line))
		} else {
			fmt.Println(colorize(White, "  "+line))
		}
	}
	printDivider()
	fmt.Println()
	fmt.Println(colorize(Dim, "  "+colorize(Green, "■")+" Audio-only formats    "+colorize(Blue, "■")+" Video formats"))
	fmt.Println()
}

// ─────────────────────────── Format Type Menu ─────────────────────────────

func selectFormatType(reader *bufio.Reader) string {
	fmt.Println(colorize(Bold+White, "  ┌─────────────────────────────────────┐"))
	fmt.Println(colorize(Bold+White, "  │       What do you want to save?     │"))
	fmt.Println(colorize(Bold+White, "  ├─────────────────────────────────────┤"))
	fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Blue, "[1]") + colorize(White, " 🎬  Video  (with audio)         ") + colorize(Bold+White, "│"))
	fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Green, "[2]") + colorize(White, " 🎵  Audio  (music/podcast)      ") + colorize(Bold+White, "│"))
	fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Cyan, "[3]") + colorize(White, " 🎞   Custom (enter format ID)   ") + colorize(Bold+White, "│"))
	fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Red, "[q]") + colorize(White, " ✕   Quit                        ") + colorize(Bold+White, "│"))
	fmt.Println(colorize(Bold+White, "  └─────────────────────────────────────┘"))
	fmt.Println()

	for {
		printPrompt("Choice (1/2/3/q)")
		choice := strings.ToLower(readInput(reader))
		switch choice {
		case "1":
			return "video"
		case "2":
			return "audio"
		case "3":
			return "custom"
		case "q", "quit", "exit":
			return "quit"
		default:
			printWarning("Please enter 1, 2, 3, or q")
		}
	}
}

// ─────────────────────────── Quality Selection ────────────────────────────

func selectQuality(reader *bufio.Reader, formatType string) (string, string) {
	fmt.Println()

	switch formatType {
	case "video":
		fmt.Println(colorize(Bold+White, "  ┌─────────────────────────────────────┐"))
		fmt.Println(colorize(Bold+White, "  │         🎬 Video Quality              │"))
		fmt.Println(colorize(Bold+White, "  ├─────────────────────────────────────┤"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Cyan, "[1]") + colorize(White, " 4K  / 2160p  (best)            ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Cyan, "[2]") + colorize(White, " 1440p                          ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Cyan, "[3]") + colorize(White, " 1080p  Full HD                 ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Cyan, "[4]") + colorize(White, " 720p   HD                      ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Cyan, "[5]") + colorize(White, " 480p                           ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Cyan, "[6]") + colorize(White, " 360p   (small)                 ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Yellow, "[7]") + colorize(White, " Best available (auto)          ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Magenta, "[c]") + colorize(White, " Custom format ID               ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  └─────────────────────────────────────┘"))
		fmt.Println()

		printPrompt("Choice")
		choice := readInput(reader)
		qualityMap := map[string]string{
			"1": "bestvideo[height<=2160]+bestaudio/best[height<=2160]",
			"2": "bestvideo[height<=1440]+bestaudio/best[height<=1440]",
			"3": "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
			"4": "bestvideo[height<=720]+bestaudio/best[height<=720]",
			"5": "bestvideo[height<=480]+bestaudio/best[height<=480]",
			"6": "bestvideo[height<=360]+bestaudio/best[height<=360]",
			"7": "bestvideo+bestaudio/best",
		}
		if format, ok := qualityMap[choice]; ok {
			return format, "mp4"
		}
		if choice == "c" {
			printPrompt("Enter format ID from the list above")
			id := readInput(reader)
			return id, "mp4"
		}
		printWarning("Invalid choice, using best available")
		return "bestvideo+bestaudio/best", "mp4"

	case "audio":
		fmt.Println(colorize(Bold+White, "  ┌─────────────────────────────────────┐"))
		fmt.Println(colorize(Bold+White, "  │         🎵 Audio Format               │"))
		fmt.Println(colorize(Bold+White, "  ├─────────────────────────────────────┤"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Green, "[1]") + colorize(White, " MP3   320kbps  (recommended)   ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Green, "[2]") + colorize(White, " MP3   128kbps  (smaller)       ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Green, "[3]") + colorize(White, " M4A   (iTunes compatible)      ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Green, "[4]") + colorize(White, " OPUS  (best compression)       ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Green, "[5]") + colorize(White, " WAV   (lossless)               ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  │  ") + colorize(Bold+Yellow, "[6]") + colorize(White, " Best available (auto)          ") + colorize(Bold+White, "│"))
		fmt.Println(colorize(Bold+White, "  └─────────────────────────────────────┘"))
		fmt.Println()

		printPrompt("Choice")
		choice := readInput(reader)
		type audioOpt struct {
			format string
			ext    string
		}
		qualityMap := map[string]audioOpt{
			"1": {"bestaudio", "mp3"},
			"2": {"bestaudio", "mp3"},
			"3": {"bestaudio", "m4a"},
			"4": {"bestaudio", "opus"},
			"5": {"bestaudio", "wav"},
			"6": {"bestaudio", "best"},
		}
		bitrateMap := map[string]string{"1": "320", "2": "128"}
		if opt, ok := qualityMap[choice]; ok {
			bitrate := bitrateMap[choice]
			return opt.format + "|" + bitrate, opt.ext
		}
		printWarning("Invalid choice, using best MP3")
		return "bestaudio|320", "mp3"

	case "custom":
		printPrompt("Enter format ID (from the list above)")
		id := readInput(reader)
		printPrompt("Output extension (e.g. mp4, mkv, mp3)")
		ext := readInput(reader)
		if ext == "" {
			ext = "mp4"
		}
		return id, ext
	}
	return "bestvideo+bestaudio/best", "mp4"
}

// ─────────────────────────── Download ─────────────────────────────────────

func getDownloadDir() string {
	// Prefer shared storage on Android
	shared := "/sdcard/Download/Faster"
	if _, err := os.Stat("/sdcard"); err == nil {
		os.MkdirAll(shared, 0755)
		return shared
	}
	// Fallback to home
	home, _ := os.UserHomeDir()
	dir := home + "/Downloads/Faster"
	os.MkdirAll(dir, 0755)
	return dir
}

func buildYtdlpArgs(url, format, ext string) []string {
	outDir := getDownloadDir()
	outTemplate := outDir + "/%(title)s.%(ext)s"

	args := []string{
		"--no-playlist",
		"--no-warnings",
		"--progress",
		"--console-title",
		"--merge-output-format", "mp4",
		"-o", outTemplate,
	}

	// Audio download
	if strings.Contains(format, "bestaudio") && (ext == "mp3" || ext == "m4a" || ext == "opus" || ext == "wav") {
		parts := strings.SplitN(format, "|", 2)
		bitrate := "192"
		if len(parts) == 2 && parts[1] != "" {
			bitrate = parts[1]
		}
		args = append(args,
			"-x",
			"--audio-format", ext,
			"--audio-quality", bitrate+"k",
			"-f", parts[0],
		)
	} else {
		// Video or custom
		args = append(args, "-f", format)
		if ext != "mp4" {
			args = append(args, "--merge-output-format", ext)
		}
	}

	// Add cookies & retries for robustness
	args = append(args,
		"--retries", "5",
		"--fragment-retries", "5",
		"--socket-timeout", "30",
		url,
	)
	return args
}

func download(url, format, ext string) error {
	args := buildYtdlpArgs(url, format, ext)

	fmt.Println()
	printDivider()
	fmt.Println(colorize(Bold+Yellow, "  ⬇  Downloading…"))
	printDivider()
	fmt.Println(colorize(Dim, "  Save directory: "+getDownloadDir()))
	fmt.Println()

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ─────────────────────────── Post-Download Info ───────────────────────────

func printCompletion(url string) {
	fmt.Println()
	printDivider()
	printSuccess("Download complete!")
	printInfo("Saved to: " + getDownloadDir())
	printDivider()
	fmt.Println()
}

// ─────────────────────────── Another Download? ────────────────────────────

func askAnother(reader *bufio.Reader) bool {
	printPrompt("Download another? (y/n)")
	ans := strings.ToLower(readInput(reader))
	return ans == "y" || ans == "yes"
}

// ─────────────────────────── Main Flow ────────────────────────────────────

func run(reader *bufio.Reader, initialURL string) {
	url := initialURL

	for {
		if url == "" {
			fmt.Println()
			printDivider()
			printInfo("Paste your YouTube / Instagram / Facebook / TikTok URL")
			printInfo("Type " + colorize(Bold+Red, "q") + " to quit")
			printDivider()
			fmt.Println()
			printPrompt("URL")
			url = readInput(reader)
		}

		if strings.ToLower(url) == "q" || strings.ToLower(url) == "quit" {
			fmt.Println()
			printInfo("Goodbye! 👋")
			fmt.Println()
			os.Exit(0)
		}

		if url == "" {
			printWarning("URL cannot be empty")
			url = ""
			continue
		}

		if !isValidURL(url) {
			printError("That doesn't look like a valid URL. Please try again.")
			url = ""
			continue
		}

		// Fetch & display formats
		raw, err := listFormats(url)
		if err != nil {
			printError(err.Error())
			printWarning("Check the URL or your internet connection and try again.")
			url = ""
			continue
		}
		printFormats(raw)

		// Format type
		formatType := selectFormatType(reader)
		if formatType == "quit" {
			fmt.Println()
			printInfo("Goodbye! 👋")
			fmt.Println()
			os.Exit(0)
		}

		// Quality
		format, ext := selectQuality(reader, formatType)

		// Download
		err = download(url, format, ext)
		if err != nil {
			// yt-dlp writes its own error; just summarise
			printError("Download failed. Check the messages above for details.")
		} else {
			printCompletion(url)
		}

		// Loop?
		if !askAnother(reader) {
			fmt.Println()
			printInfo("Goodbye! 👋")
			fmt.Println()
			break
		}
		url = ""
	}
}

// ─────────────────────────── Version / Help ───────────────────────────────

const version = "1.0.0"

func printHelp() {
	printBanner()
	fmt.Println(colorize(Bold+White, "  Usage:"))
	fmt.Println(colorize(White, "    faster              Interactive mode"))
	fmt.Println(colorize(White, "    faster <URL>        Download directly (skips URL prompt)"))
	fmt.Println(colorize(White, "    faster --version    Show version"))
	fmt.Println(colorize(White, "    faster --help       Show this help"))
	fmt.Println()
	fmt.Println(colorize(Bold+White, "  Supported Sites:"))
	fmt.Println(colorize(Dim, "    YouTube, Instagram, Facebook, TikTok, Twitter/X,"))
	fmt.Println(colorize(Dim, "    Vimeo, Dailymotion, and 1000+ more via yt-dlp"))
	fmt.Println()
}

// ─────────────────────────── Entry Point ─────────────────────────────────

func main() {
	// Handle Ctrl+C gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		fmt.Println()
		printInfo("Interrupted. Goodbye! 👋")
		fmt.Println()
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	args := os.Args[1:]

	// Flags
	if len(args) > 0 {
		switch args[0] {
		case "--version", "-v":
			fmt.Printf("  faster v%s\n", version)
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		case "--check":
			printBanner()
			printInfo("Checking dependencies…")
			fmt.Println()
			ok := checkDependencies()
			fmt.Println()
			if ok {
				printSuccess("All dependencies satisfied — ready to download!")
			} else {
				printError("Some dependencies are missing. See above.")
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	printBanner()

	// Dependency check (silent — just warn if missing)
	if !checkDependency("yt-dlp") {
		printError("yt-dlp is not installed.")
		printInfo("Install it: pip install yt-dlp")
		fmt.Println()
		os.Exit(1)
	}

	// If a URL was passed as argument (e.g. from termux-url-opener)
	initialURL := ""
	if len(args) > 0 && isValidURL(args[0]) {
		initialURL = args[0]
		fmt.Println()
		printInfo("Received URL: " + colorize(Cyan, truncateString(initialURL, 70)))
	}

	// Count downloads in session
	sessionStart := time.Now()
	_ = sessionStart

	// Validate --format flag if provided
	for i, a := range args {
		if a == "--format" && i+1 < len(args) {
			initialURL = args[len(args)-1]
		}
		_ = strconv.Itoa(i) // keep import
	}

	run(reader, initialURL)
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
