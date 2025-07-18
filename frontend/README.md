# Workshop Management System - Flutter Frontend

A comprehensive Flutter mobile application for workshop management built with Material Design 3, Riverpod state management, and responsive design.

## Features

- **Responsive Design**: Optimized for mobile, tablet, and desktop
- **Material Design 3**: Modern UI with dynamic theming
- **State Management**: Riverpod for efficient state management
- **Navigation**: GoRouter for declarative routing
- **Authentication**: JWT-based authentication with secure storage
- **Multi-Business Support**: Service workshop, sparepart sales, vehicle trading

## Architecture

### Project Structure
```
lib/
├── main.dart                    # App entry point
├── app/                         # App configuration
│   ├── app.dart                # Main app widget
│   ├── theme/                  # Material Design 3 themes
│   └── routes/                 # GoRouter configuration
├── core/                       # Core functionality
│   ├── constants/              # App constants
│   ├── services/               # Core services (auth, etc.)
│   ├── utils/                  # Utility functions
│   └── extensions/             # Dart extensions
├── data/                       # Data layer
│   ├── models/                 # Data models
│   ├── repositories/           # Repository pattern
│   └── datasources/            # API & local data sources
├── presentation/               # UI layer
│   ├── pages/                  # Screen widgets
│   ├── widgets/                # Reusable widgets
│   └── providers/              # Riverpod providers
└── shared/                     # Shared components
    ├── widgets/                # Common widgets
    └── components/             # UI components
```

### Key Technologies

- **Flutter**: 3.13.0+
- **Dart**: 3.1.0+
- **State Management**: Riverpod 2.4.9
- **Navigation**: GoRouter 12.1.3
- **HTTP Client**: Dio 5.4.0
- **Local Storage**: Hive + Shared Preferences
- **Charts**: FL Chart 0.66.2
- **Responsive Design**: ScreenUtil 5.9.0

## Getting Started

### Prerequisites

- Flutter SDK (3.13.0 or higher)
- Dart SDK (3.1.0 or higher)
- Android Studio / VS Code
- Android/iOS device or emulator

### Installation

1. Clone the repository
2. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

3. Install dependencies:
   ```bash
   flutter pub get
   ```

4. Run code generation (if needed):
   ```bash
   flutter packages pub run build_runner build
   ```

5. Run the app:
   ```bash
   flutter run
   ```

## Configuration

### Backend Connection

Update the API base URL in `lib/core/constants/app_constants.dart`:

```dart
static const String baseUrl = 'http://your-backend-url:8080/api/v1';
```

### Demo Credentials

- **Username**: admin
- **Password**: admin123

## Features

### Authentication
- JWT-based authentication
- Secure token storage
- Auto-logout on token expiry
- Role-based access control

### Dashboard
- Welcome screen with user info
- Quick stats overview
- Quick action shortcuts
- Responsive grid layout

### Navigation
- Responsive drawer navigation
- Bottom navigation for mobile
- Side navigation for tablets/desktop
- Consistent navigation patterns

### UI/UX
- Material Design 3 components
- Responsive layout system
- Dark/light theme support
- Smooth animations and transitions
- Professional color scheme

### State Management
- Riverpod providers for state
- Reactive UI updates
- Efficient rebuilds
- Type-safe state access

## Modules (Planned)

### Customer Management
- Customer registration and profiles
- Vehicle registration with details
- Service history tracking
- Customer search and filtering

### Service Job Management
- Service job creation and tracking
- Queue management system
- Technician assignment
- Real-time status updates
- Progress tracking

### Inventory Management
- Product catalog with categories
- Stock level monitoring
- Low stock alerts
- Supplier management
- Purchase order workflow

### Financial Management
- Transaction recording
- Payment method handling
- Invoice generation
- Financial reporting

### Vehicle Trading
- Vehicle listing and management
- Purchase/sales tracking
- Condition assessment
- Price history

### Reports & Analytics
- Business performance charts
- Sales and service reports
- Inventory reports
- Financial statements
- Export to PDF/Excel

## Development

### Code Generation

This project uses code generation for models and providers:

```bash
# Run code generation
flutter packages pub run build_runner build

# Watch for changes
flutter packages pub run build_runner watch
```

### Linting

The project follows strict linting rules for code quality:

```bash
flutter analyze
```

### Testing

Run tests with:

```bash
flutter test
```

## Building

### Android
```bash
flutter build apk --release
```

### iOS
```bash
flutter build ios --release
```

### Web
```bash
flutter build web --release
```

## Responsive Design

The app is designed to work seamlessly across different screen sizes:

- **Mobile** (< 600px): Bottom navigation, single column layout
- **Tablet** (600px - 1200px): Side navigation, two-column layout
- **Desktop** (> 1200px): Full sidebar, multi-column layout

## Contributing

1. Follow the established project structure
2. Use consistent naming conventions
3. Add proper documentation
4. Write tests for new features
5. Follow Material Design 3 guidelines

## License

MIT License - see LICENSE file for details