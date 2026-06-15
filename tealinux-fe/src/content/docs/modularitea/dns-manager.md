---
title: "DNS Manager"
description: "How to change your DNS provider using Modularitea's DNS Manager — including available providers, usage steps, and technical details."
category: "Modularitea"
order: 7
---

# DNS Manager

The DNS Manager in Modularitea allows you to change your system's DNS (Domain Name System) provider through a simple graphical interface. This guide explains what DNS is, why you might want to change it, and how to do so.

---

## What Is DNS?

**DNS** (Domain Name System) is the internet's phone book. When you type a website address like `google.com` into your browser, your computer needs to find the actual IP address of that server (e.g., `142.250.185.14`). DNS is the system that performs this translation.

By default, your computer uses the DNS servers provided by your **Internet Service Provider (ISP)**. However, there are several reasons you might want to use a different DNS provider.

---

## Why Change DNS?

| Reason | Description |
| :--- | :--- |
| **Speed** | Third-party DNS providers like Cloudflare and Google often resolve queries faster than ISP DNS servers. |
| **Privacy** | Some DNS providers (like Cloudflare) have strict no-logging policies, reducing tracking by your ISP. |
| **Security** | Providers like Quad9 automatically block known malicious domains, providing an extra layer of protection. |
| **Reliability** | Major DNS providers have globally distributed infrastructure with near-100% uptime. |
| **Ad Blocking** | Providers like AdGuard filter ads and trackers at the DNS level. |
| **Censorship Bypass** | If your ISP blocks certain websites via DNS, switching to a public DNS can bypass these restrictions. |

---

## Available DNS Providers

Modularitea includes the following pre-configured DNS providers:

### Cloudflare (`1.1.1.1`)

| Property | Value |
| :--- | :--- |
| Primary DNS | `1.1.1.1` |
| Secondary DNS | `1.0.0.1` |
| Focus | Speed and privacy |
| Logging Policy | No personally identifiable data logged; logs purged within 24 hours |
| Website | [cloudflare.com/dns](https://cloudflare.com/dns) |

Cloudflare DNS is widely regarded as the fastest public DNS service. It is a strong choice for users who prioritize both speed and privacy.

### Google DNS (`8.8.8.8`)

| Property | Value |
| :--- | :--- |
| Primary DNS | `8.8.8.8` |
| Secondary DNS | `8.8.4.4` |
| Focus | Reliability and speed |
| Logging Policy | Temporary logs with IP anonymization |
| Website | [developers.google.com/speed/public-dns](https://developers.google.com/speed/public-dns) |

Google DNS is one of the most well-known public DNS services, offering excellent reliability and global coverage.

### Quad9 (`9.9.9.9`)

| Property | Value |
| :--- | :--- |
| Primary DNS | `9.9.9.9` |
| Secondary DNS | `149.112.112.112` |
| Focus | Security (malware blocking) |
| Logging Policy | No personally identifiable data logged |
| Website | [quad9.net](https://quad9.net) |

Quad9 blocks access to known malicious domains using threat intelligence feeds. It is the best choice for security-focused users.

### OpenDNS (`208.67.222.222`)

| Property | Value |
| :--- | :--- |
| Primary DNS | `208.67.222.222` |
| Secondary DNS | `208.67.220.220` |
| Focus | Content filtering and parental controls |
| Logging Policy | Standard logging (owned by Cisco) |
| Website | [opendns.com](https://opendns.com) |

OpenDNS (by Cisco) offers optional content filtering features, making it popular for family and educational environments.

### AdGuard (`94.140.14.14`)

| Property | Value |
| :--- | :--- |
| Primary DNS | `94.140.14.14` |
| Secondary DNS | `94.140.15.15` |
| Focus | Ad blocking and tracker blocking |
| Logging Policy | No personally identifiable data logged |
| Website | [adguard-dns.io](https://adguard-dns.io) |

AdGuard DNS blocks ads, trackers, and known malicious domains at the DNS level, providing system-wide ad blocking without installing a browser extension.

---

## How to Change DNS

### Step-by-Step

1. Open **Modularitea**.
2. Click the **Settings** button (⚙️) at the bottom of the sidebar.
3. Locate the **DNS Configuration** card in the settings grid.
4. The card displays your current DNS provider (if one of the supported providers is detected).
5. Select a new provider from the **dropdown menu**.
6. A **PolicyKit authentication dialog** will appear — enter your administrator password.
7. Wait for the confirmation notification.
8. Your DNS is now changed.

### Verifying the Change

After changing your DNS, you can verify the change by running:

```bash
cat /etc/resolv.conf
```

You should see the nameserver entries for your selected provider. For example, if you selected Cloudflare:

```
nameserver 1.1.1.1
nameserver 1.0.0.1
```

You can also test DNS resolution speed:

```bash
dig google.com
```

The `Query time` in the output indicates how fast your DNS resolved the query.

---

## Technical Details

### How It Works

When you change the DNS provider through Modularitea:

1. The GUI calls the Tauri backend's `switch_dns` command.
2. The backend invokes `pkexec modularitea-dns-changer <provider>`.
3. The `modularitea-dns-changer` binary (from `tealinux-modularitea-libs`) runs with root privileges.
4. It modifies `/etc/resolv.conf` to contain the selected provider's nameservers.
5. The binary outputs a JSON result indicating success or failure.
6. The GUI displays the appropriate notification.

### Why Administrator Password Is Required

Changing DNS requires modifying `/etc/resolv.conf`, which is a system file owned by root. The `pkexec` tool (part of PolicyKit) provides a secure way to run a specific command as root without giving the entire application root access.

> **Note:** The `pkexec` dialog is a system-level authentication prompt. Modularitea never sees or stores your password — it is handled entirely by PolicyKit.

---

## Reverting DNS Changes

To revert to your ISP's default DNS:

1. Open a terminal.
2. Remove the custom DNS entries:

```bash
sudo rm /etc/resolv.conf
```

3. Restart the NetworkManager service to regenerate the file with your ISP's DNS:

```bash
sudo systemctl restart NetworkManager
```

Alternatively, you can switch to a different DNS provider in Modularitea's settings at any time.
