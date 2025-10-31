# Crypsis Frontend - Quick Reference

## 🎯 Project Overview

Modern React + TypeScript admin dashboard for Crypsis file encryption service.

### Tech Stack
- React 18.3 + TypeScript 5.5
- Vite (build tool)
- Tailwind CSS (styling)
- Zustand (state management)
- Axios (HTTP client)
- React Router (routing)
- Lucide React (icons)

## 🚀 Quick Start

```bash
# Install dependencies
npm install

# Start development server
npm run dev
# Opens at http://localhost:3000

# Build for production
npm run build

# Preview production build
npm run preview
```

## 📁 Project Structure

```
src/
├── components/
│   ├── features/          # Feature modules
│   │   ├── files/         # File management components
│   │   ├── admin/         # Admin management
│   │   ├── applications/  # OAuth apps
│   │   ├── logs/          # Audit logs
│   │   └── security/      # Security settings
│   ├── layout/            # App layout
│   │   ├── AppLayout.tsx
│   │   ├── Header.tsx
│   │   ├── Sidebar.tsx
│   │   └── ProtectedRoute.tsx
│   └── ui/                # Reusable UI components
│       ├── Button.tsx
│       ├── Input.tsx
│       ├── Card.tsx
│       ├── Alert.tsx
│       └── LoadingSpinner.tsx
├── pages/                 # Page components
│   ├── LoginPage.tsx
│   ├── DashboardPage.tsx
│   ├── FilesPage.tsx
│   ├── AdminsPage.tsx
│   ├── ApplicationsPage.tsx
│   ├── LogsPage.tsx
│   └── SecurityPage.tsx
├── services/              # API services
│   ├── api.ts            # Axios client
│   ├── auth.ts           # Auth service
│   └── backend.ts        # Backend API calls
├── stores/                # Zustand stores
│   ├── auth.ts           # Auth state
│   ├── files.ts          # Files state
│   ├── admins.ts         # Admins state
│   ├── applications.ts   # Apps state
│   └── logs.ts           # Logs state
├── types/                 # TypeScript types
│   └── index.ts
├── constants/             # Constants
│   └── index.ts          # API endpoints, routes
├── utils/                 # Utilities
│   └── index.ts
├── App.tsx               # Main app component
├── main.tsx              # Entry point
└── index.css             # Global styles
```

## 🎨 Styling with Tailwind

### Example Component
```tsx
export function MyComponent() {
  return (
    <div className="card">
      <h1 className="text-2xl font-bold text-gray-900 mb-4">
        Title
      </h1>
      <button className="btn-primary px-4 py-2">
        Click Me
      </button>
    </div>
  );
}
```

### Custom Classes (in index.css)
- `.card` - White background with shadow
- `.btn-primary` - Primary button style
- `.btn-secondary` - Secondary button style
- `.btn-danger` - Danger/delete button
- `.input` - Standard input field
- `.label` - Form label

## 🔌 API Integration

### Example Service Call
```tsx
import { filesService } from '@/services';

// Upload file
const handleUpload = async (file: File) => {
  try {
    const fileId = await filesService.uploadFile(file, (progress) => {
      console.log(`Upload progress: ${progress}%`);
    });
    console.log('Uploaded:', fileId);
  } catch (error) {
    console.error('Upload failed:', error);
  }
};
```

### Using Stores
```tsx
import { useFileStore } from '@/stores';

function FilesPage() {
  const { files, isLoading, fetchFiles, deleteFile } = useFileStore();
  
  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);
  
  return (
    <div>
      {isLoading ? (
        <LoadingSpinner />
      ) : (
        files.map(file => (
          <div key={file.id}>
            {file.file_name}
            <button onClick={() => deleteFile(file.id)}>Delete</button>
          </div>
        ))
      )}
    </div>
  );
}
```

## 🔐 Authentication

### Protected Routes
```tsx
import { ProtectedRoute } from '@/components/layout';

<Route
  path="/admin"
  element={
    <ProtectedRoute>
      <AdminPage />
    </ProtectedRoute>
  }
/>
```

### Using Auth Store
```tsx
import { useAuthStore } from '@/stores';

function Header() {
  const { user, logout } = useAuthStore();
  
  return (
    <div>
      <span>Welcome, {user?.username}</span>
      <button onClick={logout}>Logout</button>
    </div>
  );
}
```

## 🛠️ Development Tips

### Hot Reload
Vite provides instant hot module replacement. Changes appear immediately.

### Type Safety
Use TypeScript interfaces for all data:
```tsx
interface FileItem {
  id: string;
  file_name: string;
  file_size: number;
  // ...
}
```

### Error Handling
```tsx
try {
  await someApiCall();
} catch (error) {
  if (error instanceof ApiError) {
    alert(error.message);
  }
}
```

## 📦 Building for Production

```bash
# Build
npm run build

# Output in ./dist directory
# Optimized, minified, tree-shaken

# Preview
npm run preview
```

## 🐳 Docker Deployment

```bash
# Build image
docker build -t crypsis-frontend:latest .

# Run
docker run -p 3000:80 \
  -e VITE_API_URL=http://api.example.com \
  crypsis-frontend:latest
```

## 🔧 Environment Variables

Create `.env`:
```bash
VITE_API_URL=http://localhost:8080
VITE_ENV=development
```

Access in code:
```tsx
const apiUrl = import.meta.env.VITE_API_URL;
```

## 🎯 Common Tasks

### Add a New Page
1. Create `src/pages/MyPage.tsx`
2. Add route in `App.tsx`
3. Add navigation in `Sidebar.tsx`

### Add New API Endpoint
1. Add to `constants/index.ts`:
```tsx
export const API_ENDPOINTS = {
  // ...
  MY_ENDPOINT: '/api/my-endpoint',
};
```

2. Add service function in `services/backend.ts`:
```tsx
async myFunction() {
  const response = await apiClient.get(API_ENDPOINTS.MY_ENDPOINT);
  return response.data;
}
```

3. Use in component:
```tsx
const data = await myService.myFunction();
```

## 🐛 Debugging

### React DevTools
Install React DevTools browser extension

### Network Requests
Check browser DevTools → Network tab

### State Inspection
Zustand stores visible in React DevTools

### Console Logs
```tsx
console.log('Debug:', { user, files });
```

## 📚 Resources

- [React Docs](https://react.dev)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [Vite Guide](https://vitejs.dev/guide/)
- [Zustand Docs](https://github.com/pmndrs/zustand)

## 🆘 Troubleshooting

### Build Errors
```bash
# Clear cache and rebuild
rm -rf node_modules dist
npm install
npm run build
```

### Type Errors
```bash
# Regenerate type definitions
npm run type-check
```

### Port Already in Use
Change port in `vite.config.ts`:
```ts
server: {
  port: 3001, // Change here
}
```

## ✅ Checklist for Production

- [ ] Environment variables configured
- [ ] Build runs without errors
- [ ] All pages load correctly
- [ ] Authentication works
- [ ] API calls succeed
- [ ] Responsive on mobile
- [ ] Browser console has no errors
- [ ] Performance is acceptable

---

**For full project documentation, see [COMPLETE_GUIDE.md](../COMPLETE_GUIDE.md)**
