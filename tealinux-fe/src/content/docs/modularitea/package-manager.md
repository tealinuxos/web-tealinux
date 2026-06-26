---
title: "Package Manager"
description: "How to browse, install, remove, and update software using Modularitea's profile-based package management system."
category: "Modularitea"
order: 5
---

# Package Manager

Modularitea's package management system takes a **profile-based approach** to software installation. Instead of searching for individual packages, you install curated **software profiles** — pre-configured bundles of packages and services designed for specific use cases.

---

## App Store (Profile Installer)

### Browsing Profiles

When you open Modularitea, the **Home** page displays all available software profiles in a grid layout. Each profile card shows:

- The profile name and description
- The category (e.g., `development`, `multimedia`, `gaming`)
- The number of packages included
- A preview of the included package names

You can visually scan the grid to find profiles that match your needs.

### Available Profile Categories

TeaLinuxOS ships with **28 pre-built profiles** across the following categories:

| Category | Example Profiles | Description |
| :--- | :--- | :--- |
| **Development** | Developer Essentials, JavaScript, PHP, Python, C++/Java, Go/Rust | Programming languages, tools, and frameworks |
| **DevOps** | Container, Kubernetes & Cloud, Virtualization | Docker, Podman, k8s tools, VirtualBox, QEMU |
| **Data & AI** | Data Science, AI Pack | Jupyter, R, Python ML libraries |
| **Design** | Design | Inkscape, GIMP, and UI/UX tools |
| **Multimedia** | Multimedia, Audio & Video | Media players, editors, and recording tools |
| **Networking** | Network & Security, Internet | Security tools, browsers, communication apps |
| **Productivity** | Office, Student & Teacher | LibreOffice, educational tools |
| **System** | Backup & Optimization, Utility & Driver | System maintenance and hardware support |
| **Gaming** | Gaming | Steam, Lutris, and gaming utilities |
| **Other** | Fonts, Islamic Education, Mobile, Embedded & IoT | Specialized tool collections |

### Viewing Profile Details

Click on a profile card to see its full details, including:

- **Complete package list** — All official repository packages that will be installed
- **AUR packages** — Community packages from the Arch User Repository
- **Services** — Systemd services that will be enabled after installation

### Installing a Profile

To install a software profile:

1. **Find the profile** you want on the Home page.
2. **Click the Install button** on the profile card.
3. **Authenticate** — A PolicyKit dialog will appear asking for your administrator password.
4. **Wait for completion** — The profile card will show a loading animation during installation.
5. **Verify** — Once complete, a success notification will confirm the installation.

During installation, Modularitea:

1. Updates the package database (`pacman -Sy`).
2. Installs all official packages (`pacman -S <packages>`).
3. Enables any required systemd services (`systemctl enable <service>`).
4. Reports the result back to the UI.

> **Note:** AUR packages are currently listed in profiles but their automatic installation is a work in progress. A warning will appear for profiles containing AUR packages.

### Removing a Profile

To uninstall a previously installed profile:

1. Navigate to the profile's detail view.
2. Click the **Uninstall** button.
3. Authenticate with your administrator password.
4. Modularitea will remove the profile's packages and disable its services.

> **Important:** When removing a profile, Modularitea uses `pacman -R` by default. If a package is a dependency of another installed package, it will be skipped to prevent system breakage. The `force` option (`pacman -Rdd`) is available but should be used with caution.

### Updating Packages

To update all installed packages on your system:

1. Navigate to the **Settings** page.
2. Use the system update functionality (see [Maintenance](./maintenance) for details).

Modularitea uses `pacman -Syu` to perform a full system upgrade, which updates all installed packages — including those installed through profiles.

---

## Package Management Operations

Behind the scenes, Modularitea provides the following package management commands through its Tauri backend:

| Operation | Backend Command | Description |
| :--- | :--- | :--- |
| Install packages | `pkexec pacman -S --noconfirm <packages>` | Install one or more packages |
| Remove packages | `pkexec pacman -R --noconfirm <packages>` | Remove packages (respects dependencies) |
| Force remove | `pkexec pacman -Rdd --noconfirm <packages>` | Remove packages (ignores dependencies) |
| Update database | `pkexec pacman -Sy` | Synchronize the package database |
| Check installed | `pacman -Q <package>` | Check if a package is installed |
| Install profile | Install packages + enable services | Full profile installation workflow |
| Uninstall profile | Remove packages + disable services | Full profile removal workflow |

---

## Integration with pacman

Modularitea is a **frontend** for `pacman` — the standard Arch Linux package manager. It does not replace `pacman` or use a separate package database. This means:

- **Packages installed via Modularitea are visible in `pacman`** — You can run `pacman -Q` to see them.
- **Packages installed via terminal are visible in Modularitea** — The profile system checks installed packages in real time.
- **System updates via `pacman -Syu` include Modularitea-installed packages** — There is no separate update channel.

### How Privilege Escalation Works

Modularitea does **not** run as root. Instead, it uses PolicyKit (`pkexec`) to escalate privileges only when needed:

1. You click an action that requires root access (e.g., install packages).
2. Modularitea calls the appropriate backend binary via `pkexec`.
3. PolicyKit displays a system authentication dialog.
4. You enter your password.
5. The binary executes with root privileges and returns the result.
6. Modularitea displays the result in the UI.

This approach is safer than running the entire application as root, because only the specific operation runs with elevated privileges.

---


## Example: Installing a Web Development Environment

Here is a complete walkthrough of installing the JavaScript development profile:

1. Open Modularitea.
2. On the Home page, find the **"JavaScript"** profile card.
3. Review the included packages: `nodejs`, `npm`, `bun`, and related tools.
4. Click **Install**.
5. Enter your administrator password when prompted.
6. Wait for the installation to complete (the card shows a progress indicator).
7. Once complete, open a terminal and verify:

```bash
node --version
npm --version
```

Both commands should now return version numbers, confirming the successful installation.
