import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../presentation/pages/auth/login_page.dart';
import '../../presentation/pages/dashboard/dashboard_page.dart';
import '../../presentation/pages/customers/customers_page.dart';
import '../../presentation/pages/vehicles/vehicles_page.dart';
import '../../presentation/pages/service_jobs/service_jobs_page.dart';
import '../../presentation/pages/inventory/inventory_page.dart';
import '../../presentation/pages/financial/financial_page.dart';
import '../../presentation/pages/reports/reports_page.dart';
import '../../presentation/pages/settings/settings_page.dart';
import '../../core/services/auth_service.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  final authService = ref.watch(authServiceProvider);
  
  return GoRouter(
    initialLocation: '/login',
    redirect: (context, state) {
      final isLoggedIn = authService.isAuthenticated;
      final isLoginRoute = state.location == '/login';
      
      // If not logged in and not on login page, redirect to login
      if (!isLoggedIn && !isLoginRoute) {
        return '/login';
      }
      
      // If logged in and on login page, redirect to dashboard
      if (isLoggedIn && isLoginRoute) {
        return '/dashboard';
      }
      
      return null; // No redirect needed
    },
    routes: [
      // Auth routes
      GoRoute(
        path: '/login',
        name: 'login',
        builder: (context, state) => const LoginPage(),
      ),
      
      // Main app routes
      GoRoute(
        path: '/dashboard',
        name: 'dashboard',
        builder: (context, state) => const DashboardPage(),
      ),
      
      // Customer management
      GoRoute(
        path: '/customers',
        name: 'customers',
        builder: (context, state) => const CustomersPage(),
      ),
      
      // Vehicle management
      GoRoute(
        path: '/vehicles',
        name: 'vehicles',
        builder: (context, state) => const VehiclesPage(),
      ),
      
      // Service jobs
      GoRoute(
        path: '/service-jobs',
        name: 'service-jobs',
        builder: (context, state) => const ServiceJobsPage(),
      ),
      
      // Inventory management
      GoRoute(
        path: '/inventory',
        name: 'inventory',
        builder: (context, state) => const InventoryPage(),
      ),
      
      // Financial management
      GoRoute(
        path: '/financial',
        name: 'financial',
        builder: (context, state) => const FinancialPage(),
      ),
      
      // Reports
      GoRoute(
        path: '/reports',
        name: 'reports',
        builder: (context, state) => const ReportsPage(),
      ),
      
      // Settings
      GoRoute(
        path: '/settings',
        name: 'settings',
        builder: (context, state) => const SettingsPage(),
      ),
    ],
  );
});

// Navigation helper
class AppRouter {
  static void toLogin(BuildContext context) {
    context.goNamed('login');
  }
  
  static void toDashboard(BuildContext context) {
    context.goNamed('dashboard');
  }
  
  static void toCustomers(BuildContext context) {
    context.goNamed('customers');
  }
  
  static void toVehicles(BuildContext context) {
    context.goNamed('vehicles');
  }
  
  static void toServiceJobs(BuildContext context) {
    context.goNamed('service-jobs');
  }
  
  static void toInventory(BuildContext context) {
    context.goNamed('inventory');
  }
  
  static void toFinancial(BuildContext context) {
    context.goNamed('financial');
  }
  
  static void toReports(BuildContext context) {
    context.goNamed('reports');
  }
  
  static void toSettings(BuildContext context) {
    context.goNamed('settings');
  }
}