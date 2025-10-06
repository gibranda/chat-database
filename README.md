# 🤖 AI Database Agent - Asisten Database Cerdas

Sistem AI canggih yang mengubah bahasa natural menjadi query SQL menggunakan Large Language Model (LLM) lokal dengan kemampuan eksplorasi database yang intelligent.

## ✨ Fitur Utama

### 🎯 Kemampuan Inti
- **Koneksi Database Dinamis**: Connect ke database apapun langsung dari UI tanpa konfigurasi
- **Natural Language to SQL**: Ubah pertanyaan bahasa Indonesia/Inggris menjadi SQL query yang akurat
- **AI Reasoning**: Proses berpikir AI yang transparan dengan multi-step reasoning
- **Insight Otomatis**: AI memberikan analisis dan insight dalam Bahasa Indonesia
- **Schema Explorer**: Eksplorasi struktur database secara visual
- **Query History**: Riwayat query dengan kemampuan re-run

### 🗄️ Database Support
- **PostgreSQL** 🐘 - Full support dengan semua fitur
- **MySQL** 🐬 - Compatible dengan MySQL 5.7+
- **SQLite** 📁 - Perfect untuk development dan testing

### 🎨 User Interface
- **Modern & Responsive**: Desain clean dengan Nuxt.js + Vue 3 + TailwindCSS
- **Real-time Updates**: Live feedback saat AI memproses query
- **Collapsible Sections**: SQL dan reasoning bisa di-show/hide
- **Beautiful Tables**: Data hasil query dengan formatting otomatis
- **Dark Mode Ready**: UI yang nyaman untuk mata

### 🔒 Keamanan & Privacy
- **Local LLM**: Semua AI processing di lokal menggunakan Ollama
- **No External API**: Data tidak keluar dari infrastruktur Anda
- **Readonly Mode**: Proteksi dari query yang merusak data
- **SQL Validation**: Validasi otomatis untuk mencegah SQL injection

## 🏗️ Arsitektur Sistem

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     FRONTEND (Nuxt.js)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Connection   │  │ Chat         │  │ Schema       │      │
│  │ Modal        │  │ Interface    │  │ Explorer     │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                  │                  │              │
│         └──────────────────┼──────────────────┘              │
│                            │ HTTP/REST API                   │
└────────────────────────────┼─────────────────────────────────┘
                             │
┌────────────────────────────┼─────────────────────────────────┐
│                     BACKEND (Golang)                         │
│  ┌──────────────────────────┴───────────────────────────┐   │
│  │              API Layer (Gin Framework)               │   │
│  │  • Connection Management  • Query Processing         │   │
│  │  • Schema API            • Health Check              │   │
│  └──────────────────────────┬───────────────────────────┘   │
│                             │                                │
│  ┌──────────────────────────┴───────────────────────────┐   │
│  │                  AI Agent Engine                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐           │   │
│  │  │ Planning │→ │ SQL Gen  │→ │ Answer   │           │   │
│  │  │ Step     │  │ Step     │  │ Step     │           │   │
│  │  └──────────┘  └──────────┘  └──────────┘           │   │
│  └──────────────────────────┬───────────────────────────┘   │
│                             │                                │
│         ┌───────────────────┼───────────────────┐           │
│         │                   │                   │           │
│  ┌──────▼──────┐    ┌──────▼──────┐    ┌──────▼──────┐    │
│  │   Ollama    │    │  Database   │    │   Config    │    │
│  │   Client    │    │   Layer     │    │  Manager    │    │
│  └─────────────┘    └─────────────┘    └─────────────┘    │
└────────────────────────────────────────────────────────────┘
           │                    │
           │                    │
┌──────────▼────────┐  ┌────────▼──────────┐
│  Ollama Server    │  │   Your Database   │
│  (llama3.2, etc)  │  │  (PG/MySQL/SQLite)│
└───────────────────┘  └───────────────────┘
```

### Component Details

#### Frontend Components
```
frontend/
├── pages/
│   └── index.vue              # Main chat interface
├── components/
│   ├── ConnectionModal.vue    # Database connection UI
│   ├── ChatMessage.vue        # Message display with SQL & results
│   ├── QueryResults.vue       # Beautiful table with formatting
│   ├── ReasoningSteps.vue     # AI reasoning visualization
│   ├── SchemaExplorer.vue     # Database schema browser
│   └── QueryHistory.vue       # Query history sidebar
├── composables/
│   └── useApi.ts              # API client functions
└── stores/
    ├── chat.ts                # Chat state management
    └── schema.ts              # Schema state management
