package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/snakex21/devspace-go/internal/config"
)

// Op is what a tool is about to do with a path.
type Op int

const (
	// OpRead covers read, ls, grep and glob.
	OpRead Op = iota
	// OpWrite covers write, edit and mkdir.
	OpWrite
	// OpMove is a move/rename source or target: a write that may also remove
	// the path, so it must never be an extra root or contain a root.
	OpMove
	// OpExec is the shell working directory. The shell itself is not
	// confined, but starting a build in a read-only root would write there.
	OpExec
)

func (o Op) String() string {
	switch o {
	case OpRead:
		return "read"
	case OpWrite:
		return "write"
	case OpMove:
		return "move"
	case OpExec:
		return "exec"
	}
	return "unknown"
}

func (o Op) writes() bool { return o != OpRead }

// Origin tells where an MCP request came from.
type Origin string

const (
	// OriginLocal is a request made on this machine to the loopback listener.
	OriginLocal Origin = "local"
	// OriginRemote is a request forwarded by a tunnel or from another host.
	OriginRemote Origin = "remote"
)

// ErrAccessDenied wraps every refusal of an outside-workspace path.
var ErrAccessDenied = errors.New("access denied")

type extraRoot struct {
	configured string
	canonical  string
	writable   bool
}

type denyRule struct {
	path   string
	reason string
	// system rules protect operating-system and program folders. They are
	// lifted when the user configured an extra root inside that folder,
	// which is an explicit decision to work there. Credential rules never are.
	system bool
}

// AccessPolicy decides which paths outside the open workspace file tools may
// use. The zero value and a nil policy allow nothing outside the workspace.
//
// Order of checks for an outside path: tunnel requests are refused unless
// enabled; the path must be a plain absolute drive path (no UNC, device,
// stream or reserved names); it is canonicalised (links, junctions, 8.3
// names, case) and must lie in an extra root; built-in and configured deny
// rules apply on both the canonical and the lexical path; writes need a
// read-write root; moves may not touch a root itself.
type AccessPolicy struct {
	roots       []extraRoot
	deny        []denyRule
	allowRemote bool
	audit       *auditLog
}

// NewAccessPolicy builds the policy for cfg. Problems describe configured
// roots or deny entries that were ignored; the policy is usable regardless.
func NewAccessPolicy(cfg *config.Config) (*AccessPolicy, []string) {
	return newAccessPolicy(cfg, platformLocations())
}

func newAccessPolicy(cfg *config.Config, loc locations) (*AccessPolicy, []string) {
	policy := &AccessPolicy{allowRemote: cfg.OutsideAccess.Remote}
	var problems []string

	policy.deny = builtinDenyRules(loc)
	for _, dir := range []string{cfg.ConfigDir, cfg.StateDir} {
		if dir != "" {
			policy.deny = append(policy.deny, denyRule{path: dir, reason: "Dev Space Go configuration and state"})
		}
	}
	if cfg.ConfigDir != "" {
		parent := filepath.Dir(filepath.Clean(cfg.ConfigDir))
		for _, legacy := range []string{".webcoder", ".webcoder-state"} {
			policy.deny = append(policy.deny, denyRule{path: filepath.Join(parent, legacy), reason: "Dev Space Go configuration and state"})
		}
	}
	for _, entry := range cfg.OutsideAccess.Deny {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if !filepath.IsAbs(entry) {
			problems = append(problems, fmt.Sprintf("outsideAccess.deny entry %q is not an absolute path and is ignored", entry))
			continue
		}
		policy.deny = append(policy.deny, denyRule{path: entry, reason: "configured deny entry"})
	}
	// Compare deny rules with canonical paths: add the canonical form of
	// every rule next to its lexical one (they differ for subst drives,
	// junctioned profile folders, short names, ...).
	for _, rule := range append([]denyRule(nil), policy.deny...) {
		rule.path = filepath.Clean(rule.path)
		if canonical, err := Canonicalize(rule.path); err == nil && !samePathText(canonical, rule.path) {
			rule.path = canonical
			policy.deny = append(policy.deny, rule)
		}
	}

	for _, root := range cfg.OutsideAccess.Roots {
		configured := strings.TrimSpace(root.Path)
		access := strings.ToLower(strings.TrimSpace(root.Access))
		if access != "" && access != config.AccessRead && access != config.AccessReadWrite {
			continue // reported by config.Problems
		}
		if err := checkPlainAbsolute(configured); err != nil {
			problems = append(problems, fmt.Sprintf("extra root %q ignored: %v", configured, err))
			continue
		}
		configured = filepath.Clean(configured)
		if filepath.Dir(configured) == configured {
			problems = append(problems, fmt.Sprintf("extra root %s ignored: a whole drive or filesystem root is too broad; add the folders you need", configured))
			continue
		}
		info, err := os.Stat(configured)
		if err != nil || !info.IsDir() {
			problems = append(problems, fmt.Sprintf("extra root %s ignored: not an existing directory", configured))
			continue
		}
		canonical, err := Canonicalize(configured)
		if err != nil {
			problems = append(problems, fmt.Sprintf("extra root %s ignored: %v", configured, err))
			continue
		}
		if rule, ok := policy.matchDeny(canonical, canonical); ok && !rule.system {
			problems = append(problems, fmt.Sprintf("extra root %s ignored: it is inside a protected location (%s)", configured, rule.reason))
			continue
		}
		if reason := deniedName(canonical); reason != "" {
			problems = append(problems, fmt.Sprintf("extra root %s ignored: it is inside a protected location (%s)", configured, reason))
			continue
		}
		policy.roots = append(policy.roots, extraRoot{
			configured: configured,
			canonical:  canonical,
			writable:   access == config.AccessReadWrite,
		})
	}
	// Most specific root first, so a read-only subfolder of a read-write
	// root stays read-only.
	sort.SliceStable(policy.roots, func(i, j int) bool {
		return len(policy.roots[i].canonical) > len(policy.roots[j].canonical)
	})

	if len(policy.roots) > 0 && cfg.OutsideAccess.AuditEnabled() && cfg.StateDir != "" {
		policy.audit = newAuditLog(filepath.Join(cfg.StateDir, auditFileName))
	}
	return policy, problems
}

