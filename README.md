# Go-Git-Good

Go-Git-Good is a research project written in Go to explore Git internals by implementing core Git commands from scratch, ranging from low-level plumbing to high-level porcelain commands. This project is purely educational and not intended for production use—please use the official [Git](https://git-scm.com/) (or [Go-Git](https://github.com/go-git/go-git)) for real version control needs.

## Implemented Commands

The following commands are implemented:

- [`init [path]`](./cmd/init.go): Initializes a new Git repository at the specified path (defaults to current directory).
- [`hash-object [-w] <file>`](./cmd/hash_object.go): Computes a file's SHA-1 hash, with an option to write the blob to the object database.
- [`cat-file <object-hash>`](./cmd/cat_file.go): Displays the contents of a repository object (currently supports blobs, trees, and commits).
- [`update-index [-add | -remove] <filename>`](./cmd/update_index.go): Adds or removes a file from the index.
- [`ls-files [-s]`](./cmd/ls_files.go): Lists files in the index, with an option to show detailed stage information (mode bits, hash, and path).
- [`add <filename> | .`](./cmd/add.go): Stages a single file or all files in the working directory to the index.
- [`write-tree`](./cmd/write_tree.go): Creates a tree object from the current index and writes it to the object database.
- [`ls-tree <tree-hash>`](./cmd/ls_tree.go): Print the tree contents (supports trees and commits)
- [`commit -m <message>`](./cmd/commit.go): Record changes to the repository
- [`log`](./cmd/log.go): Show commit logs

## Setup

To explore this project locally:

1. Ensure [Go](https://go.dev/) is installed (version 1.18 or later recommended).
2. Clone the repository:
   ```bash
   git clone https://github.com/CLBRITTON2/go-git-good.git
   ```
3. Navigate to the project directory:
   ```bash
   cd go-git-good
   ```
4. Build the binary:
   ```bash
   go build
   ```
5. Rename the binary to `gitgood` (or something shorter honestly) for convenience if you just want to mess with commands:
   ```bash
   mv go-git-good gitgood
   ```

## Usage

After building, run commands using the executable (e.g., `./gitgood`) followed by a command and optional flags. Examples:

- Initialize a repository:
  ```bash
  ./gitgood init
  ```
- Hash a file:
  ```bash
  ./gitgood hash-object example.txt
  ```
- Stage a file:
  ```bash
  ./gitgood add example.txt
  ```

For detailed usage, run `./gitgood` without arguments to see the help menu, or refer to the linked source files.

## Project Structure

- `cmd/`: Contains command implementations (linked above).
- `common/`: Shared utilities for repository ([`repository.go`](./common/repository.go)) and index management ([`index.go`](./common/index.go)).
- `objects/`: Logic for creating Git objects, such as blobs.
## Notes

- This project is a learning exercise and not a replacement for Git. It lacks many features and optimizations found in the official Git implementation.
- Contributions or feedback are welcome, but the focus remains on educational exploration rather than production readiness.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
