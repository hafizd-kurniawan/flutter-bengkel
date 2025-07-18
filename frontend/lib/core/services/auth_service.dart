import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../constants/app_constants.dart';
import '../../data/models/user_model.dart';

// Auth service provider
final authServiceProvider = Provider<AuthService>((ref) {
  return AuthService();
});

class AuthService {
  static const String _accessTokenKey = AppConstants.accessTokenKey;
  static const String _refreshTokenKey = AppConstants.refreshTokenKey;
  static const String _userDataKey = AppConstants.userDataKey;
  
  // Current user state
  UserModel? _currentUser;
  String? _accessToken;
  String? _refreshToken;
  
  // Getters
  UserModel? get currentUser => _currentUser;
  String? get accessToken => _accessToken;
  String? get refreshToken => _refreshToken;
  bool get isAuthenticated => _accessToken != null && _currentUser != null;
  
  // Initialize auth service (load stored tokens)
  Future<void> initialize() async {
    final prefs = await SharedPreferences.getInstance();
    
    _accessToken = prefs.getString(_accessTokenKey);
    _refreshToken = prefs.getString(_refreshTokenKey);
    
    final userDataJson = prefs.getString(_userDataKey);
    if (userDataJson != null) {
      try {
        // _currentUser = UserModel.fromJson(json.decode(userDataJson));
        // For now, create a dummy user
        _currentUser = UserModel(
          id: 1,
          username: 'admin',
          email: 'admin@bengkel.com',
          fullName: 'System Administrator',
          roleName: 'Super Admin',
        );
      } catch (e) {
        // Handle parsing error
        await clearTokens();
      }
    }
  }
  
  // Login
  Future<bool> login(String username, String password) async {
    try {
      // TODO: Implement actual API call
      // For demo purposes, accept admin/admin123
      if (username == 'admin' && password == 'admin123') {
        _accessToken = 'dummy_access_token';
        _refreshToken = 'dummy_refresh_token';
        _currentUser = UserModel(
          id: 1,
          username: username,
          email: 'admin@bengkel.com',
          fullName: 'System Administrator',
          roleName: 'Super Admin',
        );
        
        await _storeTokens();
        return true;
      }
      
      return false;
    } catch (e) {
      return false;
    }
  }
  
  // Logout
  Future<void> logout() async {
    await clearTokens();
    _currentUser = null;
    _accessToken = null;
    _refreshToken = null;
  }
  
  // Store tokens in secure storage
  Future<void> _storeTokens() async {
    final prefs = await SharedPreferences.getInstance();
    
    if (_accessToken != null) {
      await prefs.setString(_accessTokenKey, _accessToken!);
    }
    
    if (_refreshToken != null) {
      await prefs.setString(_refreshTokenKey, _refreshToken!);
    }
    
    if (_currentUser != null) {
      // await prefs.setString(_userDataKey, json.encode(_currentUser!.toJson()));
      // For now, just store a dummy value
      await prefs.setString(_userDataKey, 'dummy_user_data');
    }
  }
  
  // Clear stored tokens
  Future<void> clearTokens() async {
    final prefs = await SharedPreferences.getInstance();
    
    await prefs.remove(_accessTokenKey);
    await prefs.remove(_refreshTokenKey);
    await prefs.remove(_userDataKey);
  }
  
  // Refresh access token
  Future<bool> refreshAccessToken() async {
    try {
      // TODO: Implement token refresh API call
      return false;
    } catch (e) {
      return false;
    }
  }
}