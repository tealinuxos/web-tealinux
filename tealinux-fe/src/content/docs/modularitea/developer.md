---
title: "Developer Guide"
description: "Technical documentation for Modularitea developers — architecture, project structure, adding modules, contributing, and bug reporting."
category: "Modularitea"
order: 13
---

# Developer Guide

This guide is intended for developers who want to understand Modularitea's architecture, contribute code, or extend the application with new modules.

---

## Architecture

Modularitea follows a layered architecture with clear separation between the frontend UI, the Tauri bridge layer, and the backend system libraries.

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     Frontend (UI)                       │
│              SvelteKit 5 + TailwindCSS 4                │
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐  │
│  │   Home   │  │   News   │  │ Settings │  │  Grub  │  │
│  │  (Profiles)│ │  (RSS)  │  │  Panel   │  │ Themes │  │
│  └──────────┘  └──────────┘  └──────────┘  └────────┘  │
├─────────────────────────────────────────────────────────┤
│                  Tauri Bridge (Rust)                     │
│           #[tauri::command] + tauri-specta               │
│                                                         │
│  ┌──────────────┐  ┌────────────┐  ┌─────────────────┐  │
│  │  installer/   │  │   grub/    │  │    settings/    │  │
│  │  - profiler   │  │  - themes  │  │  - dns, swap,   │  │
│  │  - runner     │  │  - apply   │  │    mirror, cpu  │  │
│  └──────────────┘  └────────────┘  └─────────────────┘  │
├─────────────────────────────────────────────────────────┤
│              modularitea-libs (Rust Crate)               │
│         Core business logic + system binaries            │
│                                                         │
│  ┌───────────────────────────────────────────────────┐   │
│  │  infrastructure/                                   │  │
│  │  ├── news_parser     (RSS feed aggregation)       │  │
│  │  ├── tools_utils     (DNS, Mirror, CPU Booster)   │  │
│  │  ├── PackageCacheCleaner                          │  │
│  │  └── Swap                                         │  │
│  │                                                   │  │
│  │  bin/                (pkexec-callable binaries)    │  │
│  │  ├── modularitea-dns-changer                      │  │
│  │  ├── modularitea-swap                             │  │
│  │  ├── modularitea-grub                             │  │
│  │  ├── modularitea-pacman                           │  │
│  │  ├── modularitea-package-cache-cleaner            │  │
│  │  └── modularitea-systemctl                        │  │
│  └───────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────┤
│                 Operating System Layer                   │
│     pacman · systemd · reflector · resolv.conf · grub   │
└─────────────────────────────────────────────────────────┘
```

### Frontend

- **Framework:** SvelteKit 5 with TypeScript
- **Styling:** TailwindCSS 4 with shadcn-svelte (bits-ui) components
- **State Management:** Svelte 5 runes (`$state`, `$derived`, `$effect`)
- **API Client:** Auto-generated TypeScript bindings via `tauri-specta`
- **Build Tool:** Vite + Bun

The frontend communicates with the backend exclusively through **Tauri commands** — type-safe IPC calls generated at compile time.

### Tauri Bridge (Backend)

- **Language:** Rust
- **Framework:** Tauri 2
- **Type Bindings:** `specta` + `tauri-specta` for automatic TypeScript generation
- **HTTP Plugin:** `tauri-plugin-http` for external API calls
- **State:** Tauri managed state (e.g., `GrubManager`)

The Tauri layer acts as a bridge between the frontend and the system. It:

1. Receives commands from the frontend.
2. Calls the appropriate function in `modularitea-libs` or invokes system binaries.
3. Returns structured results (via `specta` types) to the frontend.

### modularitea-libs (Backend Library)

- **Language:** Rust
- **License:** GPL-2.0
- **Package:** `modularitea-libs`

The backend library contains all core business logic. It is designed to be **headless** — it can be used independently of the GUI for testing, CLI tools, or alternative frontends.

Key modules:

| Module | Responsibility |
| :--- | :--- |
| `infrastructure::news_parser` | RSS feed fetching and parsing with 1-hour cache |
| `infrastructure::tools_utils::DnsSwitcher` | DNS provider switching via `/etc/resolv.conf` |
| `infrastructure::tools_utils::MirrorUtils` | Mirror refresh via `reflector` |
| `infrastructure::tools_utils::CpuBooster` | CPU governor profile switching |
| `infrastructure::PackageCacheCleaner` | Pacman cache cleanup |
| `infrastructure::Swap` | ZRAM swap enable/disable |
| `domain/` | Domain models and business rules |
| `loader/` | Profile TOML file loading |
| `planner/` | Installation planning logic |
| `executor/` | Package installation execution |

### Privilege Escalation Model

Operations that require root access follow a strict pattern:

1. The GUI calls a Tauri command.
2. The Tauri command invokes a **standalone binary** via `pkexec`.
3. The binary runs with root privileges, performs the operation, and outputs JSON.
4. The Tauri command parses the JSON result and returns it to the GUI.

This design ensures that:

- The main application **never** runs as root.
- Only the specific operation gets elevated privileges.
- The user explicitly authorizes each privileged operation via PolicyKit.

---

## Project Structure

### tealinux-modularity (Main Application)

```
tealinux-modularity/
├── src/                          # SvelteKit frontend source
│   ├── app.html                  # HTML template
│   ├── routes/                   # SvelteKit routes (pages)
│   │   ├── +page.svelte          # Splash screen / onboarding
│   │   ├── +layout.svelte        # Root layout
│   │   ├── (mainpage)/           # Main app layout group
│   │   │   ├── +layout.svelte    # Sidebar + content layout
│   │   │   ├── home/             # Profile installer page
│   │   │   ├── news/             # News feed page
│   │   │   ├── sysinfo/          # System information page
│   │   │   ├── aur/              # AUR packages page
│   │   │   ├── grub/             # GRUB theme manager page
│   │   │   └── settings/         # System settings page
│   │   │       └── partials/     # Settings sub-components
│   │   │           ├── package-cache.svelte
│   │   │           ├── swap-memory.svelte
│   │   │           ├── system-theme.svelte
│   │   │           ├── mirror-settings.svelte
│   │   │           ├── dns-configuration.svelte
│   │   │           └── cpu-performance.svelte
│   │   └── partials/             # Shared page partials
│   └── lib/                      # Shared library code
│       ├── commands.ts           # Auto-generated Tauri bindings
│       ├── components/           # Svelte components
│       │   ├── TeaBar.svelte     # Sidebar navigation
│       │   ├── teabar-items.ts   # Sidebar menu items
│       │   └── ui/               # shadcn-svelte UI components
│       ├── data/                 # Static data files
│       ├── hooks/                # Svelte hooks
│       ├── services/             # Business logic services
│       ├── state/                # Application state
│       ├── stores/               # Svelte stores
│       ├── types/                # TypeScript type definitions
│       └── utils/                # Utility functions
├── src-tauri/                    # Tauri Rust backend
│   ├── src/
│   │   ├── main.rs               # Entry point
│   │   ├── lib.rs                # Tauri app setup + command registration
│   │   ├── pkexec_args.rs        # pkexec helper utilities
│   │   ├── installer/            # Package installation module
│   │   │   ├── backend_runner.rs # pacman/pkexec operations
│   │   │   └── profiler.rs       # TOML profile loading
│   │   ├── grub/                 # GRUB theme module
│   │   │   ├── initialization.rs # GrubManager setup
│   │   │   └── command.rs        # Tauri commands for GRUB
│   │   ├── settings/             # System settings module
│   │   │   ├── command.rs        # Tauri commands for settings
│   │   │   ├── models.rs         # DnsProvider, CpuProfile, SwapMode
│   │   │   ├── mod.rs            # Settings logic + additional commands
│   │   │   └── mirror_countries.rs # Allowed reflector countries
│   │   ├── sysinfo/              # System information module
│   │   ├── splash/               # Splash screen logic
│   │   └── utils/                # Shared backend utilities
│   ├── Cargo.toml                # Rust dependencies
│   ├── tauri.conf.json           # Tauri configuration
│   └── capabilities/             # Tauri capability definitions
├── profiles/                     # Software profile TOML files
│   ├── developer_essentials.toml
│   ├── javascript.toml
│   ├── python.toml
│   ├── data_science.toml
│   ├── gaming.toml
│   └── ... (28 total)
├── polkit/                       # PolicyKit action definitions
├── static/                       # Static assets (images, fonts)
├── package.json                  # Frontend dependencies
├── svelte.config.js              # SvelteKit configuration
├── vite.config.js                # Vite configuration
└── tsconfig.json                 # TypeScript configuration
```

### tealinux-modularitea-libs (Backend Library)

```
tealinux-modularitea-libs/
├── src/
│   ├── lib.rs                    # Library root
│   ├── config.rs                 # Configuration constants
│   ├── error.rs                  # Error types
│   ├── domain/                   # Domain models
│   ├── infrastructure/           # System interaction layer
│   │   ├── news_parser.rs        # RSS feed parser
│   │   ├── tools_utils/          # DNS, Mirror, CPU utilities
│   │   ├── mod.rs                # Module exports
│   │   └── ...
│   ├── loader/                   # TOML profile loader
│   ├── planner/                  # Installation planner
│   ├── executor/                 # Package executor
│   ├── helper/                   # Helper utilities
│   ├── privilege/                # Privilege escalation
│   └── bin/                      # Standalone binaries
│       ├── modularitea-dns-changer.rs
│       ├── modularitea-swap.rs
│       ├── modularitea-grub.rs
│       ├── modularitea-pacman.rs
│       ├── modularitea-package-cache-cleaner.rs
│       ├── modularitea-systemctl.rs
│       ├── modularitea-settings.rs
│       └── modularitea-profile-installer.rs
├── data/                         # Bundled data (grub themes, etc.)
├── tests/                        # Integration tests
├── testfiles/                    # Test fixtures
├── Cargo.toml                    # Rust dependencies
├── PKGBUILD                      # Arch Linux package build script
└── install.sh                    # Post-install script
```

---

## Adding a New Module

To add a new feature module to Modularitea, follow these steps:

### Step 1: Create the Backend Logic

If your module needs root access, create a new binary in `tealinux-modularitea-libs`:

1. Create `src/bin/modularitea-<feature>.rs` with the main logic.
2. Add the binary target to `Cargo.toml`:

```toml
[[bin]]
name = "modularitea-<feature>"
path = "src/bin/modularitea-<feature>.rs"
```

3. Ensure the binary outputs JSON for programmatic consumption.
4. Add any shared logic to the `infrastructure/` module.

### Step 2: Create the Tauri Commands

In `tealinux-modularity/src-tauri/src/`:

1. Create a new module directory (e.g., `src/<feature>/`).
2. Implement `command.rs` with `#[tauri::command]` and `#[specta::specta]` functions.
3. Create `mod.rs` to expose the module.
4. Register the commands in `lib.rs`:

