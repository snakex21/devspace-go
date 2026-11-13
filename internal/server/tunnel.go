package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog/log"
	"github.com/snakex21/devspace-go/internal/locales"
	"github.com/snakex21/devspace-go/internal/procutil"
)

const (
	cloudflaredURLTimeout  = 30 * time.Second
	cloudflaredMaxAttempts = 2
	tunnelProbeInterval    = 500 * time.Millisecond
	maxTunnelLogLine       = 64 * 1024
)

var (
	cloudflaredURLPattern = regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com\b`)
	pinggyURLPattern      = regexp.MustCompile(`https://[a-zA-Z0-9-]+\.(a\.)?pinggy\.(link|io|xyz)\b`)
)

// tunnelProcess owns a child until Wait finishes, including on startup failure.
// err is read only after done closes.
type tunnelProcess struct {
	url    string
	done   chan struct{}
	cancel context.CancelFunc
	err    error
}

func (p *tunnelProcess) stop() {
	p.cancel()
	<-p.done
}

type tunnelProvider func(context.Context, string) (*tunnelProcess, error)

// tunnelStarter launches one tunnel. It returns errNoTunnelProvider when no
// tunnel program is installed, so the supervisor does not retry forever.
type tunnelStarter func(context.Context, string) (*tunnelProcess, error)

var errNoTunnelProvider = errors.New("no tunnel program (cloudflared or ssh) found")

func (s *Server) startTunnel(ctx context.Context, origin string) (*tunnelProcess, error) {
	return s.startTunnelWithProviders(ctx, origin, s.startCloudflared, s.startPinggy)
}

