---
title: "Troubleshooting"
description: "Solutions for common Modularitea issues — including launch failures, pkexec errors, DNS problems, mirror failures, and more."
category: "Modularitea"
order: 11
---

# Troubleshooting

This guide provides solutions for common issues you may encounter when using Modularitea. Each issue follows a structured format: **Symptoms → Cause → Solution**.

---

## Modularitea Fails to Launch

### Symptoms

- Clicking the Modularitea icon does nothing.
- Launching from the terminal produces an error or the application closes immediately.
- A blank window appears and then closes.

### Possible Causes and Solutions

#### Cause 1: Missing Dependencies

One or more required libraries are not installed.

**Diagnosis:**

```bash
ldd $(which tealinux-modularity) | grep "not found"
```

**Solution:**

```bash
sudo pacman -S webkit2gtk gtk4 tealinux-modularitea-libs
```

#### Cause 2: Corrupted Package

The Modularitea package itself may be corrupted.

**Diagnosis:**

```bash
pacman -Qk tealinux-modularity
```

If the output shows missing files, reinstall:

**Solution:**

```bash
sudo pacman -S tealinux-modularity --overwrite '*'
```

#### Cause 3: Display Protocol Issues

On some configurations, Tauri applications may have issues with Wayland.

**Diagnosis:**

```bash
echo $XDG_SESSION_TYPE
```

**Solution — Try running with X11 backend:**

```bash
GDK_BACKEND=x11 tealinux-modularity
```

If this works, you can make it permanent by editing the application's `.desktop` file.

#### Cause 4: Corrupted Configuration File

The Modularitea configuration file may be corrupted.

**Solution:**

```bash
rm -f ~/.config/tealinux-modularity
tealinux-modularity
```

This resets Modularitea to its first-launch state, showing the splash screen again.

---

## Error: pkexec / PolicyKit Failures

### Symptoms

Operations that require administrator privileges fail with errors like:

```
command exited with code 126
```

or

```
Error: PERMISSION_DENIED
```

or

```
Not authorized to perform this action
```

### Possible Causes and Solutions

#### Cause 1: PolicyKit Agent Not Running

The PolicyKit authentication agent may not be running. This is common in minimal desktop setups.

**Diagnosis:**

```bash
ps aux | grep polkit
```

**Solution:**

```bash
# Install PolicyKit agent if missing
sudo pacman -S polkit polkit-gnome

# Start the agent (for COSMIC/GNOME)
/usr/lib/polkit-gnome/polkit-gnome-authentication-agent-1 &
```

#### Cause 2: PolicyKit Rules Missing

Modularitea's PolicyKit rules may not be installed.

**Diagnosis:**

```bash
ls /usr/share/polkit-1/actions/ | grep tealinux
```

**Solution:**

```bash
sudo pacman -S tealinux-modularitea-libs
```

This reinstalls the package, which includes the PolicyKit action files.

#### Cause 3: User Cancelled Authentication

If you clicked "Cancel" on the password dialog, the operation was intentionally aborted. This is not an error — simply retry and enter your password.

#### Cause 4: User Not in `wheel` Group

The user account must be in the `wheel` group to use `pkexec`.

**Diagnosis:**

```bash
groups $USER
```

**Solution:**

```bash
sudo usermod -aG wheel $USER
```

Log out and log back in for the change to take effect.

---

## DNS Change Fails

### Symptoms

- DNS provider switch fails with an error notification.
- `/etc/resolv.conf` remains unchanged after switching.

### Possible Causes and Solutions

#### Cause 1: Binary Not Found

**Diagnosis:**

```bash
which modularitea-dns-changer
```

**Solution:**

```bash
sudo pacman -S tealinux-modularitea-libs
```

#### Cause 2: resolv.conf Is Immutable

Some system configurations set the immutable attribute on `/etc/resolv.conf`.

**Diagnosis:**

```bash
lsattr /etc/resolv.conf
```

If you see `i` in the output, the file is immutable.

**Solution:**

```bash
sudo chattr -i /etc/resolv.conf
```

Then retry the DNS change in Modularitea.

#### Cause 3: NetworkManager Overwriting resolv.conf

NetworkManager may overwrite `/etc/resolv.conf` on network changes.

**Solution:**

Create `/etc/NetworkManager/conf.d/dns.conf`:

```ini
[main]
dns=none
```

Then restart NetworkManager:

```bash
sudo systemctl restart NetworkManager
```

---