```

#### Backend Components
```
backend/
├── cmd/server/
│   └── main.go                # Application entry point
├── internal/
│   ├── agent/
│   │   └── agent.go           # AI agent with reasoning engine
│   ├── api/
│   │   ├── handlers.go        # HTTP request handlers
│   │   ├── router.go          # Route definitions
│   │   └── connection.go      # Dynamic connection handlers
│   ├── database/
│   │   └── database.go        # Database abstraction layer
│   ├── llm/
│   │   └── ollama.go          # Ollama LLM client
│   └── config/
│       └── config.go          # Configuration management
```

## 🚀 Instalasi & Setup

### Prasyarat

Pastikan sistem Anda memiliki:

- **Go** 1.21 atau lebih baru
- **Node.js** 18+ atau **Bun**
- **Ollama** terinstall dan running
- **Database** (PostgreSQL/MySQL/SQLite) - opsional, bisa connect dari UI

### Langkah 1: Clone Repository

```bash
git clone https://github.com/yourusername/chat-with-database.git
cd chat-with-database
```

### Langkah 2: Setup Otomatis

Gunakan script setup untuk instalasi cepat:

```bash
chmod +x setup.sh
./setup.sh
```

Script ini akan:
- ✅ Check semua prerequisites
- ✅ Install Go dependencies
- ✅ Install Node.js dependencies
- ✅ Create config files
- ✅ Setup environment variables

### Langkah 3: Setup Ollama

Install dan jalankan Ollama:

```bash
# Install Ollama (macOS)
brew install ollama

# Atau download dari https://ollama.ai

# Start Ollama server
ollama serve

# Pull model (di terminal baru)
ollama pull llama3.2
# atau
ollama pull llama3.1
```

**Model yang Direkomendasikan:**
- `llama3.2` - Cepat, cocok untuk development (3B/1B)
- `llama3.1` - Lebih akurat, untuk production (8B)
- `codellama` - Bagus untuk SQL generation
- `mistral` - Alternative yang cepat

### Langkah 4: Jalankan Aplikasi

#### Mode Production

**Terminal 1 - Backend:**
```bash
./start-backend.sh
```

**Terminal 2 - Frontend:**
```bash
./start-frontend.sh
```

#### Mode Development (dengan Hot Reloading) 🔥

**Terminal 1 - Backend dengan Hot Reload:**
```bash
./dev-backend.sh
```

Script ini akan:
- ✅ Auto-install `air` jika belum ada
- ✅ Watch perubahan di file `.go` dan `.yaml`
- ✅ Auto-reload server saat ada perubahan
- ✅ Menampilkan build errors secara real-time

**Terminal 2 - Frontend:**
```bash
./start-frontend.sh  # Nuxt sudah support hot reload by default
```

**Atau manual:**

```bash
# Backend (tanpa hot reload)
cd backend
go run cmd/server/main.go

# Backend (dengan hot reload)
cd backend
air  # pastikan air sudah terinstall

# Frontend (terminal baru)
cd frontend
npm run dev
```

### Langkah 5: Akses Aplikasi

Buka browser dan akses:
```
http://localhost:3000
```

Modal koneksi database akan muncul otomatis!

## 📖 Cara Penggunaan

### 1. Koneksi ke Database

Saat pertama kali membuka aplikasi, modal koneksi akan muncul:

#### PostgreSQL
```
Database Type: PostgreSQL 🐘
Host: localhost
Port: 5432
Database: nama_database
Username: postgres
Password: password_anda
SSL Mode: disable
```

#### MySQL
```
Database Type: MySQL 🐬
Host: localhost
Port: 3306
Database: nama_database
Username: root
Password: password_anda
```

#### SQLite
```
Database Type: SQLite 📁
Path: ./database.db
```

**Langkah Koneksi:**
1. Pilih tipe database
2. Isi kredensial
3. Klik **"Test Connection"** untuk validasi
4. Jika sukses, klik **"Connect & Start"**
5. Mulai chat!

### 2. Bertanya ke Database

Setelah terkoneksi, Anda bisa bertanya dalam bahasa natural:

**Contoh Pertanyaan:**

```
📊 Eksplorasi Data:
- "Tampilkan semua tabel yang ada"
- "Apa saja kolom di tabel users?"
- "Berapa jumlah data di tabel products?"

📈 Analisis:
- "Tampilkan 10 produk terlaris"
- "Berapa total penjualan bulan ini?"
- "Siapa customer dengan pembelian terbanyak?"

🔍 Query Kompleks:
- "Bandingkan penjualan Q1 vs Q2 tahun ini"
- "Tampilkan produk yang belum pernah terjual"
- "Cari customer yang tidak aktif 6 bulan terakhir"

