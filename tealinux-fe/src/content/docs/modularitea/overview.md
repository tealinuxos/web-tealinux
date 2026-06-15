---
title: "Overview"
description: "A comprehensive overview of Modularitea — its purpose, modular architecture, and advantages over traditional system management tools."
category: "Modularitea"
order: 2
---

# Modularitea Overview

## What Is Modularitea?

**Modularitea** is a native desktop application that serves as the central hub for managing your TeaLinuxOS system. Built with **Tauri** (Rust backend + SvelteKit frontend), it provides a modern, responsive GUI that wraps around the powerful command-line tools already present in Arch Linux.

Rather than being *just* an app store or *just* a settings panel, Modularitea combines both into a unified experience. You can install a complete Python data science environment, change your DNS provider, clean your package cache, and switch your CPU governor — all from the same application.

### How Modularitea Differs from a Traditional App Store

Traditional Linux app stores (like GNOME Software or KDE Discover) focus primarily on installing individual packages one at a time. Modularitea takes a different approach:

| Aspect | Traditional App Store | Modularitea |
| :--- | :--- | :--- |
| Installation unit | Individual packages | **Profiles** — curated bundles of packages + services |
| Scope | Software only | Software + system settings + maintenance |
| Backend | PackageKit / Flatpak | Direct `pacman` integration with `pkexec` privilege escalation |
| Target audience | General users | Developers and power users on TeaLinuxOS |

For example, instead of searching for and installing `nodejs`, `npm`, `docker`, `nginx`, and `mariadb` separately, you can install the **"API Development"** profile that includes all of them along with the proper systemd services.

### How Modularitea Differs from Traditional System Settings

Desktop environment settings (like GNOME Settings or KDE System Settings) typically focus on appearance and peripheral configuration. Modularitea goes deeper:

- **DNS configuration** — directly modifies `/etc/resolv.conf`
- **Mirror optimization** — uses `reflector` to find the fastest Arch Linux mirrors
- **Swap management** — configures `zram-generator` via systemd
- **CPU governor control** — adjusts the kernel scaling governor in real time
- **Package cache management** — cleans `/var/cache/pacman/pkg/`
- **GRUB theme management** — applies boot themes with `grub-mkconfig`

### Advantages Over Terminal-Based Management

While the terminal remains available for advanced use, Modularitea offers several advantages for routine tasks:

1. **Discoverability** — All available options are visible in the UI. No need to remember command syntax or flag combinations.
2. **Safety** — Dangerous operations require explicit confirmation via PolicyKit (`pkexec`). You cannot accidentally execute a destructive command.
3. **Context** — The UI provides descriptions, status indicators, and real-time feedback that raw terminal output lacks.
4. **Consistency** — Every operation follows the same interaction pattern: select, confirm, receive feedback.

---

## The Modular Concept

Modularitea is not a monolithic application. It is designed as a collection of **independent modules**, each responsible for a specific domain of system management.

### Architecture Layers

```
┌─────────────────────────────────────────────┐
│              Modularitea GUI                │
│          (SvelteKit + TailwindCSS)          │
├─────────────────────────────────────────────┤
│            Tauri Bridge (Rust)              │
│     Commands, State Management, IPC        │
├─────────────────────────────────────────────┤
│         modularitea-libs (Rust)             │
│   Core logic: DNS, Mirror, Swap, Grub,     │
│   News Parser, Package Cache, CPU Booster  │
├─────────────────────────────────────────────┤
│        System Binaries (pkexec)             │
│   modularitea-dns-changer                  │
│   modularitea-swap                         │
│   modularitea-grub                         │
│   modularitea-pacman                       │
│   modularitea-package-cache-cleaner        │
│   modularitea-systemctl                    │
├─────────────────────────────────────────────┤
│         Operating System (Arch Linux)       │
│   pacman, systemd, reflector, resolv.conf  │
└─────────────────────────────────────────────┘
```

### Module Independence

Each module in Modularitea is self-contained:

- **DNS Manager** — Can switch DNS providers independently of any other module.
- **Mirror Manager** — Can refresh mirrors without affecting package installations.
- **Swap Manager** — Can enable or disable swap without restarting the application.
- **Profile Installer** — Can install software profiles without touching system settings.
- **GRUB Manager** — Can change boot themes without affecting the running system.

### Benefits of Modular Design

| Benefit | Description |
| :--- | :--- |
| **Easy Maintenance** | A bug in the DNS module does not affect the package installer. Fixes can be isolated and deployed independently. |
| **Easy Development** | New contributors can work on a single module without understanding the entire codebase. |
| **Easy Extension** | Adding a new feature (e.g., firewall management) means creating a new module without modifying existing ones. |
| **Consistent UX** | Every module follows the same UI patterns, so users learn one interaction model and apply it everywhere. |
| **Independent Testing** | Each module has its own test suite in the `modularitea-libs` package. |

### Profile System

One of Modularitea's most distinctive features is the **profile system**. Profiles are defined as TOML files that bundle:

- **Packages to install** (from official Arch repositories)
- **AUR packages** (community packages)
- **Services to enable** (via systemd)

Example profile (`developer_essentials.toml`):

```toml
[meta]
name = "Developer Essentials Pack"
description = "Essential developer tools for Linux."
version = "1.0.0"
author = "TeaLinuxOS Team"
category = "development"
icon = "developer_essentials.svg"

[packages]
install = ["git", "curl", "wget", "base-devel", "unzip", "zip",
           "nano", "vim", "htop", "tree", "openssh", "neovim",
           "cmake", "python", "nodejs", "npm"]
aur = ["visual-studio-code-bin", "neovim"]

[services]
enable = []
```

TeaLinuxOS ships with **28 pre-built profiles** covering categories such as:

- Web Development (JavaScript, PHP, API Development)
- Systems Programming (C++/Java, Go/Rust)
- Data Science & AI
- Design & Multimedia
- Gaming & Entertainment
- Networking & Security
- Office & Education
- DevOps (Containers, Kubernetes, Virtualization)
- And more

---

## Summary

Modularitea transforms TeaLinuxOS from a standard Arch-based distribution into a **managed developer workstation**. By combining software installation, system configuration, and maintenance into a single modular application, it eliminates the friction of day-to-day Linux administration while preserving full access to the underlying system for those who want it.
