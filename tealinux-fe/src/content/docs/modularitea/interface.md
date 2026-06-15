---
title: "Interface Guide"
description: "A detailed walkthrough of the Modularitea user interface — sidebar, dashboard, header, notifications, and navigation flow."
category: "Modularitea"
order: 4
---

# Interface Guide

Modularitea features a clean, modern interface designed with discoverability and ease of use in mind. This guide walks you through every section of the application.

---

## Application Layout

The Modularitea interface is divided into three main areas:

```
┌──────────┬──────────────────────────────────┐
│          │          Header / Trigger         │
│          ├──────────────────────────────────┤
│ Sidebar  │                                  │
│ (TeaBar) │        Main Content Area         │
│          │                                  │
│          │     (Dashboard / Settings /       │
│          │      Profile Detail / etc.)       │
│          │                                  │
├──────────┤                                  │
│ Settings │                                  │
│  Button  │                                  │
└──────────┴──────────────────────────────────┘
```

---

## Sidebar (TeaBar)

The sidebar — internally called **TeaBar** — is the primary navigation element of Modularitea. It is always visible on the left side of the application.

### Sidebar Header

At the top of the sidebar, you will find:

- The **Modularitea logo** — a distinctive green icon.
- The application name **"Modularitea"** displayed in TeaLinuxOS brand green (`#26A768`).

### Navigation Menu

The sidebar contains the following navigation items:

| Menu Item | Icon | Description |
| :--- | :--- | :--- |
| **Home** | 🏠 House | The main dashboard showing all available software profiles |
| **Latest News** | 📈 Trending | Aggregated open-source news from popular RSS feeds |
| **System Information** | ℹ️ Info | Hardware and software details about your machine |
| **AUR Packages** | 📦 Package | Browse and manage Arch User Repository packages |
| **Boot Theme** | 🎚️ Sliders | Preview and apply GRUB boot themes |

### Settings Button

At the bottom of the sidebar, separated by a horizontal divider, is the **Settings** button (⚙️ Cog icon). This provides access to all system configuration options.

### Responsive Behavior

The sidebar automatically adapts to your window size:

- **Expanded mode** (window width > 1024px): Full sidebar with icons and text labels.
- **Collapsed mode** (window width ≤ 1024px): Compact sidebar showing only icons, saving screen space.

You can manually toggle the sidebar state using the **Sidebar Trigger** button at the top of the main content area.

---

## Home Dashboard

The Home page is the default view when you open Modularitea. It displays:

### Hero Section

A prominent banner at the top showing:

- The welcome message and brief description of Modularitea.
- The total number of available software profiles.

### Profile Grid

Below the hero section, you will find a grid of **Profile Cards**. Each card represents an installable software bundle and displays:

| Element | Description |
| :--- | :--- |
| **Title** | The name of the profile (e.g., "Developer Essentials Pack") |
| **Description** | A brief summary of what the profile provides |
| **Category** | The category tag (e.g., `development`, `multimedia`, `gaming`) |
| **Package Count** | The total number of packages included in the profile |
| **Package List** | A preview of the included packages |
| **Services Count** | Number of systemd services that will be enabled |
| **Install Button** | Action button to install the entire profile |

### Profile Install States

Each profile card reflects its current status:

| State | Appearance | Description |
| :--- | :--- | :--- |
| `idle` | Default | Profile is available for installation |
| `installing` | Loading animation | Installation is in progress |
| `success` | Success indicator | Profile was installed successfully |
| `error` | Error indicator | Installation encountered an error |

---

## System Information Page

The System Information page provides a detailed overview of your hardware and software:

### Computer Information

- **Processor** — CPU model and specifications
- **Memory** — Total RAM available
- **Operating System** — TeaLinuxOS version
- **Kernel Version** — Linux kernel version
- **Username** — Current logged-in user(s)

### Display Information

- **Monitor Name** — Connected display(s)
- **Graphics Cards** — Detected GPU(s)
- **Display Protocol** — Wayland or X11
- **Window Manager** — Active compositor/WM

### Audio Information

- **Audio Devices** — List of detected audio output/input devices

---

## News Feed

The Latest News page aggregates articles from popular open-source and technology news sources via RSS feeds. Each news item displays:

- **Title** — The article headline
- **Description** — A brief summary of the article
- **Thumbnail** — An article preview image (when available)
- **Source URL** — Click to open the full article in your browser

> **Note:** News data is cached for up to **1 hour** to improve performance and reduce network requests.

---

## Boot Theme Manager

The Boot Theme page allows you to customize your GRUB bootloader appearance:

- **Theme Gallery** — Browse all available GRUB themes with preview images.
- **Theme Details** — View the theme name, author, version, and GitHub URL.
- **Apply Button** — Apply a theme with a single click (requires administrator password).

Available themes include the default TeaLinuxOS theme ("Lorem Loader") and community themes.

---

## Notifications

Modularitea uses a toast notification system (powered by `svelte-sonner`) to communicate the result of user actions:

### Success Notifications

- Appear as a green toast at the bottom of the screen.
- Confirm that an action completed successfully.
- Example: *"DNS provider switched to Cloudflare successfully."*

### Error Notifications

- Appear as a red toast with an error description.
- Include actionable information about what went wrong.
- Example: *"Failed to refresh mirror: network timeout."*

### Confirmation Dialogs

Some actions that require administrator privileges will trigger a **PolicyKit authentication dialog**:

1. The system displays a password prompt.
2. Enter your administrator password.
3. The action proceeds if authentication is successful.
4. If you cancel the dialog, the action is aborted safely.

---

## First-Time User Flow

Here is the typical experience when opening Modularitea for the first time:

1. **Launch** — Open Modularitea from the application menu or terminal.
2. **Splash Screen** — The onboarding wizard introduces key features (first launch only).
3. **Complete Onboarding** — Follow the splash screen prompts to finish initial setup.
4. **Home Dashboard** — The main interface loads, showing all available software profiles.
5. **Browse Profiles** — Scroll through the profile grid to find software bundles that match your needs.
6. **Install a Profile** — Click the install button on a profile card, authenticate with your password, and wait for installation to complete.
7. **Explore Settings** — Click the Settings button in the sidebar to configure DNS, mirrors, swap, and other system options.
8. **Check System Info** — Navigate to System Information to verify your hardware configuration.

---

## Keyboard Navigation

Modularitea supports standard keyboard navigation:

| Key | Action |
| :--- | :--- |
| `Tab` | Move focus to the next interactive element |
| `Shift + Tab` | Move focus to the previous interactive element |
| `Enter` | Activate the focused button or link |
| `Escape` | Close modal dialogs |
