---
title: "Mirror Manager"
description: "How to update and optimize your Arch Linux package mirrors using Modularitea's Mirror Manager."
category: "Modularitea"
order: 8
---

# Mirror Manager

The Mirror Manager in Modularitea allows you to refresh and optimize your Arch Linux package mirrors to ensure the fastest possible download speeds when installing or updating software.

---

## What Are Mirrors?

**Mirrors** are servers around the world that host copies of the Arch Linux package repositories. When you install a package with `pacman`, it downloads the package file from one of these mirrors.

The list of mirrors your system uses is stored in `/etc/pacman.d/mirrorlist`. The order matters — `pacman` tries mirrors from top to bottom and uses the first one that responds.

---

## Why Mirrors Need to Be Updated

The mirror landscape changes over time:

| Situation | Impact |
| :--- | :--- |
| A nearby mirror goes offline | Package downloads fail or time out |
| A new, faster mirror comes online | You miss out on faster downloads |
| A mirror falls behind on synchronization | You may get outdated packages |
| Your network route changes | A previously fast mirror becomes slow |

Refreshing your mirrors periodically ensures that your system always uses the fastest, most up-to-date, and most reliable servers available.

---

## How Mirror Selection Works

When you use the Mirror Manager in Modularitea, it leverages **`reflector`** — the official Arch Linux mirror ranking tool — to:

1. **Fetch** the current list of all available Arch Linux mirrors from the Arch Mirror Status page.
2. **Filter** mirrors by country (if specified), protocol, and synchronization status.
3. **Test** each mirror by measuring download speed and response time.
4. **Rank** mirrors from fastest to slowest.
5. **Write** the ranked list to `/etc/pacman.d/mirrorlist`.

This process typically takes **30–60 seconds** depending on your internet connection and the number of mirrors being tested.

---

## How to Refresh Mirrors

### Step-by-Step

1. Open **Modularitea**.
2. Click the **Settings** button gear  at the bottom of the sidebar.
3. Locate the **Mirror Settings** card in the settings grid.
4. Optionally, select a **country** to prioritize mirrors geographically close to you (e.g., "Indonesia").
5. Click the **Refresh Mirror** button.
6. Authenticate with your administrator password when prompted.
7. Wait for the process to complete — a progress indicator will be shown.
8. A success notification confirms the mirror list has been updated.

### Country Selection

Selecting a country filters mirrors to those hosted in or near that region. This is recommended because:

- Mirrors geographically closer to you are generally faster.
- It reduces the number of mirrors to test, making the process quicker.

If you do not select a country, `reflector` will test mirrors worldwide and pick the fastest overall.

---

## Verifying the Mirror Update

After refreshing mirrors, you can verify the new list:

```bash
cat /etc/pacman.d/mirrorlist
```

You should see a list of mirror URLs, ordered by speed, with the fastest at the top. The file will contain a timestamp comment indicating when it was last generated.

You can test the speed of the top mirror by syncing the database:

```bash
sudo pacman -Sy
```

If the sync completes quickly, your new mirrors are working well.

---

## Technical Details

### reflector

`reflector` is a Python script that retrieves the most recent mirror list from the Arch Linux mirror status page, filters and sorts mirrors by various criteria, and outputs a formatted mirror list.

The Modularitea Mirror Manager uses `reflector` with the following typical parameters:

- **Sort by:** Download speed
- **Protocol:** HTTPS preferred
- **Country:** As specified by the user (or worldwide if not specified)
- **Number of mirrors:** Top results saved to mirrorlist

### API Usage

Internally, Modularitea calls the `MirrorUtils` API from `modularitea-libs`:

```rust
let mirror = MirrorUtils::set_country(Some("Indonesia".to_string()));
mirror.refresh_fastest_mirror();
```

This runs `reflector` under the hood with `pkexec` for root access (since `/etc/pacman.d/mirrorlist` is owned by root).

---

## Best Practices

| Practice | Reason |
| :--- | :--- |
| Refresh mirrors monthly | Mirror performance changes over time |
| Refresh after network changes | New ISP or location may have different optimal mirrors |
| Refresh before large updates | Ensures the fastest download for major system upgrades |
| Select your country | Reduces test time and usually yields faster mirrors |

---

## Warnings

> **⚠️ Stable internet connection required.** The mirror refresh process needs to download test data from multiple servers. An unstable connection may cause the process to fail or produce inaccurate results. Ensure you have a reliable internet connection before starting.

> **⚠️ Temporary package download issues.** Immediately after refreshing mirrors, `pacman` may briefly encounter issues if a selected mirror is mid-sync. If you experience download errors after a refresh, wait a few minutes and try again, or refresh mirrors once more.

> **⚠️ Do not interrupt the process.** Interrupting a mirror refresh mid-operation may leave `/etc/pacman.d/mirrorlist` in an incomplete state. If this happens, you can manually restore it by running `reflector` from the terminal or refreshing through Modularitea again.
