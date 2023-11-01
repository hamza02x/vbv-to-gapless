### Quick implementation for (Quranic_Audio) verse by verse to gapless

# Requirements

- ffmepg

# Install

```bash
go install github.com/hamza72x/vbv-to-gapless@latest
```

# Usage

```
Usage of ./vbv-to-gapless:
  -n string
    	database name (required) (ex: husary)
  -o string
    	output directory path (required)
  -t int
    	number of threads (default 10)
  -vd string
    	verse by verse audio directory (required)
```

# Example

```bash
vbv-to-gapless -t 10 -n hudhaify -vd ~/Downloads/dbs/hudhaify/ -o ~/Downloads/dbs/hudhaify/gapless

# -vd directory must contain all verse by verse audio(s)
# from 001001.mp3 to 114006.mp3
```

# Output File/Folder structure


- `o/$name/<001...114.mp3>` # gapless audio files
- `o/build/<001...114.txt>` # temporary build files
- `o/$name.db` # database timing file
