# Sesh - A Session Management Tool

**Sesh** is a terminal-based session management tool that integrates with `tmux` to help users organize and switch between their sessions quickly. It utilizes a database to store session directories and provides a simple CLI for interacting with sessions.

## Features

- Stores session paths in a persistent database.
- Allows easy management of sessions with commands to add, remove, and search for session directories.
- Automatically starts a new `tmux` session or attaches to an existing one based on the selected directory.
- Uses `fzf` for interactive searching of session paths.

## Installation

This project is built using **Nix** and can be installed via a Nix shell or directly as a package.

### Prerequisites

Make sure you have the following installed:

- [Nix](https://nixos.org/)
- [fzf](https://github.com/junegunn/fzf) (for fuzzy searching session directories)
- [tmux](https://github.com/tmux/tmux) (for session management)

### Installation via Nix

1. Clone the repository:

   ```bash
   git clone https://github.com/NewDawn0/sesh.git
   cd sesh
   ```

2. Build and install the package using Nix:

   ```bash
   nix build
   ```

3. To start a session with `nix` shell:

   ```bash
   nix develop .#
   ```

4. Optionally, you can install it globally if needed:

   ```nix
   {
       # Add the input
       inputs.sesh.url = "github:NewDawn0/sesh";

       # Add the overlay
       import nixpkgs = {
            # Other pgks setup
           overlays = [
               sesh.overlays.default
           ];
       }
   }

   ```

## Usage

Once you have `sesh` set up, you can use it to manage your `tmux` sessions based on directories:

### 1. Start or Attach to a Session

To start a new session or attach to an existing one:

```bash
sesh
```

### 2. Add a Directory to the Session List

To add a directory to the session list:

```bash
sesh add /path/to/directory
```

### 3. Remove a Directory from the Session List

To remove a directory from the session list:

```bash
sesh rm /path/to/directory
```

### 4. Find and Select a Directory Using `fzf`

The program allows you to interactively select directories using `fzf`. When you run `sesh`, it will show a list of available directories that you can choose from.

### 5. Custom Shell Hook

The tool is configured to use a `shellHook` script that defines the `sesh` function. This function starts a tmux session in the directory provided by the script. It’s automatically sourced when you enter the Nix shell.

## Configuration

### File Paths

The session database is stored in `~/.cache/sesh.db` by default. This file contains the paths to the session directories and their corresponding tmux sessions.

If you'd like to change the location of the database file, you can modify the `file` field in the `DB` struct inside the `main.go` file.

### Customizing the Shell Hook

If you want to modify the behavior of the `sesh` function (e.g., changing the default session name or behavior), you can edit the `hooks/shellHook` file. This script is sourced whenever you enter the Nix shell environment.

## Example Workflow

1. **Starting a Session**:

   - Run `nix develop` to enter the environment and automatically source the `shellHook`.
   - The `sesh` function will be available, allowing you to start or attach to a tmux session.

2. **Adding a Directory**:

   - Add a directory to the session list by running `sesh add /path/to/your/project`.

3. **Finding a Directory**:

   - Run `sesh` to search through the stored directories interactively with `fzf`.

4. **Removing a Directory**:
   - Run `sesh rm /path/to/your/project` to remove a session directory from the list.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
