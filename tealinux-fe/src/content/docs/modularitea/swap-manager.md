---
title: "Swap Manager"
description: "How to manage swap memory using Modularitea — including enabling, disabling, and understanding ZRAM swap."
category: "Modularitea"
order: 9
---

# Swap Manager

The Swap Manager in Modularitea allows you to enable or disable swap memory on your TeaLinuxOS system with a single toggle. This guide explains what swap is, when you need it, and how to manage it.

---

## What Is Swap?

**Swap** is a form of virtual memory that extends your physical RAM by using a portion of your storage (disk or RAM-backed block device) as overflow space. When your RAM is full, the operating system moves less-used data from RAM to swap, freeing RAM for active processes.

TeaLinuxOS uses **ZRAM** for swap, which is different from traditional disk-based swap:

| Feature | Traditional Swap | ZRAM Swap (TeaLinuxOS) |
| :--- | :--- | :--- |
| Location | Disk partition or swap file | Compressed block device in RAM |
| Speed | Limited by disk I/O speed | Near-RAM speed (much faster) |
| Compression | None | ZSTD compression (typically 2–3x ratio) |
| Disk usage | Consumes disk space | Zero disk space used |
| Configuration | Partition or file + fstab entry | `zram-generator.conf` + systemd |

ZRAM swap works by creating a compressed region in RAM. For example, with a 2:1 compression ratio, 4 GB of data can fit in 2 GB of physical RAM — effectively expanding your usable memory without touching the disk.

---

## When Is Swap Needed?

### You Should Enable Swap If:

- **Your system has 4 GB or less of RAM** — Swap provides critical overflow space to prevent out-of-memory crashes.
- **You run memory-intensive applications** — IDEs, browsers with many tabs, virtual machines, and compilation jobs can easily exceed physical RAM.
- **You want system stability** — Without swap, the Linux OOM (Out of Memory) killer may terminate applications abruptly when RAM is exhausted.
- **You use hibernation** — Some hibernation mechanisms require swap space.

### You Might Disable Swap If:

- **Your system has 16 GB or more of RAM** — With ample RAM, you may never need swap under normal workloads.
- **You are benchmarking** — Disabling swap ensures all operations run in physical RAM for consistent performance measurements.
- **You understand the risks** — Running without swap means the OOM killer will activate immediately when RAM is exhausted.

---

## Advantages and Disadvantages

### Advantages of Swap

| Advantage | Description |
| :--- | :--- |
| **Prevents crashes** | Applications are less likely to be killed by the OOM killer |
| **Extends usable memory** | ZRAM compression effectively increases available memory |
| **Zero disk impact** | ZRAM swap does not use disk space or cause disk wear |
| **Fast** | Compressed RAM is much faster than disk-based swap |
| **Automatic** | Once enabled, swap works transparently in the background |

### Disadvantages of Swap

| Disadvantage | Description |
| :--- | :--- |
| **CPU overhead** | Compression and decompression use CPU cycles |
| **Reduced available RAM** | The ZRAM device itself uses some RAM for the compressed store |
| **Not infinite** | ZRAM cannot compress all data equally; incompressible data fills swap quickly |

---

## How to Use the Swap Manager

### Checking Swap Status

1. Open **Modularitea**.
2. Click the **Settings** button () at the bottom of the sidebar.
3. Locate the **Swap Memory** card.
4. The card displays the current swap status: **Enabled** or **Disabled**.

### Enabling Swap

1. If swap is currently disabled, toggle the **switch to On**.
2. A **PolicyKit authentication dialog** will appear — enter your administrator password.
3. Wait for the confirmation notification.
4. Swap is now active.

### Disabling Swap

1. If swap is currently enabled, toggle the **switch to Off**.
2. Authenticate with your administrator password.
3. Wait for the confirmation notification.
4. Swap is now disabled.

---

## Verifying Swap Status

You can verify swap status from the terminal:

```bash
swapon --show
```

If swap is enabled, you will see output similar to:

```
NAME       TYPE      SIZE USED PRIO
/dev/zram0 partition   4G   0B  100
```

If swap is disabled, this command produces no output.

You can also check the ZRAM configuration file:

```bash
ls -la /etc/systemd/zram-generator.conf
```

If the file exists, swap is configured. If the file does not exist, swap is disabled.

---

## Technical Details

### How It Works

When you toggle swap in Modularitea:

1. The GUI calls the Tauri backend's `set_swap_mode` command with `"enable"` or `"disable"`.
2. The backend invokes `pkexec modularitea-swap <on|off>`.
3. The `modularitea-swap` binary (from `tealinux-modularitea-libs`) runs with root privileges.
4. **Enable:** Creates `/etc/systemd/zram-generator.conf` and starts the ZRAM device via systemd.
5. **Disable:** Removes `/etc/systemd/zram-generator.conf` and deactivates the ZRAM device.
6. The binary outputs a JSON result indicating success or failure.
7. The GUI displays the appropriate notification.

### Configuration File

When swap is enabled, the configuration file at `/etc/systemd/zram-generator.conf` defines the ZRAM parameters, including:

- The compression algorithm (ZSTD)
- The size of the ZRAM device

### Why Administrator Password Is Required

Creating and removing files in `/etc/systemd/` and manipulating block devices requires root access. Modularitea uses `pkexec` to safely escalate privileges only for this specific operation.

---

## Warnings

> ** Risk of system instability.** Disabling swap on a system with limited RAM (4 GB or less) may cause the Linux OOM (Out of Memory) killer to terminate running applications without warning. This can result in data loss. Only disable swap if you have sufficient RAM for your workload.

> ** ZRAM uses RAM.** While ZRAM is faster than disk swap, it still uses a portion of your physical RAM for the compressed block device. On very low-RAM systems (2 GB), enabling ZRAM swap may slightly reduce the RAM available for applications, though the net effect of compression is usually positive.

> ** Changes are persistent.** Enabling or disabling swap through Modularitea is a persistent change that survives reboots. You do not need to toggle swap again after restarting your system.
