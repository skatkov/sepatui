# SEPA TUI - SEPA XML Viewer

A Terminal User Interface (TUI) application for viewing and analyzing SEPA XML payment files. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## 🚀 Features

- **Interactive Table View**: Browse SEPA XML data in an organized, categorized table format
- **Copy to Clipboard**: Easily copy any field value with a single keypress
- **Keyboard Navigation**: Intuitive vim-style navigation (j/k or arrow keys)
- **Clean Interface**: Minimalist TUI design focused on readability

## 📦 Installation

### From Source

```bash
git clone https://github.com/skatkov/sepatui.git
cd sepatui
go build -o sepa
```

## From homebrew
`brew install skatkov/tap/sepa`

### Usage

```bash
./sepa <filepath>
```

**Example:**
```bash
./sepa example/SEPA_Example_2024.xml
```

## 🎯 Supported SEPA Formats

Currently supports:
- **pain.001.001.03** - Customer Credit Transfer Initiation

The application parses and displays the following information categories:

| Category | Fields |
|----------|--------|
| **Group Header** | Message ID, Creation Date Time, Number of Transactions, Control Sum, Initiating Party |
| **Payment Info** | Payment Info ID, Payment Method, Batch Booking, Service Level, Category Purpose, Execution Date, Charge Bearer |
| **Debtor** | Name, IBAN, Currency, BIC |
| **Transaction** | End-to-End ID, Amount |
| **Creditor** | Name, IBAN, BIC |
| **Remittance** | Reference Type, Issuer, Reference |

## 🎮 Controls

| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `c` | Copy selected field value to clipboard |
| `q/esc` | Quit application |

## 📋 Example Output

When you run the application, you'll see a table like this:

```
SEPA Payment Information

┌──────────────────────┬─────────────────────────┬──────────────────────────────────────────────────┐
│ Category             │ Field                   │ Value                                            │
├──────────────────────┼─────────────────────────┼──────────────────────────────────────────────────┤
│ Group Header         │ Message ID              │ 123456789012345678                               │
│ Group Header         │ Creation Date Time      │ 2024-03-15 14:30:00                             │
│ Group Header         │ Number of Transactions  │ 1                                                │
│ Group Header         │ Control Sum             │ 250.75                                           │
│ Group Header         │ Initiating Party        │ ACME CORP B.V.                                   │
│ Payment Info         │ Payment Info ID         │ 20240315143000-7891234                           │
│ Debtor               │ Name                    │ ACME CORP B.V.                                   │
│ Debtor               │ IBAN                    │ NL91ABNA0417164300                               │
│ Creditor             │ Name                    │ WIDGET SUPPLIES LTD                              │
│ Transaction          │ Amount                  │ 250.75 EUR                                       │
└──────────────────────┴─────────────────────────┴──────────────────────────────────────────────────┘
```

## 📄 SEPA XML Format

SEPA XML files are standardized payment instruction files defined by the European Payments Council (EPC). They enable automated euro-denominated payments within the Single Euro Payments Area.

**Key characteristics:**
- **Standardized format** across all SEPA countries
- **UTF-8 encoding** required
- **EUR currency** for core SEPA schemes
- **ISO 20022 standard** compliance

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Resources

- [European Payments Council (EPC)](https://www.europeanpaymentscouncil.eu/)
- [ISO 20022 Official Standards](https://www.iso20022.org/)
- [SEPA Implementation Guidelines](https://www.europeanpaymentscouncil.eu/what-we-do/sepa-payment-schemes)