func (s *Server) startTunnelWithProviders(ctx context.Context, origin string, cloudflared, pinggy tunnelProvider) (*tunnelProcess, error) {
	providers := []tunnelProvider{}
	for attempt := 0; attempt < cloudflaredMaxAttempts; attempt++ {
		providers = append(providers, cloudflared)
	}
	providers = append(providers, pinggy)
	var lastErr error
	for _, provider := range providers {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		tunnel, err := provider(ctx, origin)
		if err == nil && tunnel != nil {
			if ctx.Err() != nil {
				tunnel.stop()
				return nil, ctx.Err()
			}
			printTunnelURL(tunnel.url)
			return tunnel, nil
		}
		if err != nil {
			lastErr = err
			if ctx.Err() == nil {
				log.Warn().Err(err).Msg("tunnel startup failed")
			}
		}
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if lastErr == nil {
		return nil, errNoTunnelProvider
	}
	log.Warn().Msg("no tunnel is ready; server remains available locally")
	return nil, lastErr
}

// backoffPolicy spaces out tunnel restarts: exponential growth from base up
// to max, with "equal jitter" (a random delay in [d/2, d]) so that many
// instances do not retry in lockstep. A tunnel that stayed up for at least
// resetAfter counts as healthy, and the next failure starts from base again.
type backoffPolicy struct {
	base       time.Duration
	max        time.Duration
	resetAfter time.Duration
}

var defaultTunnelBackoff = backoffPolicy{base: 5 * time.Second, max: 5 * time.Minute, resetAfter: 2 * time.Minute}

func (p backoffPolicy) delay(attempt int, random func(int64) int64) time.Duration {
	d := p.base
	for i := 0; i < attempt && d < p.max; i++ {
		d *= 2
	}
	if d > p.max {
		d = p.max
	}
	half := d / 2
	if half <= 0 {
		return d
	}
	return half + time.Duration(random(int64(half)+1))
}

// superviseTunnel keeps a public tunnel running until ctx is canceled. When
// the tunnel process exits (network loss, sleep/hibernation, a crash of the
// tunnel program) it is restarted with backoff. A restarted Quick Tunnel gets
// a new URL, which is printed in a new TUNNEL ACTIVE banner.
func superviseTunnel(ctx context.Context, origin string, start tunnelStarter, policy backoffPolicy) {
	attempt := 0
	for {
		tunnel, err := start(ctx, origin)
		if ctx.Err() != nil {
			if tunnel != nil {
				tunnel.stop()
			}
			return
		}
		if errors.Is(err, errNoTunnelProvider) {
			log.Info().Msg("no tunnel program found; server is available locally only")
			return
		}
		if tunnel != nil {
			startedAt := time.Now()
			select {
			case <-ctx.Done():
				tunnel.stop()
				return
			case <-tunnel.done:
				tunnel.stop()
				log.Error().Err(tunnel.err).Msg("tunnel exited; public URL is no longer active; server remains available locally")
			}
			if time.Since(startedAt) >= policy.resetAfter {
				attempt = 0
			}
		}
		delay := policy.delay(attempt, rand.Int64N)
		attempt++
		log.Warn().Dur("retry_in", delay).Int("attempt", attempt).Msg("restarting tunnel")
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (s *Server) startCloudflared(ctx context.Context, origin string) (*tunnelProcess, error) {
	executable := findCloudflaredExecutable()
	if executable == "" {
		return nil, nil
	}
	fmt.Printf("\n🔗  %s\n    %s\n\n", locales.T("tunnel.starting_cloudflared"), executable)
	return startTunnelProcess(ctx, func(ctx context.Context) *exec.Cmd {
		return exec.CommandContext(ctx, executable, "tunnel", "--url", origin)
	}, cloudflaredURLPattern, cloudflaredURLTimeout, waitForTunnelHealth, func(line string) { fmt.Println(line) })
}

func (s *Server) startPinggy(ctx context.Context, origin string) (*tunnelProcess, error) {
	executable, err := exec.LookPath("ssh")
	if err != nil {
		return nil, nil
	}
	originURL, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}
	// Forward to the actual bound address, including an ephemeral port or IPv6.
	target := "0:" + net.JoinHostPort(originURL.Hostname(), originURL.Port())
	fmt.Printf("\n🔗  %s\n\n", locales.T("tunnel.starting_pinggy"))
	return startTunnelProcess(ctx, func(ctx context.Context) *exec.Cmd {
		return exec.CommandContext(ctx, executable, "-p", "443",
			"-o", "StrictHostKeyChecking=accept-new",
			"-o", "BatchMode=yes",
			"-o", "ExitOnForwardFailure=yes",
			"-o", "ConnectTimeout=10",
			"-o", "ServerAliveInterval=30", "-o", "ServerAliveCountMax=3",
			"-R", target, "a.pinggy.io")
	}, pinggyURLPattern, 15*time.Second, waitForTunnelHealth, func(line string) { fmt.Println(line) })
}

// startTunnelProcess continuously drains both output streams. A URL only means
// it was allocated; require a successful origin health check before advertising it.
func startTunnelProcess(parent context.Context, command func(context.Context) *exec.Cmd, pattern *regexp.Regexp, timeout time.Duration, probe func(context.Context, string) error, output func(string)) (*tunnelProcess, error) {
	ctx, cancel := context.WithCancel(parent)
	p := &tunnelProcess{cancel: cancel, done: make(chan struct{})}
	urls := make(chan string, 1)
	onLine := func(line string) {
		if match := pattern.FindString(line); match != "" {
			select {
			case urls <- match:
			default:
			}
		}
		if output != nil {
			output(line)
		}
	}
	cmd := command(ctx)
	// Let os/exec own and drain the pipes through Wait, instead of abandoning
	// StdoutPipe/StderrPipe readers as soon as the first URL is found.
	cmd.Stdout = &tunnelOutput{onLine: onLine}
	cmd.Stderr = &tunnelOutput{onLine: onLine}
	cmd.WaitDelay = 2 * time.Second
	// The tunnel must never outlive this process: an orphaned tunnel would keep
	// publishing the port (and a later instance on it) under the old URL.
	// The job object / parent-death signal also covers crashes and hard kills.
	tree := procutil.Prepare(cmd, true)
	if err := tree.Start(); err != nil {
		cancel()
		return nil, err
	}
	go func() { p.err = cmd.Wait(); tree.Release(); close(p.done) }()
	keep := false
	defer func() {
		if !keep {
			p.stop()
		}
	}()
	startupCtx, startupCancel := context.WithTimeout(ctx, timeout)
	defer startupCancel()
	select {
	case p.url = <-urls:
	case <-p.done:
		return nil, tunnelExitError(p.err)
	case <-startupCtx.Done():
		return nil, fmt.Errorf("waiting for tunnel URL: %w", startupCtx.Err())
	}
	// Give process exit/cancellation precedence over a buffered stale URL.
	select {
	case <-p.done:
		return nil, tunnelExitError(p.err)
	default:
	}
	if err := startupCtx.Err(); err != nil {
		return nil, err
	}
	ready := make(chan error, 1)
	go func() { ready <- probe(startupCtx, p.url) }()
	select {
	case err := <-ready:
		if err != nil {
			return nil, fmt.Errorf("tunnel URL allocated but not ready: %w", err)
		}
	case <-p.done:
		return nil, tunnelExitError(p.err)
	case <-startupCtx.Done():
		return nil, fmt.Errorf("waiting for tunnel readiness: %w", startupCtx.Err())
	}
	select {
	case <-p.done:
		return nil, tunnelExitError(p.err)
	default:
	}
	if err := startupCtx.Err(); err != nil {
		return nil, err
	}
	keep = true
	return p, nil
}

func tunnelExitError(err error) error {
	if err == nil {
		return errors.New("tunnel process exited before becoming ready")
	}
	return fmt.Errorf("tunnel process exited: %w", err)
}

// A bounded line collector keeps draining even if a child writes an oversized
// line. Each output stream has its own writer, called serially by os/exec.
type tunnelOutput struct {
	line      []byte
	oversized bool
	onLine    func(string)
}

func (w *tunnelOutput) Write(p []byte) (int, error) {
	n := len(p)
	for len(p) > 0 {
		end := bytes.IndexByte(p, '\n')
		chunk := p
		if end >= 0 {
			chunk = p[:end]
		}
		if !w.oversized {
			if len(w.line)+len(chunk) > maxTunnelLogLine {
				w.oversized = true
				w.line = w.line[:0]
			} else {
				w.line = append(w.line, chunk...)
			}
		}
		if end < 0 {
			break
		}
		if !w.oversized {
			w.onLine(strings.TrimSuffix(string(w.line), "\r"))
		}
		w.line = w.line[:0]
		w.oversized = false
		p = p[end+1:]
	}
	return n, nil
}

func waitForTunnelHealth(ctx context.Context, publicURL string) error {
	client := &http.Client{
		Timeout: 3 * time.Second,
		// A login/error-page redirect is not proof that this origin is reachable.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	defer client.CloseIdleConnections()
	return waitForTunnelHealthWithClient(ctx, publicURL, client, tunnelProbeInterval)
}

func waitForTunnelHealthWithClient(ctx context.Context, publicURL string, client *http.Client, interval time.Duration) error {
	var lastErr error
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, publicURL+"/healthz", nil)
		if err != nil {
			return err
		}
		req.Header.Set("Cache-Control", "no-cache")
		response, err := client.Do(req)
		if err == nil {
			var health struct {
				OK   bool   `json:"ok"`
				Name string `json:"name"`
			}
			err = json.NewDecoder(io.LimitReader(response.Body, 4096)).Decode(&health)
			response.Body.Close()
			if response.StatusCode == http.StatusOK && err == nil && health.OK && health.Name == "devspace-go" {
				return nil
			}
			err = fmt.Errorf("unexpected health response (HTTP %d)", response.StatusCode)
		}
		lastErr = err
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("%w (last check: %v)", ctx.Err(), lastErr)
		case <-timer.C:
		}
	}
}

