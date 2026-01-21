# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability, please email the maintainer directly rather than opening a public issue.

## Security Considerations

### API Key Storage

**Current Implementation:**
- API keys are stored in plaintext in `~/.unraidcli/config.yaml`
- File permissions are set to `0600` (owner read/write only)

**Risks:**
- Keys are not encrypted at rest
- Vulnerable if an attacker gains access to your user account
- Visible in system backups if not excluded

**Best Practices:**
1. **Use least-privilege API keys** - Create API keys with only the permissions you need
2. **Rotate keys regularly** - Periodically regenerate API keys
3. **Exclude from backups** - Add `~/.unraidcli/` to backup exclusions
4. **Use environment variables** - For sensitive environments, consider using `UNRAID_API_KEY` env var (feature request)

### Command Line Exposure

**Risk:** When using `--apikey` flag, the key is visible in:
- Shell history (`~/.bash_history`, `~/.zsh_history`)
- Process listings (`ps aux`)

**Mitigation:**
```bash
# Option 1: Use interactive prompt (future feature)
unraidcli config set --url http://192.168.1.100
# Would prompt for API key without echoing

# Option 2: Clear shell history after setup
unraidcli config set --url http://192.168.1.100 --apikey YOUR_KEY
history -d $(history 1)

# Option 3: Read from stdin (future feature)
echo "YOUR_API_KEY" | unraidcli config set --url http://192.168.1.100 --apikey-stdin
```

### Transport Security

**Warning:** Using HTTP URLs exposes your API key and data in transit.

**Recommendation:**
- Always use HTTPS when accessing Unraid remotely
- For local networks, HTTP is acceptable if you trust the network
- Consider using Tailscale or WireGuard for remote access

### Network Security

**Risks:**
- Man-in-the-middle attacks on HTTP connections
- API key interception on untrusted networks

**Best Practices:**
1. Use HTTPS for all remote connections
2. Use VPN (Tailscale, WireGuard) for remote management
3. Only use HTTP on trusted local networks
4. Avoid using unraidcli on public Wi-Fi networks

### Multi-User Systems

**Risk:** On shared systems, other users with elevated privileges could:
- Read your config file
- View running processes with API keys in command line

**Mitigation:**
- Only use on systems you control
- Be aware that root/admin users can access your config
- Use dedicated user accounts for Unraid management

### Secure Configuration Example

```bash
# 1. Initial setup with HTTPS
unraidcli config set --url https://unraid.example.com --apikey YOUR_KEY

# 2. Clear the command from history
history -d $(history 1)

# 3. Verify secure storage
ls -la ~/.unraidcli/config.yaml
# Should show: -rw------- (600 permissions)

# 4. Verify HTTPS is configured
unraidcli config show
# Should show URL starting with https://
```

### Dependency Security

We regularly monitor dependencies for known vulnerabilities. Run:
```bash
go list -json -m all | nancy sleuth
```

## Security Features

### Current
✅ Config file permissions (0600)
✅ API key masking in output
✅ Connection testing before saving credentials
✅ HTTPS support

### Planned
⏳ Interactive API key prompt (no command line exposure)
⏳ Environment variable support (`UNRAID_API_KEY`)
⏳ API key encryption at rest (optional)
⏳ HTTPS enforcement warnings

## Responsible Disclosure

We take security seriously. If you discover a vulnerability:

1. **Do not** open a public GitHub issue
2. Email the maintainer with details
3. Allow reasonable time for a fix before public disclosure
4. We will credit you in the security advisory (if desired)

## Security Updates

Security updates will be released as soon as possible and marked clearly in release notes.
