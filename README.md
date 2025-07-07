\# 🌌 StarWatch



\*\*StarWatch\*\* is a lightweight, terminal-based Starlink monitoring tool with attitude. Built in Go and powered by \[`go-starlink`](https://github.com/b0ch3nski/go-starlink), it gives you real-time insight into your Dishy's performance, alignment, and signal quality — all from a clean terminal UI.



> Shellz Code Solutions by w0rmer  

> Version: `v1.0.0 alpha`  

> Repo: \[github.com/6869677a/starwatch](https://github.com/6869677a/starwatch)



---



\## 🎯 Features



\- 📡 Live Starlink dish status \& signal monitoring

\- 🛰️ Uplink / downlink throughput, POP latency, GPS status

\- 🧭 Real-time tilt \& rotation feedback with alignment scoring

\- 🧾 Event log panel with scrollback

\- 📊 CSV logging toggle (`-log`) for data hoarders

\- 📚 Built-in help modal (press `H`)

\- 🧠 Uses \[tview](https://github.com/rivo/tview) for clean, keyboard-friendly UI



---



\## 🛠️ Installation



\### 🔧 Prerequisites



\- \[Go 1.21+](https://go.dev/dl/)

\- A working Starlink Dishy (Standard/High Performance)

\- Access to the Dishy’s default IP (`192.168.100.1`)



\### 📥 Build From Source



```bash

git clone https://github.com/6869677a/starwatch.git

cd starwatch

go build -o starwatch.exe main.go

```



---



\## 🚀 Usage



```bash

./starwatch.exe \[-log] \[-version]

```



\### Flags:

\- `-log` – Log live data to `starwatch-log.csv`

\- `-version` – Print version and exit



---



\## ⌨️ Controls



| Key | Action                 |

|-----|------------------------|

| `Q` | Quit                   |

| `L` | Toggle CSV Logging     |

| `H` | Show Help / Legend     |



---



\## 📓 CSV Log Format



If logging is enabled, `starwatch-log.csv` will contain:



```csv

Timestamp,Uptime(s),POP\_Latency(ms),Downlink(kbps),Uplink(kbps),Obstruct(%),Tilt,Rotation

```



---



\## 📷 UI Preview



```

🌌 StarWatch v1.0.0 alpha by w0rmer

Dish State: ONLINE   Uptime: 12.4 hrs   Boot Count: 4   Signal Quality: ✅ Above Noise Floor



POP Latency : 43.2 ms

Downlink    : 20540.7 kbps

Uplink      : 1390.5 kbps

Obstruct %  : 1.23%



Alignment   : OKAY

Tilt        : 45.20° → 46.00° (0.80°↑)

Rotation    : 180.10° → 180.00° (0.10°↺)

Uncertainty : 1.24°



Logging: true | Time: 13:45:27 MST 07/07

```



---



\## ⚠️ Warnings \& Gotchas



\- If your dish is unreachable, the tool will show `\[red]DISH UNREACHABLE`.

\- Obstruction events >3% are logged and flagged.

\- Use responsibly — don't use this to spam the Starlink firmware endpoints.



---



\## 💡 Credits



\- Built by \[@w0rmer](https://twitter.com/0x686967)

\- Powered by \[`go-starlink`](https://github.com/b0ch3nski/go-starlink)

\- UI: \[`tview`](https://github.com/rivo/tview) + \[`tcell`](https://github.com/gdamore/tcell)



---



\## 📜 License



MIT – Fork it, mod it, embed it in a bunker UI.  

Just don’t pretend you wrote it unless you’re logging that lie in CSV.



