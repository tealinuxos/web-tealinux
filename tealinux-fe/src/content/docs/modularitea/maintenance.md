---
title: "Maintenance"
description: "System maintenance features in Modularitea — including system updates, cache cleaning, removing unused packages, and health checks."
category: "Modularitea"
order: 10
---

# Maintenance

Modularitea includes several maintenance features to help keep your TeaLinuxOS system clean, updated, and running smoothly. This guide covers each maintenance operation and provides recommendations for when and how often to use them.

---

## System Update

### What It Does

A system update synchronizes your package database and upgrades all installed packages to their latest available versions. This is equivalent to running `pacman -Syu` in the terminal.

### Why It Matters

Arch Linux is a **rolling release** distribution, which means:

- There are no major version upgrades — updates are delivered continuously.
- Keeping your system up to date ensures you have the latest security patches, bug fixes, and features.
- Delaying updates for too long can cause package conflicts and dependency issues.

### How to Update

1. Open **Modularitea**.
2. Navigate to the **Settings** page.
3. The system will use the `update_db` command to sync the package database.
4. Authenticate with your administrator password.
5. Wait for the update to complete.
6. Reboot if kernel or system-critical packages were updated.

### Recommended Frequency

| Scenario | Recommended Update Interval |
| :--- | :--- |
| Daily driver workstation | Weekly |
| Development machine | Every 1–2 weeks |
| Server or production system | Monthly, after testing |

> **Tip:** Always check the [Arch Linux news page](https://archlinux.org/news/) before a major update for any manual intervention notices.

---

## Clean Package Cache

### What It Does

Cleans the pacman package cache located at `/var/cache/pacman/pkg/`. This directory stores every package file you have ever downloaded through `pacman`, which can accumulate to several gigabytes over time.

### Why It Matters

| Cache Size | Impact |
| :--- | :--- |
| < 1 GB | Normal — no action needed |
| 1–5 GB | Moderate — consider cleaning |
| 5+ GB | Large — recommended to clean |
| 10+ GB | Very large — should clean soon |

### How to Clean

1. Open **Modularitea** and navigate to **Settings**.
2. Locate the **Package Cache** card.
3. The current cache size is displayed (e.g., "3.7 GB").
4. Click **Clean Cache**.
5. Authenticate with your administrator password.
6. A success notification confirms the cache was cleaned.

### What Gets Removed

The cache cleaner removes downloaded package files (`.pkg.tar.zst`). It does **not**:

- Remove installed packages
- Affect your system configuration
- Remove currently installed package versions

### Recommended Frequency

| Scenario | Recommended Cleaning Interval |
| :--- | :--- |
| Limited disk space (< 40 GB) | Monthly |
| Moderate disk space (40–100 GB) | Every 2–3 months |
| Ample disk space (> 100 GB) | Every 6 months |

> **Note:** If you frequently need to downgrade packages, keep the cache intact — cached files allow you to roll back to previous versions without re-downloading.

---

## Remove Unused Packages

### What It Does

Identifies and removes **orphan packages** — packages that were installed as dependencies but are no longer required by any installed package.

### Why It Matters

Over time, as you install and remove software, dependency packages can become "orphaned." These packages:

- Consume disk space unnecessarily
- May receive security updates that you never benefit from
- Clutter the package database

### How to Remove Orphans (Terminal)

Currently, orphan removal can be performed via the terminal:

```bash
# List orphaned packages
pacman -Qdt

# Remove orphaned packages
sudo pacman -Rns $(pacman -Qdtq)
```

> **Tip:** Always review the list of orphan packages before removing them. In rare cases, a package you use may be incorrectly identified as an orphan if it was installed as a dependency rather than explicitly.

### Recommended Frequency

Clean up orphan packages every 1–2 months, or after uninstalling large software profiles.

---

## Health Check

### What It Does

A health check assesses the overall state of your system, looking for common issues that could cause problems.

### Items to Check

| Check | What It Verifies | Terminal Command |
| :--- | :--- | :--- |
| **Failed systemd services** | Services that failed to start | `systemctl --failed` |
| **Pacman database consistency** | Package database integrity | `pacman -Dk` |
| **Disk space** | Available storage on all partitions | `df -h` |
| **ZRAM swap status** | Whether swap is active and functioning | `swapon --show` |
| **Journal log errors** | Recent critical system errors | `journalctl -p 3 -b` |
| **DNS resolution** | Whether DNS queries resolve correctly | `dig google.com` |

### Performing a Health Check (Terminal)

You can run a quick system health check from the terminal:

```bash
# Check for failed services
echo "=== Failed Services ==="
systemctl --failed

# Check pacman database
echo "=== Package Database ==="
pacman -Dk

# Check disk space
echo "=== Disk Space ==="
df -h /

# Check swap status
echo "=== Swap Status ==="
swapon --show

# Check recent errors
echo "=== Recent Errors ==="
journalctl -p 3 -b --no-pager | tail -20
```

### Recommended Frequency

Run a health check monthly, or immediately if you notice unusual behavior such as:

- Applications crashing unexpectedly
- Slow system performance
- Package installation failures
- Network connectivity issues

---

## Maintenance Schedule Recommendation

Here is a suggested maintenance schedule for TeaLinuxOS users:

| Task | Frequency | Priority |
| :--- | :--- | :--- |
| System update | Weekly | **High** |
| Clean package cache | Monthly | Medium |
| Remove orphan packages | Every 2 months | Low |
| Health check | Monthly | Medium |
| Mirror refresh | Monthly | Medium |

### Automated Maintenance

For users who prefer automated maintenance, you can create a simple systemd timer. However, for most desktop users, performing these tasks manually through Modularitea provides better control and awareness of what is happening on your system.

---

## Tips for System Longevity

1. **Update regularly** — Do not let updates accumulate for more than 2 weeks on a rolling release distribution.
2. **Read the Arch News** — Before updating, check for any manual intervention notices.
3. **Keep some cache** — If you have the disk space, keeping the last 2–3 versions of cached packages allows for easy rollbacks.
4. **Monitor disk space** — A full disk can cause various system failures including failed updates and application crashes.
5. **Review journal logs** — Periodically check `journalctl -p 3 -b` for critical errors that may indicate hardware or software issues.
