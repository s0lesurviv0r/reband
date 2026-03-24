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
reband [--debug] <command> [flags]
```

### Global flags

| Flag | Description |
|---|---|
| `--debug` | Enable debug logging |

---

### `decode`

Decode a channel list file and print its contents.

```sh
reband decode --format <format> --path <file>
```

| Flag | Required | Description |
|---|---|---|
| `--format` | yes | Source format (see supported formats below) |
| `--path` | yes | Path to the input file |
| `--on-error` | no | How to handle row errors: `exit` (default), `skip`, or `empty` |

**Example:**
```sh
reband decode --format bc125py --path channels.csv
```

---

### `convert`

Convert a channel list from one format to another.

```sh
reband convert --from <format> --to <format> --input <file> [--output <file>]
reband convert --from <format> --to <format> --input <file> --split-output-size <n> --output-dir <dir>
```

| Flag | Required | Description |
|---|---|---|
| `--from` | yes | Source format |
| `--to` | yes | Destination format |
| `--input` | yes | Path to the input file |
| `--output` | no | Path to the output file (defaults to stdout) |
| `--split-output-size` | no | Split output into multiple files with at most this many channels each |
| `--output-dir` | no | Directory to write split output files (required with `--split-output-size`) |
| `--on-error` | no | How to handle row errors: `exit` (default), `skip`, or `empty` |

When `--split-output-size` is set, output files are written to `--output-dir` named `01.csv`, `02.csv`, etc., zero-padded to the number of digits needed. `--output` and `--split-output-size` are mutually exclusive.

**Examples:**
```sh
# Convert BC125PY to Reband CSV, write to file
reband convert --from bc125py --to reband --input scanner.csv --output out.csv

# Convert CHIRP to Reband CSV, print to stdout
reband convert --from chirp --to reband --input chirp.csv

# Convert Reband CSV to BC125PY
reband convert --from reband --to bc125py --input channels.csv --output scanner.csv

# Convert UV-Pro CSV to Reband CSV, split into files of 50 channels each
reband convert --from uv-pro --to reband --input channels.csv --split-output-size 50 --output-dir out/
```

---

## Supported formats

| Format key | Description | Decode | Encode |
|---|---|:---:|:---:|
| `reband` | Reband CSV (all channel fields) | ✓ | ✓ |
| `bc125py` | Uniden BC125AT/BC125XLT scanner | ✓ | ✓ |
| `chirp` | CHIRP radio programming software | ✓ | ✓ |
| `uv-pro` | UV-Pro | ✓ | ✓ |
| `radioreference` | RadioReference | ✓ | ✓ |
| `sdrtrunk` | SDRTrunk | | |

### Reband CSV format

The `reband` format is a lossless interchange format that represents all channel fields. It is the recommended intermediate when converting between two formats that do not share a direct converter.

A sample file is provided at [`samples/reband.csv`](samples/reband.csv).

#### Columns

| # | Column | Type | Description |
|---|---|---|---|
| 1 | `Index` | integer | Channel number (1-based position in the channel list) |
| 2 | `Name` | string | Channel name, typically up to 16 characters |
| 3 | `AlphaTag` | string | Short display tag shown on radio LCD, typically up to 8 characters |
| 4 | `Comment` | string | Free-form human-readable comment |
| 5 | `Frequency` | decimal (MHz) | Receive frequency in MHz, 4 decimal places (e.g. `462.5500`) |
| 6 | `Duplex` | enum | Repeater offset direction: empty for simplex, `+` for positive offset, `-` for negative offset |
| 7 | `Offset` | decimal (MHz) | Repeater transmit offset in MHz, 4 decimal places (e.g. `0.6000`); `0.0000` when unused |
| 8 | `ToneType` | enum | Squelch tone type: `none`, `ctcss`, or `dcs` |
| 9 | `ToneValue` | integer | Tone code; for `ctcss`: frequency in 1/10 Hz units (e.g. `1567` = 156.7 Hz); for `dcs`: DCS code number (e.g. `71`); `0` when `ToneType` is `none` |
| 10 | `Modulation` | enum | Modulation mode (see table below) |
| 11 | `Bandwidth` | decimal (kHz) | Channel bandwidth in kHz (e.g. `12.5` for narrowband FM, `25.0` for wideband FM); `0.0` when not applicable |
| 12 | `Power` | integer (W) | Transmit power in watts; `0` when unknown or receive-only |
| 13 | `Delay` | integer (s) | Squelch hang time in seconds after signal drops |
| 14 | `Lockout` | boolean | `true` to exclude the channel from scanning, `false` to include it |
| 15 | `Priority` | boolean | `true` to mark the channel as a priority scan channel, `false` otherwise |

#### Modulation values

| Value | Mode | Notes |
|---|---|---|
| `fm` | FM | Use `Bandwidth` to distinguish narrowband (`12.5` kHz) from wideband (`25.0` kHz) |
| `am` | AM | Used for aviation and some shortwave bands |
| `wfm` | Wide FM | Broadcast FM radio |
| `lsb` | LSB | Lower sideband SSB |
| `usb` | USB | Upper sideband SSB |
| `cw` | CW | Morse code / continuous wave |
| `c4fm` | C4FM | Yaesu digital voice (System Fusion) |
| `dstar` | D-STAR | Icom digital voice |
| `p25` | P25 | APCO Project 25 digital |
| `nxdn` | NXDN | Icom/Kenwood digital voice |
| `dmr` | DMR | Digital Mobile Radio |
| `ysf` | YSF | Yaesu System Fusion |
| `fusion` | Fusion | Yaesu Fusion (alias) |
| `pocsag` | POCSAG | Paging protocol |
| `dpmr` | dPMR | Digital Private Mobile Radio |
| `tetra` | TETRA | Terrestrial Trunked Radio |

#### Example row

```
Index,Name,AlphaTag,Comment,Frequency,Duplex,Offset,ToneType,ToneValue,Modulation,Bandwidth,Power,Delay,Lockout,Priority
5,2m Repeater,2M RPT,Local 2m repeater,146.9400,+,0.6000,ctcss,1000,fm,25.0,50,2,false,false
```

This represents a wideband FM repeater on 146.940 MHz with a +600 kHz offset, CTCSS tone of 100.0 Hz (ToneValue `1000` = 100.0 Hz), 50W output, 2-second squelch delay.

---