```rust
mod <feature>;

// In the Builder::new().commands() call:
<feature>::command::your_command_name,
```

### Step 3: Create the Frontend UI

In `tealinux-modularity/src/`:

1. Create a new route directory: `routes/(mainpage)/<feature>/`
2. Create `+page.svelte` with the page UI.
3. Add the page to the sidebar in `lib/components/teabar-items.ts`:

```typescript
{
    href: '/<feature>',
    label: 'Feature Name',
    icon: YourIcon
}
```

### Step 4: Generate TypeScript Bindings

Run the application in development mode to auto-generate the TypeScript bindings:

```bash
bunx tauri dev
```

The `src/lib/commands.ts` file will be automatically updated with type-safe bindings for your new commands.

### Step 5: Testing

1. Test the backend binary independently:

```bash
cargo test -p modularitea-libs
```

2. Test the Tauri commands in development mode:

```bash
TEALINUX_BUILD=dev bunx tauri dev
```

3. Verify the UI works end-to-end.

---

## Development Setup

### Prerequisites

| Tool | Version | Purpose |
| :--- | :--- | :--- |
| Rust | Latest stable | Backend compilation |
| Bun | Latest | Frontend package management and dev server |
| Node.js | 18+ | Alternative to Bun (optional) |
| GTK4 + WebKitGTK | System packages | Tauri rendering |
| pkexec | System package | Privilege escalation |