// RootInfo describes one usable extra root.
type RootInfo struct {
	Path     string
	Writable bool
}

// Roots lists the extra roots available to requests from origin.
func (p *AccessPolicy) Roots(origin Origin) []RootInfo {
	if p == nil || (origin != OriginLocal && !p.allowRemote) {
		return nil
	}
	infos := make([]RootInfo, 0, len(p.roots))
	for _, root := range p.roots {
		infos = append(infos, RootInfo{Path: root.configured, Writable: root.writable})
	}
	return infos
}

// Guard binds the policy to one request origin and tool for auditing.
func (p *AccessPolicy) Guard(origin Origin, tool string) *Guard {
	return &Guard{policy: p, origin: origin, tool: tool}
}

// Guard resolves tool paths: inside the workspace exactly as ResolveInRoot,
// outside it through the extra-roots policy. A nil Guard confines to the
// workspace.
type Guard struct {
	policy *AccessPolicy
	origin Origin
	tool   string
}

// Resolve maps a tool path argument to an absolute path. outside reports
// that the path was granted by an extra root (and is then canonical).
func (g *Guard) Resolve(wsRoot, input string, op Op) (path string, outside bool, err error) {
	path, err = ResolveInRoot(wsRoot, input)
	if err == nil || !errors.Is(err, ErrOutsideRoot) || g == nil || g.policy == nil || len(g.policy.roots) == 0 {
		return path, false, err
	}
	path, denial := g.policy.resolveOutside(wsRoot, input, op, g.origin)
	g.policy.audit.record(auditEntry{
		Origin: string(g.origin),
		Tool:   g.tool,
		Op:     op.String(),
		Input:  input,
		Path:   path,
		Allow:  denial == "",
		Reason: denial,
	})
	if denial != "" {
		return "", false, fmt.Errorf("path %s is outside the workspace and %s: %w", input, denial, ErrAccessDenied)
	}
	return path, true, nil
}

// Hidden reports whether a path met while walking (grep, glob, ls) inside an
// extra root must be skipped because a deny rule covers it. The walk itself
// starts from a canonical path and does not follow links.
func (g *Guard) Hidden(path string) bool {
	if g == nil || g.policy == nil {
		return true
	}
	root, ok := g.policy.rootFor(path)
	if !ok {
		return true
	}
	if rule, ok := g.policy.matchDeny(path, path); ok && !(rule.system && IsPathInsideRoot(root.canonical, rule.path)) {
		return true
	}
	return deniedName(path) != ""
}