## Mirror Refresh Fails

### Symptoms

- Mirror refresh times out or fails with an error.
- Error message mentions "network" or "timeout".

### Possible Causes and Solutions

#### Cause 1: No Internet Connection

**Diagnosis:**

```bash
ping -c 3 archlinux.org
```

**Solution:** Ensure your internet connection is active and stable.

#### Cause 2: reflector Not Installed

**Diagnosis:**

```bash
which reflector
```

**Solution:**

```bash
sudo pacman -S reflector
```

#### Cause 3: Firewall Blocking HTTPS

Some firewalls block outgoing HTTPS connections needed by `reflector`.

**Diagnosis:**

```bash
curl -sI https://archlinux.org/mirrors/status/
```

**Solution:** Ensure your firewall allows outgoing connections on port 443.

#### Cause 4: Invalid Country Name

If you specified a country that `reflector` does not recognize.

**Solution:** Leave the country field empty to test worldwide mirrors, or use the exact country name as listed on the [Arch Mirror Status page](https://archlinux.org/mirrors/status/).

---

## Swap Toggle Fails

### Symptoms

- Swap enable/disable fails with an error notification.
- Swap status does not change after toggling.

### Possible Causes and Solutions

#### Cause 1: Binary Not Found

**Diagnosis:**

```bash
which modularitea-swap
```

**Solution:**

```bash
sudo pacman -S tealinux-modularitea-libs
```

#### Cause 2: systemd-zram-generator Not Available

**Diagnosis:**

```bash
pacman -Q zram-generator
```

**Solution:**

```bash
sudo pacman -S zram-generator
```

#### Cause 3: Conflicting Swap Configuration

Another swap configuration (e.g., a swap partition or swap file) may be interfering.

**Diagnosis:**

```bash
swapon --show
cat /etc/fstab | grep swap
```

**Solution:** Disable conflicting swap entries in `/etc/fstab` and try again.

---

## Profile Installation Fails

### Symptoms

- A profile installation fails partway through.
- Some packages install but others fail.
- Error messages mention "package not found" or "dependency conflicts".

### Possible Causes and Solutions

#### Cause 1: Outdated Package Database

**Diagnosis:** Check when the database was last synced:

```bash
stat /var/lib/pacman/sync/core.db
```

**Solution:**

```bash
sudo pacman -Sy
```

Then retry the installation.

#### Cause 2: Package Name Changed or Removed

A package in the profile TOML may have been renamed or removed from the repository.

**Solution:** Check if the package exists:

```bash
pacman -Ss <package-name>
```

If the package no longer exists, it may need to be removed from the profile. Report this as a bug.

#### Cause 3: Dependency Conflicts

Two packages in the profile may have conflicting dependencies.

**Solution:**

```bash
# Check for conflicts
sudo pacman -S <package1> <package2>
```

Review the conflict messages and resolve manually, or report the issue.

#### Cause 4: Insufficient Disk Space

**Diagnosis:**

```bash
df -h /
```

**Solution:** Free up disk space by cleaning the package cache or removing unused packages (see [Maintenance](./maintenance)).

---

## GRUB Theme Fails to Apply

### Symptoms

- GRUB theme change appears to succeed but the old theme shows on next boot.
- Error message during theme application.

### Possible Causes and Solutions

#### Cause 1: grub-mkconfig Fails

**Diagnosis:**

```bash
sudo grub-mkconfig -o /boot/grub/grub.cfg 2>&1 | tail -20
```

**Solution:** Review the error output and fix any configuration issues in `/etc/default/grub`.

#### Cause 2: Theme Files Missing

**Solution:** Reinstall the GRUB themes:

```bash
sudo pacman -S tealinux-modularitea-libs
```

---

## General Debugging Steps

If none of the above solutions resolve your issue, follow these general debugging steps:

1. **Run Modularitea from the terminal** to see error output:

```bash
RUST_LOG=debug tealinux-modularity
```

2. **Check system logs:**

```bash
journalctl -b --no-pager | grep -i modularitea | tail -50
```

3. **Verify all Modularitea packages are installed:**

```bash
pacman -Q | grep -i tealinux
pacman -Q | grep -i modularitea
```

4. **Check for broken packages:**

```bash
pacman -Qk tealinux-modularity tealinux-modularitea-libs 2>&1 | grep -v "0 missing"
```

5. **Report the bug** if the issue persists (see [Developer Guide → Bug Reporting](./developer#bug-reporting)).
