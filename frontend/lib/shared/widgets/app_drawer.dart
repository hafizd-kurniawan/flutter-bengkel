import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gap/gap.dart';

import '../../core/services/auth_service.dart';
import '../../core/constants/app_constants.dart';
import '../../app/routes/app_router.dart';

class AppDrawer extends ConsumerWidget {
  const AppDrawer({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authService = ref.watch(authServiceProvider);
    final user = authService.currentUser;

    return Drawer(
      child: Column(
        children: [
          // User info header
          UserAccountsDrawerHeader(
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primary,
            ),
            accountName: Text(
              user?.fullName ?? 'User',
              style: const TextStyle(
                fontWeight: FontWeight.bold,
              ),
            ),
            accountEmail: Text(user?.email ?? ''),
            currentAccountPicture: CircleAvatar(
              backgroundColor: Colors.white,
              child: Text(
                user?.fullName.substring(0, 1).toUpperCase() ?? 'U',
                style: TextStyle(
                  fontSize: 20.sp,
                  fontWeight: FontWeight.bold,
                  color: Theme.of(context).colorScheme.primary,
                ),
              ),
            ),
          ),
          
          // Navigation items
          Expanded(
            child: ListView(
              padding: EdgeInsets.zero,
              children: [
                _buildDrawerItem(
                  context,
                  Icons.dashboard,
                  AppStrings.dashboard,
                  () => AppRouter.toDashboard(context),
                ),
                _buildDrawerItem(
                  context,
                  Icons.people,
                  AppStrings.customers,
                  () => AppRouter.toCustomers(context),
                ),
                _buildDrawerItem(
                  context,
                  Icons.directions_car,
                  AppStrings.vehicles,
                  () => AppRouter.toVehicles(context),
                ),
                _buildDrawerItem(
                  context,
                  Icons.build,
                  AppStrings.serviceJobs,
                  () => AppRouter.toServiceJobs(context),
                ),
                _buildDrawerItem(
                  context,
                  Icons.inventory,
                  AppStrings.inventory,
                  () => AppRouter.toInventory(context),
                ),
                _buildDrawerItem(
                  context,
                  Icons.account_balance_wallet,
                  AppStrings.financial,
                  () => AppRouter.toFinancial(context),
                ),
                _buildDrawerItem(
                  context,
                  Icons.analytics,
                  AppStrings.reports,
                  () => AppRouter.toReports(context),
                ),
                const Divider(),
                _buildDrawerItem(
                  context,
                  Icons.settings,
                  AppStrings.settings,
                  () => AppRouter.toSettings(context),
                ),
              ],
            ),
          ),
          
          // Logout button
          Padding(
            padding: EdgeInsets.all(16.w),
            child: ElevatedButton.icon(
              onPressed: () async {
                await authService.logout();
                if (context.mounted) {
                  Navigator.of(context).pop(); // Close drawer
                  AppRouter.toLogin(context);
                }
              },
              icon: const Icon(Icons.logout),
              label: const Text(AppStrings.logout),
              style: ElevatedButton.styleFrom(
                backgroundColor: Theme.of(context).colorScheme.error,
                foregroundColor: Theme.of(context).colorScheme.onError,
              ),
            ),
          ),
          
          // App version
          Padding(
            padding: EdgeInsets.only(bottom: 16.h),
            child: Text(
              'Version ${AppConstants.appVersion}',
              style: TextStyle(
                fontSize: 12.sp,
                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.5),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDrawerItem(
    BuildContext context,
    IconData icon,
    String title,
    VoidCallback onTap,
  ) {
    return ListTile(
      leading: Icon(icon),
      title: Text(title),
      onTap: () {
        Navigator.of(context).pop(); // Close drawer
        onTap();
      },
    );
  }
}