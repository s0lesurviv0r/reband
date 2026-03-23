# reband
Channel conversion CLI that converts between different radio programming CSV formats

## Installation

```sh
git clone https://github.com/s0lesurviv0r/reband
cd reband
make build
```

The binary will be placed at `build/reband`.

## Usage

```
freq-conv [--debug] <command> [flags]
```

### Global flags

| Flag | Description |
|---|---|
| `--debug` | Enable debug logging |

---

### `decode`

Decode a channel list file and print its contents.

```sh
freq-conv decode --format <format> --path <file>
```

| Flag | Required | Description |
|---|---|---|
| `--format` | yes | Source format (see supported formats below) |
| `--path` | yes | Path to the input file |

**Example:**
```sh
freq-conv decode --format bc125py --path channels.csv
```

---

### `convert`

Convert a channel list from one format to another.

```sh
freq-conv convert --from <format> --to <format> --input <file> [--output <file>]
```

| Flag | Required | Description |
|---|---|---|
| `--from` | yes | Source format |
| `--to` | yes | Destination format |
| `--input` | yes | Path to the input file |
| `--output` | no | Path to the output file (defaults to stdout) |

**Examples:**
```sh
# Convert BC125PY to Reband CSV, write to file
freq-conv convert --from bc125py --to reband --input scanner.csv --output out.csv

# Convert CHIRP to Reband CSV, print to stdout
freq-conv convert --from chirp --to reband --input chirp.csv

# Convert Reband CSV to BC125PY
freq-conv convert --from reband --to bc125py --input channels.csv --output scanner.csv
```

---

## Supported formats

| Format key | Description | Decode | Encode |
|---|---|:---:|:---:|
| `reband` | Reband CSV (all channel fields) | ✓ | ✓ |
| `bc125py` | Uniden BC125AT/BC125XLT scanner | ✓ | ✓ |
| `chirp` | CHIRP radio programming software | ✓ | |
| `uv-pro` | UV-Pro | | |
| `radioreference` | RadioReference | | |
| `sdrtrunk` | SDRTrunk | | |
| `gqrx` | GQRX | | |

### Reband CSV format

The `reband` format is a lossless interchange format that represents all channel fields. It can be used as an intermediate when converting between two formats that do not share a direct converter.

| Column | Description |
|---|---|
| `Index` | Channel number |
| `Name` | Channel name |
| `AlphaTag` | Short display tag |
| `Comment` | Free-form comment |
| `Frequency` | Frequency in MHz |
| `Duplex` | Duplex direction: empty, `+`, or `-` |
| `Offset` | Repeater offset in MHz |
| `ToneType` | Squelch tone type: `none`, `ctcss`, or `dcs` |
| `ToneValue` | Tone value (CTCSS: 1/10 Hz units; DCS: code number) |
| `Modulation` | Modulation mode (see below) |
| `Power` | Transmit power in watts |
| `Delay` | Squelch delay in seconds |
| `Lockout` | `true` or `false` |
| `Priority` | `true` or `false` |

Supported modulation values: `fm`, `nfm`, `am`, `wfm`, `lsb`, `usb`, `cw`, `c4fm`, `dstar`, `p25`, `nxdn`, `dmr`, `ysf`, `fusion`, `pocsag`, `dpmr`, `tetra`

---

## Notes

- UV-PRO should allow an option to split into multiple CSVs, one per channel group
- Generally the CLI should allow splitting into smaller CSVs per request