func (p *AccessPolicy) resolveOutside(wsRoot, input string, op Op, origin Origin) (string, string) {
	if origin != OriginLocal && !p.allowRemote {
		return "", "extra roots are disabled for requests through the tunnel (outsideAccess.remote is false)"
	}
	cleanRoot, err := filepath.Abs(filepath.Clean(wsRoot))
	if err != nil {
		return "", "the workspace root cannot be resolved"
	}
	candidate, err := lexicalCandidate(cleanRoot, input)
	if err != nil {
		return "", err.Error()
	}
	if err := checkPlainAbsolute(candidate); err != nil {
		return "", err.Error()
	}
	canonical, err := Canonicalize(candidate)
	if err != nil {
		return "", "it cannot be resolved: " + err.Error()
	}
	root, ok := p.rootFor(canonical)
	if !ok {
		return canonical, "not inside a configured extra root"
	}
	if rule, ok := p.matchDeny(canonical, candidate); ok && !(rule.system && IsPathInsideRoot(root.canonical, rule.path)) {
		return canonical, "it is a protected location (" + rule.reason + ")"
	}
	if reason := deniedName(canonical); reason != "" {
		return canonical, "it is a protected location (" + reason + ")"
	}
	if reason := deniedName(candidate); reason != "" {
		return canonical, "it is a protected location (" + reason + ")"
	}
	if op.writes() && !root.writable {
		return canonical, "extra root " + root.configured + " is read-only"
	}
	if op == OpMove {
		for _, other := range p.roots {
			if IsPathInsideRoot(other.canonical, canonical) {
				return canonical, "an extra root cannot be moved or replaced"
			}
		}
		if realWs, err := Canonicalize(cleanRoot); err == nil && IsPathInsideRoot(realWs, canonical) {
			return canonical, "it contains the workspace root"
		}
	}
	return canonical, ""
}

func (p *AccessPolicy) rootFor(path string) (extraRoot, bool) {
	for _, root := range p.roots {
		if IsPathInsideRoot(path, root.canonical) {
			return root, true
		}
	}
	return extraRoot{}, false
}

func (p *AccessPolicy) matchDeny(canonical, lexical string) (denyRule, bool) {
	for _, rule := range p.deny {
		if IsPathInsideRoot(canonical, rule.path) || IsPathInsideRoot(lexical, rule.path) {
			return rule, true
		}
	}
	return denyRule{}, false
}

// checkPlainAbsolute accepts only ordinary absolute paths. On Windows that is
// a drive-letter path: UNC paths would make Windows authenticate to a remote
// server (leaking the user's NTLM hash) while being resolved, \\?\ and \\.\
// paths bypass normalisation or open devices, a second colon names an NTFS
// alternate data stream, trailing dots and spaces are silently dropped by
// Windows (so "id_rsa." opens "id_rsa"), and reserved names such as CON or
// NUL open devices in every folder.
func checkPlainAbsolute(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path %s is not absolute", path)
	}
	if runtime.GOOS != "windows" {
		return nil
	}
	return checkPlainWindowsPath(path)
}

func checkPlainWindowsPath(path string) error {
	volume := filepath.VolumeName(path)
	if len(volume) != 2 || volume[1] != ':' {
		return fmt.Errorf("UNC, device and \\\\?\\ paths are not allowed outside the workspace")
	}
	rest := path[2:]
	if strings.ContainsAny(rest, `:*?"<>|`) {
		return fmt.Errorf("path %s contains characters that name streams or devices", path)
	}
	for _, part := range strings.FieldsFunc(rest, func(r rune) bool { return r == '\\' || r == '/' }) {
		if part != strings.TrimRight(part, ". ") && part != "." && part != ".." {
			return fmt.Errorf("path component %q ends with a dot or space, which Windows ignores", part)
		}
		if isReservedWindowsName(part) {
			return fmt.Errorf("path component %q is a reserved device name", part)
		}
	}
	return nil
}

func isReservedWindowsName(name string) bool {
	base := strings.ToUpper(name)
	if i := strings.IndexByte(base, '.'); i >= 0 {
		base = base[:i]
	}
	base = strings.TrimRight(base, " ")
	switch base {
	case "CON", "PRN", "AUX", "NUL", "CONIN$", "CONOUT$":
		return true
	}
	if len(base) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) {
		return base[3] >= '0' && base[3] <= '9'
	}
	return false
}

// Names that hold credentials wherever they appear inside an extra root.
var (
	deniedDirNames = map[string]string{
		".ssh": "SSH keys", ".gnupg": "GnuPG keys", ".aws": "AWS credentials",
		".azure": "Azure credentials", ".kube": "Kubernetes credentials",
		".docker": "Docker credentials", ".password-store": "password store",
		".devspace": "Dev Space Go configuration", ".devspace-state": "Dev Space Go state",
		".webcoder": "Dev Space Go configuration", ".webcoder-state": "Dev Space Go state",
	}
	deniedFileNames = map[string]string{
		".git-credentials": "git credentials", ".netrc": "netrc credentials",
		"_netrc": "netrc credentials", ".npmrc": "npm token", ".pypirc": "PyPI token",
		".pgpass": "PostgreSQL passwords",
	}
	deniedKeyPrefixes = []string{"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519"}
)

