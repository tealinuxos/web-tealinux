---
title: "Installation"
description: "How to install, verify, and launch Modularitea on TeaLinuxOS."
category: "Modularitea"
order: 3
---

# Installation

## Pre-Installed on TeaLinuxOS

Modularitea comes **pre-installed** on all editions of TeaLinuxOS. If you are running a standard TeaLinuxOS installation (Cosmic or KDE edition), Modularitea is already available on your system and ready to use.

No additional installation steps are required for most users.

---

## Verifying Modularitea Is Installed

To confirm that Modularitea is available on your system, open a terminal and run:

```bash
which tealinux-modularity
```
---

## Launching Modularitea

There are several ways to open Modularitea:

### From the Application Menu

1. Open your desktop's application launcher (press the `Super` key or click the application menu).
2. Search for **"Modularitea"**.
3. Click the Modularitea icon to launch the application.

### From the Terminal

```bash
tealinux-modularity
```

### First Launch

When you launch Modularitea for the first time, you will see a **splash screen** (onboarding wizard). This splash screen introduces you to the key features and allows initial configuration. After completing the onboarding, subsequent launches will skip directly to the main interface.

> **Note:** The splash screen behavior can be controlled with the environment variable `TEALINUX_BUILD_SHOW_SPLASH=true` to force it to appear, or it can be reset by deleting the configuration file at `~/.config/tealinux-modularity`.

---



## Dependencies

Modularitea relies on the following system components. These are automatically installed as dependencies when you install the `tealinux-modularity` package:

### Core Dependencies

| Dependency | Purpose |
| :--- | :--- |
| `tealinux-modularitea-libs` | Backend library containing system management binaries |
| `gtk4` | GTK4 libraries for the native window rendering |
| `webkit2gtk` | Web rendering engine used by Tauri |
| `pkexec` (PolicyKit) | Privilege escalation for administrative tasks |

### Feature-Specific Dependencies

| Dependency | Feature | Purpose |
| :--- | :--- | :--- |
| `reflector` | Mirror Manager | Finds and ranks the fastest Arch Linux mirrors |
| `grub` | Boot Theme Manager | Required for applying GRUB themes via `grub-mkconfig` |
| `systemd` | Swap Manager | Manages `zram-generator` for swap configuration |

### Backend Binaries

The `tealinux-modularitea-libs` package installs the following helper binaries to `/usr/bin/`:

| Binary | Purpose |
| :--- | :--- |
| `modularitea-dns-changer` | Switches DNS provider (runs via `pkexec`) |
| `modularitea-swap` | Enables/disables swap (runs via `pkexec`) |
| `modularitea-grub` | Applies GRUB themes (runs via `pkexec`) |
| `modularitea-package-cache-cleaner` | Cleans the pacman package cache |
| `modularitea-systemctl` | Manages systemd services (runs via `pkexec`) |

> **Note:** These binaries are not intended to be used directly by end users. They are invoked internally by the Modularitea GUI through `pkexec` for proper privilege escalation.

---

## System Requirements

| Component | Minimum | Recommended |
| :--- | :--- | :--- |
| CPU | 64-bit x86_64 | 2 GHz Dual Core or better |
| RAM | 4 GB | 8 GB |
| Storage | 20 GB | 40 GB |
| Display | 800×600 | 1280×768 or higher |
| Network | Required for package installation and mirror updates | Stable broadband connection |

Modularitea itself is lightweight — the application window requires approximately **80–120 MB of RAM** during operation. However, operations like installing packages or refreshing mirrors will require an active internet connection.

---

## Environment Variables

Modularitea supports the following environment variables for advanced configuration:

| Variable | Values | Description |
| :--- | :--- | :--- |
| `TEALINUX_BUILD` | `dev` or `prod` | Controls whether dangerous operations (changing GRUB theme, disabling swap) are available. Set to `dev` during development to avoid accidental system changes. |
| `TEALINUX_BUILD_SHOW_SPLASH` | `true` or `false` | Forces the splash screen to show or hide, overriding the stored configuration in `~/.config/tealinux-modularity`. |
