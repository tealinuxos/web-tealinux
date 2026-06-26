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


Modularitea transforms TeaLinuxOS from a standard Arch-based distribution into a **managed developer workstation**. By combining software installation, system configuration, and maintenance into a single modular application, it eliminates the friction of day-to-day Linux administration while preserving full access to the underlying system for those who want it.
