# Changelog

## [Unreleased] - 2025-10-02

### ✨ Improved UI/UX

#### Chat Message Display
- **Insight Section**: Menampilkan ringkasan dan insight di bagian atas dengan icon checkmark
- **Collapsible SQL**: Generated SQL sekarang bisa di-show/hide dengan tombol toggle
- **Collapsible Reasoning**: Agent reasoning process default collapsed untuk UI yang lebih bersih
- **Better Error Display**: Error messages dengan styling yang lebih jelas (red border-left)

#### Query Results Table
- **Enhanced Header**: Header dengan gradient background dan summary statistics
- **Data Summary**: Menampilkan total rows dan columns count
- **Better Formatting**: 
  - Numbers dengan thousand separator (format Indonesia)
  - NULL values dengan styling italic gray
  - Boolean values dengan color coding (green/red)
  - Alternating row colors untuk readability
- **Sticky Header**: Table header tetap visible saat scroll
- **Footer Info**: Menampilkan "Showing X of Y rows" jika data > 10 rows

#### Reasoning Steps
- **Compact Design**: Button-style toggle dengan step count badge
- **Line Clamp**: Observation text di-truncate untuk menghemat space
- **Better Icons**: Menambahkan lightbulb icon untuk reasoning

### 🧠 Backend Improvements

#### Agent Prompts
- **Enhanced Answer Generation**: Prompt yang lebih baik untuk menghasilkan insight
- **Indonesian Language**: Response dalam Bahasa Indonesia
- **Structured Insights**: Meminta AI untuk memberikan:
  1. Summary (2-3 kalimat)
  2. Key insights/patterns
  3. Important numbers/statistics
  4. Interesting observations

#### Logging
- **Detailed Query Logging**: Menambahkan log untuk debugging query processing
- **Error Tracking**: Better error messages untuk troubleshooting

### 🎨 Visual Improvements

- **Gradient Accents**: Primary-purple gradient untuk visual appeal
- **Hover Effects**: Smooth transitions pada interactive elements
- **Better Spacing**: Improved padding dan margins
- **Color Coding**: Semantic colors untuk different data types

### 🐛 Bug Fixes

- Fixed Nuxt UI module conflict (removed @nuxt/ui dependency)
- Fixed Tailwind configuration issues
- Improved error handling in API handlers

### 📝 Documentation

- Added detailed troubleshooting for Ollama model issues
- Updated README with Nuxt.js frontend information
- Added QUICKSTART.md in Indonesian

## Known Issues

- Frontend masih ada TypeScript warnings (tidak mempengaruhi functionality)
- Ollama model `llama3.1` perlu di-pull terlebih dahulu
- Readonly mode error perlu di-handle dengan lebih baik di UI

## Next Steps

- [ ] Add pagination untuk large result sets
- [ ] Add data visualization charts
- [ ] Implement query caching
- [ ] Add user authentication
- [ ] Support for multiple database connections
- [ ] Export results to Excel/PDF
- [ ] Query performance metrics
- [ ] Advanced filtering and sorting in results table