func findCloudflaredExecutable() string {
	names := []string{"cloudflared"}
	if runtime.GOOS == "windows" {
		names = []string{"cloudflared.exe", "cloudflared"}
	}
	var dirs []string

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		for dir := exeDir; dir != ""; dir = filepath.Dir(dir) {
			dirs = append(dirs, filepath.Join(dir, "tools"), dir)
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(wd, "tools"), wd)
	}

	seen := map[string]bool{}
	for _, dir := range dirs {
		cleanDir := filepath.Clean(dir)
		if seen[cleanDir] {
			continue
		}
		seen[cleanDir] = true
		for _, name := range names {
			candidate := filepath.Join(cleanDir, name)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() && (runtime.GOOS == "windows" || info.Mode().Perm()&0111 != 0) {
				return candidate
			}
		}
	}

	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}

func printTunnelURL(url string) {
	mcpURL := url + "/mcp"
	supportsSSE := !strings.HasSuffix(strings.TrimPrefix(url, "https://"), ".trycloudflare.com")
	lines := []string{
		"🌐 " + locales.T("tunnel.active"),
		mcpURL,
		"",
		locales.T("tunnel.paste_chatgpt"),
		mcpURL,
	}
	if supportsSSE {
		lines = append(lines, url+"/sse", locales.T("tunnel.try_sse"))
	}
	width := 54
	for _, line := range lines {
		if lineWidth := utf8.RuneCountInString(line); lineWidth > width {
			width = lineWidth
		}
	}

	fmt.Println()
	fmt.Printf("╔%s╗\n", strings.Repeat("═", width+4))
	for _, line := range lines {
		padding := strings.Repeat(" ", width-utf8.RuneCountInString(line))
		fmt.Printf("║  %s%s  ║\n", line, padding)
	}
	fmt.Printf("╚%s╝\n", strings.Repeat("═", width+4))
	fmt.Println()
}
