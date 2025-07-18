class ApiConstants {
  // Base URL - Update this to match your backend URL
  static const String baseUrl = 'http://localhost:8080/api/v1';
  
  // Auth endpoints
  static const String login = '/auth/login';
  static const String refresh = '/auth/refresh';
  static const String logout = '/auth/logout';
  
  // Core endpoints
  static const String users = '/users';
  static const String customers = '/customers';
  static const String vehicles = '/vehicles';
  static const String services = '/services';
  static const String products = '/products';
  static const String serviceJobs = '/service-jobs';
  static const String transactions = '/transactions';
  static const String payments = '/payments';
  
  // Master data endpoints
  static const String serviceCategories = '/master-data/service-categories';
  static const String productCategories = '/master-data/product-categories';
  static const String suppliers = '/master-data/suppliers';
  static const String unitTypes = '/master-data/unit-types';
  static const String paymentMethods = '/master-data/payment-methods';
}

class AppConstants {
  // App info
  static const String appName = 'Bengkel Management';
  static const String appVersion = '1.0.0';
  
  // Storage keys
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userDataKey = 'user_data';
  
  // Pagination
  static const int defaultPageSize = 10;
  static const int maxPageSize = 100;
  
  // Timeouts
  static const int connectTimeout = 30000; // 30 seconds
  static const int receiveTimeout = 30000; // 30 seconds
  
  // UI constants
  static const double cardElevation = 2.0;
  static const double borderRadius = 8.0;
  static const double defaultPadding = 16.0;
}

class AppStrings {
  // Auth
  static const String login = 'Login';
  static const String logout = 'Logout';
  static const String username = 'Username';
  static const String password = 'Password';
  static const String forgotPassword = 'Forgot Password?';
  static const String loginSuccess = 'Login successful';
  static const String loginFailed = 'Login failed';
  static const String invalidCredentials = 'Invalid username or password';
  
  // Navigation
  static const String dashboard = 'Dashboard';
  static const String customers = 'Customers';
  static const String vehicles = 'Vehicles';
  static const String serviceJobs = 'Service Jobs';
  static const String inventory = 'Inventory';
  static const String financial = 'Financial';
  static const String reports = 'Reports';
  static const String settings = 'Settings';
  
  // Common
  static const String save = 'Save';
  static const String cancel = 'Cancel';
  static const String delete = 'Delete';
  static const String edit = 'Edit';
  static const String add = 'Add';
  static const String search = 'Search';
  static const String filter = 'Filter';
  static const String loading = 'Loading...';
  static const String noData = 'No data available';
  static const String error = 'Error';
  static const String success = 'Success';
  static const String warning = 'Warning';
  static const String info = 'Info';
  
  // Status
  static const String pending = 'Pending';
  static const String inProgress = 'In Progress';
  static const String completed = 'Completed';
  static const String cancelled = 'Cancelled';
  static const String onHold = 'On Hold';
  
  // Validation
  static const String fieldRequired = 'This field is required';
  static const String invalidEmail = 'Please enter a valid email';
  static const String invalidPhone = 'Please enter a valid phone number';
  static const String passwordTooShort = 'Password must be at least 6 characters';
}