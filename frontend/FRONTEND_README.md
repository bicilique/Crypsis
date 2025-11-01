# E-Crypt Frontend 🔐

A modern, professional React-based frontend for the E-Crypt secure file encryption and storage service. Built with TypeScript, Tailwind CSS, and designed with enterprise-grade security in mind.

## 🚀 Features

### 🔒 Security-First Design
- **Modern Authentication**: JWT-based authentication with automatic token refresh
- **Role-Based Access Control**: Admin and user role management
- **Secure API Communication**: Axios interceptors with automatic retry and error handling
- **Protected Routes**: Authentication guards for all sensitive pages

### 💼 Professional UI/UX
- **Aikido.dev Inspired Design**: Clean, minimalist security-focused interface
- **Responsive Layout**: Mobile-first design with collapsible navigation
- **Accessibility Compliant**: WCAG 2.1 AA standards with proper ARIA labels
- **Modern Components**: Reusable UI components with consistent styling

### 📁 File Management
- **Drag & Drop Upload**: Intuitive file upload with progress tracking
- **Encryption Support**: Visual indicators for encrypted files
- **Bulk Operations**: Select and manage multiple files
- **File Preview**: Support for various file types (planned)
- **Search & Filter**: Advanced filtering and search capabilities

### ⚙️ Admin Panel
- **User Management**: Create, update, and manage admin users
- **OAuth2 Applications**: Register and manage client applications
- **System Monitoring**: Real-time dashboard with key metrics
- **Audit Logs**: Comprehensive activity tracking
- **Security Alerts**: System security status monitoring

## 🛠️ Tech Stack

### Frontend Framework
- **React 19** with TypeScript for type safety
- **Vite** for fast development and building
- **React Router v6** for client-side routing

### State Management
- **Zustand** for lightweight state management
- **Persistent storage** for authentication state

### Styling & UI
- **Tailwind CSS** for utility-first styling
- **Headless UI** for accessible component primitives
- **Heroicons** for consistent iconography
- **Inter Font** for professional typography

### Forms & Validation
- **React Hook Form** for performant form handling
- **Zod** for schema validation
- **Real-time validation** with error feedback

### HTTP & API
- **Axios** with interceptors for API communication
- **Automatic token refresh** and error handling
- **Request/Response transformation**

## 📁 Project Structure

```
src/
├── components/           # Reusable UI components
│   ├── ui/              # Base UI components
│   │   ├── Button.tsx   # Reusable button component
│   │   ├── Card.tsx     # Container component
│   │   ├── Input.tsx    # Form input component
│   │   ├── Alert.tsx    # Alert/notification component
│   │   └── index.ts     # Component exports
│   ├── layout/          # Layout components
│   │   ├── AppLayout.tsx      # Main app layout
│   │   ├── Header.tsx         # Top navigation
│   │   ├── Sidebar.tsx        # Side navigation
│   │   ├── ProtectedRoute.tsx # Route protection
│   │   └── index.ts           # Layout exports
│   └── features/        # Feature-specific components
├── pages/               # Page components
│   ├── LoginPage.tsx    # Authentication page
│   ├── DashboardPage.tsx # Main dashboard
│   └── index.ts         # Page exports
├── stores/              # State management
│   ├── auth.ts         # Authentication state
│   ├── files.ts        # File management state
│   ├── admin.ts        # Admin operations state
│   └── index.ts        # Store exports
├── services/            # API services
│   ├── api.ts          # Base API client
│   ├── auth.ts         # Authentication service
│   ├── files.ts        # File operations service
│   ├── admin.ts        # Admin operations service
│   └── index.ts        # Service exports
├── types/               # TypeScript definitions
│   └── index.ts        # Type exports
├── utils/               # Utility functions
│   └── index.ts        # Utility exports
├── constants/           # App constants
│   └── index.ts        # Constants exports
├── App.tsx             # Main app component
├── main.jsx            # App entry point
└── index.css           # Global styles
```

## 🚦 Getting Started

### Prerequisites
- Node.js 18+ and npm
- E-Crypt backend API running on port 8080

### Installation

1. **Navigate to the project:**
   ```bash
   cd crypsis
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start the development server:**
   ```bash
   npm run dev
   ```

4. **Open your browser:**
   Navigate to `http://localhost:5173`

### Default Login
- **Username:** `admin`
- **Password:** `password123`

## 🎨 Design System

### Color Palette
- **Primary Blue:** `#3b82f6` - Trust and security
- **Dark Blue:** `#1e40af` - Professional depth
- **Success Green:** `#10b981` - Successful operations
- **Warning Orange:** `#f59e0b` - Important notices
- **Error Red:** `#ef4444` - Critical issues

### Typography
- **Font Family:** Inter (Google Fonts)
- **Responsive sizing:** Mobile-first approach
- **Readable hierarchy:** Clear information structure

## 🔧 Development

### Available Scripts

```bash
# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Run linting
npm run lint
```

## 🔐 Authentication Flow

1. **Login:** User submits credentials
2. **Token Storage:** JWT tokens stored securely
3. **API Requests:** Automatic token attachment
4. **Auto Refresh:** Seamless token renewal
5. **Route Protection:** Unauthorized access prevention

## 📱 Responsive Design

### Breakpoints
- **Mobile:** 0-640px
- **Tablet:** 641-1024px
- **Desktop:** 1025px+

### Features
- **Collapsible sidebar** on mobile
- **Touch-friendly** interaction targets
- **Optimized layouts** for all screen sizes

## 🛡️ Security Features

### Frontend Security
- **XSS Prevention:** Input sanitization
- **CSRF Protection:** Request validation
- **Secure Storage:** Encrypted local storage
- **Content Security Policy:** Header protection

### API Security
- **Token Validation:** Server-side verification
- **Request Encryption:** HTTPS only
- **Rate Limiting:** Abuse prevention
- **Audit Logging:** Activity tracking

## 📄 API Integration

### Endpoints Used
```
Authentication:
POST /api/admin/login
GET  /api/admin/logout
GET  /api/admin/refresh-token

File Operations:
POST   /api/files
GET    /api/files/list
GET    /api/files/:id/download
DELETE /api/files/:id/delete

Admin Operations:
GET    /api/admin/list
POST   /api/admin/add
GET    /api/admin/apps
POST   /api/admin/apps
```

## 🤝 Contributing

### Development Workflow
1. Create feature branch
2. Implement changes with tests
3. Ensure TypeScript compliance
4. Update documentation
5. Submit pull request

### Code Style
- Follow TypeScript best practices
- Use meaningful component names
- Add JSDoc comments for functions
- Implement proper error handling

## 📚 Learning Resources

### For Beginners
- **React Documentation:** https://react.dev
- **TypeScript Handbook:** https://www.typescriptlang.org/docs
- **Tailwind CSS:** https://tailwindcss.com/docs
- **Zustand Guide:** https://zustand-demo.pmnd.rs

### Architecture Concepts
- **Clean Architecture:** Separation of concerns
- **Component Composition:** Reusable building blocks
- **State Management:** Centralized application state
- **API Integration:** Service layer pattern

---

**Built with ❤️ for enterprise security**

This frontend provides a solid foundation for the E-Crypt file encryption service with room for growth and customization based on specific requirements.