📊 Agregasi:
- "Hitung rata-rata harga per kategori"
- "Grouping penjualan per bulan"
- "Top 5 sales person tahun ini"
```

### 3. Memahami Response

Setiap response AI terdiri dari:

#### ✅ Insight (Bahasa Indonesia)
```
Berdasarkan data yang saya temukan, terdapat 150 produk 
dalam database. Yang menarik adalah kategori Electronics 
memiliki produk terbanyak (45 items). Secara angka, rata-rata 
harga produk adalah Rp 250.000.
```

#### 📊 Query Results
Tabel data dengan formatting otomatis:
- Numbers dengan thousand separator
- NULL values dengan styling khusus
- Boolean dengan color coding
- Sticky header untuk scroll

#### 💻 Generated SQL (Collapsible)
```sql
SELECT category, COUNT(*) as total, AVG(price) as avg_price
FROM products
GROUP BY category
ORDER BY total DESC
LIMIT 10;
```

#### 🧠 Agent Reasoning (Collapsible)
Proses berpikir AI step-by-step:
1. Analyze Schema
2. Plan Query
3. Generate SQL
4. Execute & Validate
5. Generate Answer

### 4. Fitur Tambahan

#### Schema Explorer
- Lihat semua tabel di database
- Expand untuk melihat kolom
- Klik tabel untuk info detail

#### Query History
- Semua query tersimpan otomatis
- Klik untuk re-run query
- Clear history kapan saja

#### Collapsible Sections
- SQL dan Reasoning default collapsed
- Klik untuk expand/collapse
- UI lebih bersih dan fokus

## ⚙️ Konfigurasi

### Backend Configuration

Edit `backend/config.yaml`:

```yaml
# Ollama LLM Configuration
ollama:
  host: http://localhost:11434
  model: llama3.2              # Model yang digunakan
  temperature: 0.1             # Kreativitas (0.0-1.0)
  timeout: 120                 # Timeout dalam detik

# Agent Configuration
agent:
  max_iterations: 5            # Max reasoning steps
  enable_query_validation: true
  readonly_mode: true          # Proteksi dari UPDATE/DELETE
  max_results: 100             # Limit hasil query

# Server Configuration
server:
  host: 0.0.0.0
  port: 8080
  debug: true
```

### Frontend Configuration

Edit `frontend/.env`:

```bash
# API Backend URL
NUXT_PUBLIC_API_URL=http://localhost:8080

# App Configuration
NUXT_PUBLIC_APP_NAME=AI Database Agent
```

### Environment Variables

```bash
# Backend
export CONFIG_PATH=./config.yaml
export OLLAMA_HOST=http://localhost:11434

# Frontend
export NUXT_PUBLIC_API_URL=http://localhost:8080
```

## 🔧 Development

### Hot Reloading Setup 🔥

Untuk development yang lebih produktif, gunakan hot reloading:

#### Backend Hot Reload (Air)

**Install Air:**
```bash
go install github.com/air-verse/air@latest
```

**Jalankan dengan hot reload:**
```bash
./dev-backend.sh
# atau
cd backend && air
```

**Konfigurasi Air** (sudah tersedia di `backend/.air.toml`):
- Watch file: `*.go`, `*.yaml`, `*.yml`
- Exclude: `*_test.go`, `tmp/`, `vendor/`
- Auto rebuild on file changes
- Build errors ditampilkan real-time

#### Frontend Hot Reload (Nuxt)

Frontend sudah support hot reload by default:
```bash
./start-frontend.sh
# atau
cd frontend && npm run dev
```

**Fitur Hot Reload:**
- ✅ Component auto-reload
- ✅ CSS hot update
- ✅ State preservation
- ✅ Fast refresh

### Project Structure

```
chat-with-database/
├── backend/                 # Golang backend
│   ├── cmd/server/         # Main application
│   ├── internal/           # Internal packages
│   │   ├── agent/         # AI agent logic
│   │   ├── api/           # HTTP handlers
│   │   ├── database/      # DB abstraction
│   │   ├── llm/           # LLM client
│   │   └── config/        # Configuration
│   ├── .air.toml          # Air hot reload config
│   ├── config.yaml        # Backend config
│   └── go.mod             # Go dependencies
│
├── frontend/               # Nuxt.js frontend
│   ├── pages/             # Route pages
│   ├── components/        # Vue components
│   ├── composables/       # Composable functions
│   ├── stores/            # Pinia stores
│   ├── assets/            # CSS & static files
│   └── nuxt.config.ts     # Nuxt configuration
│
├── setup.sh               # Auto setup script
├── start-backend.sh       # Backend starter (production)
├── dev-backend.sh         # Backend starter (development + hot reload)
├── start-frontend.sh      # Frontend starter
└── README.md              # This file
```

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm run test
```

### Building for Production

```bash
# Backend
cd backend
go build -o bin/server cmd/server/main.go

# Frontend
cd frontend
npm run build
```

### Docker Deployment

```bash
# Build images
docker-compose build

# Run services
docker-compose up -d

# View logs
docker-compose logs -f
```

## 🎯 Use Cases

