# morph
## Features
- Offline only, privacy-focused CLI
- Flags-based: `morph --in file --to format [--out file] [--from format]`
- Many formats: images, text/data (json/yaml/toml/csv/ndjson), html/md, pdf, office (docx/xlsx/pptx/odt/ods/odp), epub, rtf, svg, archives (zip/tar/gz/tar.xz/tar.bz2), etc.
- Registry-based design for easy extension

## Install
```bash
# Build
go build -o morph ./cmd/morph

# Install
go install ./cmd/morph
```

## Usage
```bash
morph --in file.pdf --to docx
morph --in image.png --to jpg
morph doctor
morph --list
morph --list-json
morph completion bash > /etc/bash_completion.d/morph
```

## Notes
- Offline only, privacy-focused
- Flags-based CLI
- Supports many formats across documents, images, data, archives, office, ebooks
