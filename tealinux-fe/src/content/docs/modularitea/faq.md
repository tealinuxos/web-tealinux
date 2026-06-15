---
title: "FAQ"
description: "Frequently asked questions about Modularitea — covering safety, compatibility, features, and common concerns."
category: "Modularitea"
order: 12
---

# Frequently Asked Questions

## General

### 1. What is Modularitea?

Modularitea is a graphical system management application for TeaLinuxOS. It combines the functionality of an app store and system settings panel into a single modern interface, allowing you to install software, configure system settings, and perform maintenance without using the terminal.

### 2. Is Modularitea safe to use?

Yes. Modularitea is developed by the TeaLinuxOS team and is part of the official TeaLinuxOS distribution. It does not run with root privileges — administrative operations are performed through PolicyKit (`pkexec`), which requires your explicit password confirmation before making any system changes.

### 3. Can I still use the terminal alongside Modularitea?

Absolutely. Modularitea is a **frontend** for existing system tools (`pacman`, `reflector`, `systemd`, etc.). Everything it does can also be done via the terminal, and changes made through either method are fully compatible. You can freely mix GUI and terminal workflows.

### 4. Does Modularitea replace pacman?

No. Modularitea uses `pacman` under the hood. It does not have its own package database or package format. Packages installed through Modularitea are standard Arch Linux packages and are visible in `pacman -Q`.

### 5. Is Modularitea available on all editions of TeaLinuxOS?

Yes. Modularitea is included in both the **Cosmic** and **KDE** editions of TeaLinuxOS. It is installed by default during the operating system installation.

---

## Privacy & Security

### 6. Does Modularitea collect any data?

No. Modularitea does not collect telemetry, usage data, or personal information. The news feed feature fetches RSS data from public open-source news sites, but this is a standard HTTP request — no identifying information is sent.

### 7. Why do some actions require my password?

Certain operations — such as installing packages, changing DNS, or managing swap — modify system files that are protected by root permissions. The password prompt is provided by PolicyKit, a standard Linux security framework. Modularitea never sees or stores your password.

### 8. Can Modularitea damage my system?

Modularitea is designed with safety in mind. However, as with any system management tool, some operations (like cleaning the package cache or disabling swap) have real effects on your system. Each action includes appropriate warnings and requires confirmation before proceeding.

---

## Features

### 9. What are profiles?

Profiles are curated bundles of software packages and system services defined in TOML files. Instead of installing packages one by one, you can install an entire development environment (e.g., "JavaScript Development" or "Data Science") with a single click.

### 10. Can I create my own profiles?

Yes. Profiles are simple `.toml` files stored in the `profiles/` directory. You can create custom profiles by following the TOML format documented in the [Package Manager](./package-manager#profile-toml-format) page.

### 11. Does Modularitea support AUR packages?

Profiles can list AUR packages, but automatic AUR installation is currently a work in progress. Modularitea will display a warning for profiles containing AUR packages. In the meantime, you can install AUR packages manually using `yay` or another AUR helper.

### 12. Can I change my DNS back to the ISP default?

Yes. You can either select a different DNS provider in Modularitea or revert to your ISP's DNS by running:

```bash
sudo rm /etc/resolv.conf
sudo systemctl restart NetworkManager
```

### 13. What happens if I disable swap?

Disabling swap means your system loses its memory overflow protection. If your RAM fills up completely, the Linux OOM (Out of Memory) killer will terminate applications. This is generally safe if you have ample RAM (16 GB+), but risky on systems with 4 GB or less.

### 14. Does mirror refresh affect installed packages?

No. Refreshing mirrors only changes which servers `pacman` downloads packages from. It does not modify, update, or remove any installed packages.

### 15. What is ZRAM and why does TeaLinuxOS use it?

ZRAM is a Linux kernel feature that creates a compressed block device in RAM. It is used as swap space and is much faster than traditional disk-based swap. TeaLinuxOS uses ZRAM because it provides memory overflow protection without consuming disk space or causing disk wear.

---

## Troubleshooting

### 16. Modularitea shows a blank screen. What should I do?

Try the following:

1. Close Modularitea and relaunch from the terminal: `tealinux-modularity`
2. Check for errors in the terminal output.
3. If using Wayland, try: `GDK_BACKEND=x11 tealinux-modularity`
4. Reset the configuration: `rm ~/.config/tealinux-modularity`

### 17. I cancelled the password dialog. How do I retry?

Simply click the same button again. Modularitea will re-prompt for authentication. Cancelling the dialog does not break anything — it simply aborts the requested operation.

### 18. An installation failed partway through. Is my system broken?

Most likely not. `pacman` is transactional — if a package fails to install, it rolls back that specific package. Your system should remain in a consistent state. You can verify by running:

```bash
pacman -Dk
```

If this shows no errors, your package database is healthy.

### 19. How do I report a bug in Modularitea?

Visit the TeaLinuxOS GitHub repository and create a new issue. Include:

- Your TeaLinuxOS version and edition
- Steps to reproduce the issue
- Any error messages from the terminal or UI
- Output of `pacman -Q | grep modularitea`

See the [Developer Guide](./developer#bug-reporting) for detailed instructions.

---

## Compatibility

### 20. Can I use Modularitea on non-TeaLinuxOS distributions?

Modularitea is designed specifically for TeaLinuxOS and depends on the TeaLinuxOS repository and package infrastructure. It is not officially supported on other distributions, though technically it could work on any Arch-based distribution with the appropriate packages installed.

### 21. Does Modularitea work on Wayland and X11?

Yes. Modularitea is built with Tauri (which uses WebKitGTK), and supports both Wayland and X11 display protocols. The application automatically detects your display protocol and adapts accordingly.

### 22. What is the minimum screen resolution for Modularitea?

The minimum supported window size is **800×600 pixels**. The recommended resolution is **1280×768 or higher** for the best experience. The sidebar automatically collapses at narrow window widths for optimal space usage.

---

## Development

### 23. Is Modularitea open source?

Yes. Modularitea is open source. The main application (`tealinux-modularity`) is licensed under the **MIT License**, and the backend library (`tealinux-modularitea-libs`) is licensed under **GPL-2.0**.

### 24. What technologies does Modularitea use?

- **Frontend:** SvelteKit 5 + TailwindCSS 4
- **Desktop Framework:** Tauri 2 (Rust)
- **Backend Library:** Rust (modularitea-libs)
- **Type Safety:** TypeScript bindings auto-generated via `tauri-specta`
- **UI Components:** shadcn-svelte (bits-ui)
- **Package Runtime:** Bun

### 25. How can I contribute to Modularitea?

See the [Developer Guide](./developer) for detailed contribution instructions including how to set up the development environment, code standards, and the pull request process.
