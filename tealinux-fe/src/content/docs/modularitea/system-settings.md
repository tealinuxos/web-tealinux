---
title: "System Settings"
description: "Configure system preferences including theme, DNS, mirrors, swap, CPU performance, and package cache through the Modularitea settings panel."
category: "Modularitea"
order: 6
---

# System Settings

The System Settings page in Modularitea provides centralized access to key system configuration options. Access it by clicking the **Settings** button (⚙️) at the bottom of the sidebar.

The settings panel is organized into a responsive grid layout with the following modules:

---

## Package Cache

The Package Cache module shows the current size of your pacman package cache (`/var/cache/pacman/pkg/`) and allows you to clean it with a single click.

### What Is Package Cache?

Every time you install or update a package with `pacman`, the downloaded package file is stored in `/var/cache/pacman/pkg/`. Over time, this cache can grow to several gigabytes, consuming valuable disk space.

### How to Use

1. The module displays the current cache size (e.g., "2.3 GB").
2. Click the **Clean Cache** button to remove cached package files.
3. Authenticate with your administrator password.
4. A success notification confirms the operation.

### Impact on System

- **Safe to clean** — Cleaning the cache removes downloaded `.pkg.tar.zst` files. It does not remove installed packages.
- **Trade-off** — After cleaning, if you need to downgrade a package, you will need to re-download it. For most users, this is not a concern.

---

## Swap Memory

The Swap Memory module allows you to enable or disable swap memory on your system.

### What Is Swap?

Swap is a section of storage used as virtual memory when your physical RAM is full. TeaLinuxOS uses **ZRAM** — a compressed block device in RAM — which is faster than traditional disk-based swap.

### How to Use

1. The module displays the current swap status (**Enabled** or **Disabled**).
2. Toggle the switch to enable or disable swap.
3. Authenticate with your administrator password.
4. The change takes effect immediately.

### How It Works

- **Enable** — Creates a `zram-generator.conf` file in `/etc/systemd/` and activates the ZRAM swap device.
- **Disable** — Removes the configuration file and deactivates the swap device.

> **Warning:** If your system has limited RAM (4 GB or less), disabling swap may cause applications to crash or the system to become unresponsive when memory usage is high. Keep swap enabled unless you have ample RAM.

For detailed information, see the [Swap Manager](./swap-manager) page.

---

## System Theme

The System Theme module allows you to switch between **Light** and **Dark** mode for the Modularitea interface.

### How to Use

1. Click the theme toggle to switch between light and dark mode.
2. The change is applied immediately — no restart required.

> **Note:** This setting affects only the Modularitea application. Your desktop environment theme is managed separately through your DE's settings.

---

## Mirror Settings

The Mirror Settings module allows you to refresh your Arch Linux package mirrors to find the fastest available server.

### What Are Mirrors?

Mirrors are servers that host copies of the Arch Linux package repositories. The speed of your package downloads depends on which mirror you are connected to. Over time, mirrors change — some become slow, others come online.

### How to Use

1. Click the **Refresh Mirror** button.
2. Optionally select a specific country to prioritize (e.g., "Indonesia").
3. Authenticate with your administrator password.
4. Wait for the process to complete — this may take 30–60 seconds.

### How It Works

Behind the scenes, Modularitea uses `reflector` — the official Arch Linux mirror ranking tool — to:

1. Fetch the current list of available mirrors.
2. Test each mirror for speed and availability.
3. Rank mirrors by download speed.
4. Write the fastest mirrors to `/etc/pacman.d/mirrorlist`.

For detailed information, see the [Mirror Manager](./mirror-manager) page.

---

## DNS Configuration

The DNS Configuration module allows you to switch your system's DNS provider with a single click.

### Available Providers

| Provider | Primary DNS | Secondary DNS |
| :--- | :--- | :--- |
| **Cloudflare** | `1.1.1.1` | `1.0.0.1` |
| **Google** | `8.8.8.8` | `8.8.4.4` |
| **Quad9** | `9.9.9.9` | `149.112.112.112` |
| **OpenDNS** | `208.67.222.222` | `208.67.220.220` |
| **AdGuard** | `94.140.14.14` | `94.140.15.15` |

### How to Use

1. The module displays your current DNS provider (if detected).
2. Select a new provider from the dropdown.
3. Authenticate with your administrator password.
4. A success notification confirms the change.

### Impact on System

- Changes are applied by modifying `/etc/resolv.conf`.
- The change takes effect immediately for all new DNS queries.
- Existing connections may continue using the old DNS until they are re-established.

For detailed information, see the [DNS Manager](./dns-manager) page.

---

## CPU Performance

The CPU Performance module allows you to adjust the CPU frequency scaling governor, controlling the balance between performance and power consumption.

### Available Profiles

| Profile | Description | Best For |
| :--- | :--- | :--- |
| **Performance** | CPU runs at maximum frequency at all times | Compilation, rendering, gaming |
| **Powersave** | CPU runs at minimum frequency to save energy | Battery life, quiet operation |
| **Ondemand** | CPU dynamically adjusts frequency based on load | General daily use (recommended) |

### How to Use

1. The module displays the current CPU governor profile.
2. Select a new profile from the dropdown.
3. The change is applied immediately — no restart required.

### How It Works

Modularitea writes the selected governor to `/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor`, which changes the behavior of all CPU cores.

> **Note:** The governor setting is not persistent across reboots by default. You may need to re-apply it after restarting your system, or configure it persistently through your desktop environment's power settings.

---

## Settings Summary Table

| Setting | Requires Password | Persistent | Reversible |
| :--- | :--- | :--- | :--- |
| Package Cache Clean | Yes | N/A (one-time action) | No (files are deleted) |
| Swap Toggle | Yes | Yes | Yes |
| System Theme | No | Yes | Yes |
| Mirror Refresh | Yes | Yes | Yes (can refresh again) |
| DNS Provider | Yes | Yes | Yes (can switch back) |
| CPU Profile | No | No (resets on reboot) | Yes |