// deniedName checks every path component against credential names. A .git
// directory's config file is denied because remote URLs often embed tokens.
func deniedName(path string) string {
	parts := strings.FieldsFunc(path, func(r rune) bool { return r == '/' || r == filepath.Separator })
	for i, part := range parts {
		name := part
		if pathsAreCaseInsensitive {
			name = strings.ToLower(name)
		}
		if reason, ok := deniedDirNames[name]; ok {
			return reason
		}
		if reason, ok := deniedFileNames[name]; ok {
			return reason
		}
		for _, prefix := range deniedKeyPrefixes {
			if strings.HasPrefix(name, prefix) {
				return "SSH private key"
			}
		}
		if name == ".git" && i+1 < len(parts) && strings.EqualFold(parts[i+1], "config") {
			return "git config (may contain credentials)"
		}
	}
	return ""
}

// locations are the per-user and system folders the built-in deny list is
// derived from. Tests inject their own.
type locations struct {
	goos         string
	home         string
	appData      string
	localAppData string
	systemRoot   string
	programDirs  []string
}

func platformLocations() locations {
	home, _ := os.UserHomeDir()
	loc := locations{goos: runtime.GOOS, home: home}
	if runtime.GOOS == "windows" {
		loc.appData = os.Getenv("APPDATA")
		loc.localAppData = os.Getenv("LOCALAPPDATA")
		loc.systemRoot = os.Getenv("SystemRoot")
		if loc.systemRoot == "" {
			loc.systemRoot = `C:\Windows`
		}
		for _, name := range []string{"ProgramFiles", "ProgramFiles(x86)", "ProgramW6432", "ProgramData"} {
			if dir := os.Getenv(name); dir != "" {
				loc.programDirs = append(loc.programDirs, dir)
			}
		}
	}
	return loc
}

func builtinDenyRules(loc locations) []denyRule {
	var rules []denyRule
	add := func(base, reason string, system bool, parts ...string) {
		if base == "" {
			return
		}
		rules = append(rules, denyRule{path: filepath.Join(append([]string{base}, parts...)...), reason: reason, system: system})
	}
	credentials := func(base string, parts ...string) { add(base, "credentials or browser profile", false, parts...) }

	credentials(loc.home, ".config", "gh")
	credentials(loc.home, ".claude", ".credentials.json")
	credentials(loc.home, ".codex", "auth.json")
	switch loc.goos {
	case "windows":
		for _, sub := range [][]string{
			{"Microsoft", "Credentials"}, {"Microsoft", "Protect"}, {"Microsoft", "Crypto"},
			{"Microsoft", "Vault"}, {"Microsoft", "SystemCertificates"},
			{"Mozilla"}, {"Opera Software"}, {"GitHub CLI"}, {"gh"},
		} {
			credentials(loc.appData, sub...)
		}
		for _, sub := range [][]string{
			{"Microsoft", "Credentials"}, {"Microsoft", "Vault"}, {"Microsoft", "Edge", "User Data"},
			{"Google", "Chrome", "User Data"}, {"BraveSoftware"}, {"Chromium"}, {"Vivaldi"},
			{"Mozilla"},
		} {
			credentials(loc.localAppData, sub...)
		}
		add(loc.systemRoot, "Windows system folder", true)
		for _, dir := range loc.programDirs {
			add(dir, "installed programs", true)
		}
	case "darwin":
		credentials(loc.home, "Library", "Keychains")
		credentials(loc.home, "Library", "Cookies")
		credentials(loc.home, "Library", "Application Support", "Google", "Chrome")
		credentials(loc.home, "Library", "Application Support", "Firefox")
		credentials(loc.home, "Library", "Application Support", "BraveSoftware")
		credentials(loc.home, "Library", "Application Support", "Microsoft Edge")
		credentials(loc.home, "Library", "Safari")
		for _, dir := range []string{"/System", "/Library", "/private", "/etc", "/usr", "/bin", "/sbin"} {
			add(dir, "system folder", true)
		}
	default:
		credentials(loc.home, ".mozilla")
		credentials(loc.home, ".config", "google-chrome")
		credentials(loc.home, ".config", "chromium")
		credentials(loc.home, ".config", "BraveSoftware")
		credentials(loc.home, ".local", "share", "keyrings")
		for _, dir := range []string{"/etc", "/usr", "/bin", "/sbin", "/lib", "/lib64", "/boot", "/dev", "/proc", "/sys", "/var", "/run", "/root"} {
			add(dir, "system folder", true)
		}
	}
	return rules
}

func samePathText(a, b string) bool {
	if pathsAreCaseInsensitive {
		return strings.EqualFold(a, b)
	}
	return a == b
}