### Building from Source

```bash
# Clone the repository
git clone git@github.com:tealinuxos/tealinux-modularity.git
cd tealinux-modularity

# Install frontend dependencies
bun install

# Run in development mode
TEALINUX_BUILD=dev bunx tauri dev

# Build for production
bunx tauri build
```

### Environment Variables

| Variable | Value | Purpose |
| :--- | :--- | :--- |
| `TEALINUX_BUILD` | `dev` | Prevents accidental system changes during development |
| `TEALINUX_BUILD` | `prod` | Enables all features for production |
| `TEALINUX_BUILD_SHOW_SPLASH` | `true/false` | Force show/hide splash screen |
| `RUST_LOG` | `debug` | Enable verbose backend logging |

---

## Contributing

### Workflow

1. **Fork** the repository on GitHub.
2. **Clone** your fork locally.
3. **Create a branch** for your feature or fix: `git checkout -b feature/my-feature`
4. **Make changes** following the code standards below.
5. **Test** your changes thoroughly.
6. **Commit** with clear, descriptive messages. Keep commits small and focused.
7. **Push** your branch and create a **Pull Request**.

### Code Standards

#### Rust (Backend)

- Follow standard Rust formatting (`cargo fmt`).
- Run `cargo clippy` and address all warnings.
- Document public functions with `///` doc comments.
- Handle all errors explicitly — no `.unwrap()` in production code.
- All `pkexec` binaries must output JSON.

