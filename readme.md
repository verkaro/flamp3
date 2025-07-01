# flamp3

**Convert FLAC files to MP3** in batch, with adjustable LAME quality  
A companion Go utility using a shell pipeline (`flac -c -d | lame`).

---

## Features

- **Selective**: only processes `.flac` files  
- **Recursive** mode (`-recursive`) preserves directory structure  
- **Skip existing** MP3s to avoid re-encoding  
- **Adjustable LAME VBR quality** (`-quality`, 0=best → 9=worst)  

---

## Requirements

- [flac](https://xiph.org/flac/)  
- [lame](http://lame.sourceforge.net/)  
- `sh` (for the pipeline)

---

## Installation

```bash
git clone https://github.com/verkaro/flamp3.git
cd flamp3
go build -o flamp3 main.go
````

---

## Usage

```bash
# Single FLAC to MP3 at default quality (6)
./flamp3 -out mp3s album.flac

# Convert entire folder to high quality MP3s
./flamp3 -out mp3s -recursive -quality 2 ~/music_library
```

**Flags:**

* `-out` — Output root directory (creates it if needed)
* `-recursive` — Recurse into input directories
* `-quality` — LAME VBR level (0=best, 9=worst; default `6`)

---

## Contributing

Please open issues for feature requests or bugs. PRs should include tests and follow Go style.

---

## License

MIT License

---

## Credits

* Utility designed and iteratively refined with support from [ChatGPT](https://openai.com).
* Audio conversion powered by `flac` and `lame`.