### 1. Business Intelligence
```
"Berapa revenue bulan ini dibanding bulan lalu?"
"Tampilkan trend penjualan 6 bulan terakhir"
"Produk apa yang paling menguntungkan?"
```

### 2. Data Exploration
```
"Apa saja tabel yang ada?"
"Tampilkan struktur tabel customers"
"Berapa total record di semua tabel?"
```

### 3. Customer Analytics
```
"Siapa customer dengan lifetime value tertinggi?"
"Berapa customer baru bulan ini?"
"Customer mana yang churn?"
```

### 4. Inventory Management
```
"Produk apa yang stoknya menipis?"
"Tampilkan produk yang belum pernah terjual"
"Berapa nilai total inventory?"
```

### 5. Sales Analysis
```
"Top 10 sales person bulan ini"
"Bandingkan penjualan per region"
"Apa produk yang paling sering dibeli bersamaan?"
```

## �� Keamanan

### Best Practices

1. **Readonly Mode**: Aktifkan untuk mencegah perubahan data
2. **SQL Validation**: Validasi otomatis untuk SQL injection
3. **Connection Encryption**: Gunakan SSL untuk database connection
4. **Environment Variables**: Jangan hardcode credentials
5. **Access Control**: Implement user authentication (coming soon)

### Readonly Mode

Dalam mode readonly, hanya query SELECT yang diizinkan:

```yaml
agent:
  readonly_mode: true  # Hanya SELECT, SHOW, DESCRIBE, EXPLAIN
```

Query yang diblok:
- INSERT, UPDATE, DELETE
- DROP, TRUNCATE, ALTER
- CREATE, GRANT, REVOKE

## 🚨 Troubleshooting

### Backend Issues

**Error: "Failed to connect to Ollama"**
```bash
# Pastikan Ollama running
ollama serve

# Check status
curl http://localhost:11434/api/tags
```

**Error: "Model not found"**
```bash
# Pull model yang dibutuhkan
ollama pull llama3.2
```

**Error: "Database connection failed"**
- Check kredensial database
- Pastikan database service running
- Verify network connectivity

### Frontend Issues

**Error: "Failed to connect to backend"**
- Check backend running di port 8080
- Verify NUXT_PUBLIC_API_URL di .env
- Check CORS configuration

**Modal tidak muncul**
- Clear browser cache
- Check console untuk errors
- Verify isConnected state

### Performance Issues

**Query terlalu lambat**
- Gunakan model Ollama yang lebih kecil (llama3.2-1b)
- Reduce temperature untuk response lebih cepat
- Add database indexes
- Limit result dengan LIMIT clause

**Memory usage tinggi**
- Reduce max_results di config
- Clear query history
- Restart Ollama service

## �� Performance

### Benchmarks

| Operation | Time | Notes |
|-----------|------|-------|
| Connection Test | < 1s | Database ping |
| Schema Loading | 1-3s | Depends on DB size |
| Query Planning | 2-5s | LLM inference |
| SQL Generation | 3-8s | LLM inference |
| Query Execution | < 1s | Depends on query |
| Answer Generation | 2-5s | LLM inference |
| **Total** | **8-20s** | End-to-end |

### Optimization Tips

1. **Use Faster Models**: llama3.2-1b vs llama3.1-8b
2. **Enable Caching**: Schema caching enabled by default
3. **Limit Results**: Set max_results to reasonable number
4. **Database Indexes**: Add indexes for frequently queried columns
5. **Connection Pooling**: Reuse database connections

## 🗺️ Roadmap

### Version 1.1 (Current)
- ✅ Dynamic database connection
- ✅ Beautiful UI with Nuxt.js
- ✅ Indonesian language support
- ✅ Collapsible SQL & reasoning
- ✅ Enhanced data formatting

### Version 1.2 (Next)
- [ ] User authentication & authorization
- [ ] Multi-user support
- [ ] Query result export (Excel, PDF)
- [ ] Data visualization charts
- [ ] Query performance metrics

### Version 2.0 (Future)
- [ ] Multiple database connections
- [ ] Scheduled queries & reports
- [ ] Email notifications
- [ ] API key management
- [ ] Advanced analytics dashboard

## 🤝 Contributing

Kontribusi sangat diterima! Silakan:

1. Fork repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Ollama** - Local LLM inference
- **Nuxt.js** - Modern Vue framework
- **Gin** - Fast Go web framework
- **TailwindCSS** - Utility-first CSS
- **Pinia** - Vue state management

## 📧 Contact & Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/chat-with-database/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/chat-with-database/discussions)
- **Email**: your.email@example.com

## 🌟 Star History

Jika project ini membantu Anda, berikan ⭐ di GitHub!

---

**Made with ❤️ by [Your Name]**

*Transforming natural language into database insights, one query at a time.*