#### TypeScript/Svelte (Frontend)

- Follow the project's Prettier configuration.
- Run `bun run lint` before committing.
- Use TypeScript types from the auto-generated `commands.ts`.
- Use Svelte 5 runes (`$state`, `$derived`, `$effect`) instead of legacy stores.

#### Commit Messages

- Use clear, descriptive commit messages.
- Keep commits small — one logical change per commit.
- Prefix with the affected module: `[settings] Fix DNS provider detection`, `[installer] Add retry logic`.

### AI-Assisted Code Policy

The use of AI for producing or modifying code is allowed under strict rules:

- **Human review required** — Every AI-assisted change must be reviewed and approved by a human maintainer.
- **No silent behavioral changes** — AI-generated code must not alter existing behavior without explicit documentation.
- **High-risk areas must be done by human** — Privilege escalation, system file modifications, and security-critical code.
- **Small refactors only** — Large architectural changes require a design document and team sign-off.
- **Revert policy** — Code merged without complying with these rules will be reverted.

---

## Bug Reporting

### What to Include

When reporting a bug, include the following information:

| Information | How to Get It |
| :--- | :--- |
| TeaLinuxOS version | `cat /etc/os-release` |
| Modularitea version | `pacman -Q tealinux-modularity` |
| Libs version | `pacman -Q tealinux-modularitea-libs` |
| Desktop environment | `echo $XDG_CURRENT_DESKTOP` |
| Display protocol | `echo $XDG_SESSION_TYPE` |
| Steps to reproduce | Detailed step-by-step instructions |
| Expected behavior | What you expected to happen |
| Actual behavior | What actually happened |
| Error messages | Terminal output or UI error text |
| System logs | `journalctl -b \| grep modularitea` |

### How to Create an Issue

1. Go to the [TeaLinuxOS GitHub repository](https://github.com/tealinuxos).
2. Navigate to the relevant project (`tealinux-modularity` or `tealinux-modularitea-libs`).
3. Click **"New Issue"**.
4. Select the **Bug Report** template if available.
5. Fill in all the information listed above.
6. Add relevant labels (e.g., `bug`, `settings`, `installer`).
7. Submit the issue.

### Security Issues

If you discover a security vulnerability (e.g., privilege escalation bypass, data exposure), **do not** create a public issue. Instead, contact the TeaLinuxOS security team directly through the DOSCOM (Dinus Open Source Community) channels.
