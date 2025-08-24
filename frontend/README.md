# URL Shortener Frontend

A modern React frontend for the URL Shortener service.

## Features

- **Dual API Key Management**: Create new API keys OR use existing ones
- **Simple Onboarding**: Create API keys with just your email
- **Existing User Support**: Input your existing API key to access your account
- **Link Creation**: Create short URLs with optional custom codes
- **Link Management**: View, copy, and delete your short URLs
- **Click Tracking**: Monitor usage statistics
- **Session Persistence**: API keys are remembered between visits
- **Modern UI**: Clean, responsive design with Tailwind CSS

## Setup

### Prerequisites
- Node.js 18+ 
- npm or yarn

### Installation
```bash
npm install
```

### Environment Configuration
Create a `.env` file in the frontend directory:
```env
VITE_API_BASE_URL=http://localhost:8080
```

### Development
```bash
npm run dev
```

The app will run on `http://localhost:5173`

### Build
```bash
npm run build
```

## Project Structure

```
src/
├── components/          # React components
│   ├── CreateApiKey.tsx    # Dual-mode API key management
│   ├── CreateLink.tsx      # Short URL creation form
│   ├── Header.tsx          # Navigation header
│   └── LinkList.tsx        # User's links display
├── services/            # API communication
│   └── api.ts              # API service layer
├── App.tsx              # Main application component
└── main.tsx             # Application entry point
```

## Usage Flow

### **New Users:**
1. **First Visit**: User sees welcome page with API key creation form
2. **Create API Key**: Enter email → receive API key
3. **Access App**: Automatically logged in with new API key

### **Existing Users:**
1. **Return Visit**: User can input existing API key
2. **Validate Key**: API key is validated against backend
3. **Access App**: Logged in with existing API key

### **All Users:**
1. **Create Links**: Make short URLs with optional custom codes
2. **Manage Links**: View, copy, delete links
3. **Session Persistence**: API key remembered between visits
4. **Logout**: Clear API key and return to welcome page

## API Key Management

### **Create New API Key**
- **Input**: Email address
- **Process**: Backend creates account and generates 32-character hex key
- **Result**: User logged in with new API key

### **Use Existing API Key**
- **Input**: 32-character API key
- **Process**: Key validated by attempting to fetch user data
- **Result**: User logged in with existing API key

### **Session Persistence**
- API keys stored in localStorage
- Automatically validated on app startup
- Invalid keys are automatically removed
- Users stay logged in between browser sessions

## API Integration

The frontend communicates with the backend using:
- **Create API Key**: `POST /create-api-key`
- **Create Link**: `POST /api/links`
- **Get Links**: `GET /api/links`
- **Delete Link**: `DELETE /api/links/:code`

## Styling

Built with Tailwind CSS for a modern, responsive design.

## Development Notes

- Uses React 19 with TypeScript
- Vite for fast development and building
- Axios for HTTP requests with automatic API key injection
- Local storage for API key persistence
- Responsive design for mobile and desktop
- Automatic API key validation and cleanup
