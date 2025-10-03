# AI Database Agent - Frontend

Modern web interface built with Nuxt 3, Vue 3, and TailwindCSS.

## Features

- 🎨 Modern, gradient-based UI design
- 💬 Real-time chat interface
- 📊 Interactive data tables with CSV export
- 🔍 Database schema explorer
- 📜 Query history tracking
- 🧠 Agent reasoning visualization
- ⚡ Fast and responsive
- 🎯 Type-safe with TypeScript

## Tech Stack

- **Nuxt 3**: Vue 3 framework with SSR support
- **Vue 3**: Progressive JavaScript framework
- **TailwindCSS**: Utility-first CSS framework
- **Pinia**: State management
- **TypeScript**: Type safety
- **@nuxt/ui**: UI component library

## Setup

### Prerequisites

- Node.js 18+ or Bun
- npm, yarn, pnpm, or bun

### Installation

```bash
# Install dependencies
npm install
# or
yarn install
# or
pnpm install
# or
bun install
```

### Configuration

1. Copy environment file:
```bash
cp .env.example .env
```

2. Edit `.env` if backend is not on localhost:8080:
```
NUXT_PUBLIC_API_BASE=http://your-backend-url:8080/api
```

## Development

Start the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

The app will be available at `http://localhost:3000`

## Production

Build for production:

```bash
npm run build
# or
yarn build
# or
pnpm build
# or
bun run build
```

Preview production build:

```bash
npm run preview
# or
yarn preview
# or
pnpm preview
# or
bun run preview
```

## Project Structure

```
frontend/
├── assets/
│   └── css/
│       └── main.css          # Global styles and Tailwind
├── components/
│   ├── ChatMessage.vue       # Chat message component
│   ├── QueryResults.vue      # Results table component
│   ├── ReasoningSteps.vue    # Agent reasoning display
│   ├── SchemaExplorer.vue    # Database schema browser
│   └── QueryHistory.vue      # Query history sidebar
├── composables/
│   └── useApi.ts            # API client composable
├── pages/
│   └── index.vue            # Main chat page
├── stores/
│   ├── chat.ts              # Chat state management
│   └── schema.ts            # Schema state management
├── app.vue                  # Root component
├── nuxt.config.ts          # Nuxt configuration
└── tailwind.config.js      # Tailwind configuration
```

## Components

### ChatMessage
Displays user and assistant messages with SQL, results, and reasoning.

### QueryResults
Interactive table with CSV export functionality.

### ReasoningSteps
Collapsible view of agent's reasoning process.

### SchemaExplorer
Browse database tables, columns, and relationships.

### QueryHistory
View and reuse previous queries.

## State Management

### Chat Store
- Messages history
- Query history
- Loading state
- Current query

### Schema Store
- Database schema
- Table selection
- Loading state
- Last update timestamp

## API Integration

All API calls are handled through the `useApi` composable:

```typescript
const api = useApi()

// Send query
const result = await api.sendQuery('Show me all customers')

// Get schema
const schema = await api.getSchema()

// Check health
const health = await api.checkHealth()
```

## Styling

### Color Scheme
- Primary: Purple gradient (#667eea to #764ba2)
- Background: Light gray with gradient accents
- Text: Gray scale for hierarchy

### Custom Classes
- `.gradient-text`: Gradient text effect
- `.card`: White card with shadow
- `.btn-primary`: Primary gradient button
- `.btn-secondary`: Secondary outline button
- `.input-field`: Styled input field
- `.code-block`: Code display block

## Features in Detail

### Real-time Chat
- Message streaming
- Auto-scroll to latest
- Timestamp display
- Error handling

### Query Results
- Sortable columns
- CSV export
- Responsive table
- Null value handling

### Schema Explorer
- Expandable tables
- Column details
- Primary/Foreign key badges
- Row count display
- Refresh capability

### Query History
- Recent queries list
- Success/failure indicators
- Click to reuse
- Clear history option

## Keyboard Shortcuts

- `Enter`: Send message
- `Shift + Enter`: New line in input

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## Performance

- Lazy loading components
- Optimized re-renders
- Efficient state updates
- Minimal bundle size

## Troubleshooting

### Backend not connecting
- Check if backend is running on port 8080
- Verify `NUXT_PUBLIC_API_BASE` in `.env`
- Check browser console for CORS errors

### Styles not loading
- Clear `.nuxt` cache: `rm -rf .nuxt`
- Reinstall dependencies
- Restart dev server

### TypeScript errors
- Run `npm run postinstall` to generate types
- Check `tsconfig.json` extends `.nuxt/tsconfig.json`

## Contributing

1. Follow Vue 3 Composition API style
2. Use TypeScript for type safety
3. Follow TailwindCSS utility-first approach
4. Keep components small and focused
5. Add proper error handling

## License

MIT
