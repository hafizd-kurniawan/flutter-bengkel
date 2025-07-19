import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gap/gap.dart';

import '../../../core/services/auth_service.dart';
import '../../../core/constants/app_constants.dart';
import '../../../app/routes/app_router.dart';
import '../../../shared/widgets/app_drawer.dart';

class DashboardPage extends ConsumerWidget {
  const DashboardPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authService = ref.watch(authServiceProvider);
    final user = authService.currentUser;

    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.dashboard),
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications),
            onPressed: () {
              // TODO: Show notifications
            },
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () async {
              await authService.logout();
              if (context.mounted) {
                AppRouter.toLogin(context);
              }
            },
          ),
        ],
      ),
      drawer: const AppDrawer(),
      body: SingleChildScrollView(
        padding: EdgeInsets.all(16.w),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Welcome card
            Card(
              child: Padding(
                padding: EdgeInsets.all(16.w),
                child: Row(
                  children: [
                    CircleAvatar(
                      radius: 30.r,
                      backgroundColor: Theme.of(context).colorScheme.primary,
                      child: Text(
                        user?.fullName.substring(0, 1).toUpperCase() ?? 'U',
                        style: TextStyle(
                          fontSize: 24.sp,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                      ),
                    ),
                    Gap(16.w),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Welcome back!',
                            style: TextStyle(
                              fontSize: 16.sp,
                              color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
                            ),
                          ),
                          Text(
                            user?.fullName ?? 'User',
                            style: TextStyle(
                              fontSize: 20.sp,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          Text(
                            user?.roleName ?? 'User',
                            style: TextStyle(
                              fontSize: 14.sp,
                              color: Theme.of(context).colorScheme.primary,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            Gap(24.h),
            
            // Quick stats
            Text(
              'Quick Stats',
              style: TextStyle(
                fontSize: 20.sp,
                fontWeight: FontWeight.bold,
              ),
            ),
            Gap(16.h),
            
            GridView.count(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisCount: 4,
              childAspectRatio: 1.2,
              crossAxisSpacing: 16.w,
              mainAxisSpacing: 16.h,
              children: [
                _buildStatCard(
                  context,
                  'Service Jobs',
                  '24',
                  Icons.build,
                  Colors.blue,
                ),
                _buildStatCard(
                  context,
                  'Customers',
                  '156',
                  Icons.people,
                  Colors.green,
                ),
                _buildStatCard(
                  context,
                  'Vehicles Stock',
                  '98',
                  Icons.directions_car,
                  Colors.orange,
                ),
                _buildStatCard(
                  context,
                  'Sales Today',
                  '5',
                  Icons.sell,
                  Colors.purple,
                ),
                _buildStatCard(
                  context,
                  'Low Stock Items',
                  '12',
                  Icons.warning,
                  Colors.red,
                ),
                _buildStatCard(
                  context,
                  'Profit This Month',
                  'Rp 125M',
                  Icons.trending_up,
                  Colors.teal,
                ),
                _buildStatCard(
                  context,
                  'Commission Due',
                  'Rp 8.5M',
                  Icons.payments,
                  Colors.indigo,
                ),
                _buildStatCard(
                  context,
                  'Pending Jobs',
                  '7',
                  Icons.pending,
                  Colors.amber,
                ),
              ],
            ),
            Gap(24.h),
            
            // Quick actions
            Text(
              'Quick Actions',
              style: TextStyle(
                fontSize: 20.sp,
                fontWeight: FontWeight.bold,
              ),
            ),
            Gap(16.h),
            
            GridView.count(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisCount: 4,
              childAspectRatio: 2,
              crossAxisSpacing: 16.w,
              mainAxisSpacing: 16.h,
              children: [
                _buildActionCard(
                  context,
                  'New Service Job',
                  Icons.add_business,
                  () => AppRouter.toServiceJobs(context),
                ),
                _buildActionCard(
                  context,
                  'Buy Vehicle',
                  Icons.car_rental,
                  () => AppRouter.toVehicles(context),
                ),
                _buildActionCard(
                  context,
                  'Sell Vehicle',
                  Icons.sell,
                  () => AppRouter.toVehicles(context),
                ),
                _buildActionCard(
                  context,
                  'Add Customer',
                  Icons.person_add,
                  () => AppRouter.toCustomers(context),
                ),
                _buildActionCard(
                  context,
                  'Check Inventory',
                  Icons.inventory,
                  () => AppRouter.toInventory(context),
                ),
                _buildActionCard(
                  context,
                  'View Reports',
                  Icons.analytics,
                  () => AppRouter.toReports(context),
                ),
                _buildActionCard(
                  context,
                  'Vehicle Photos',
                  Icons.camera_alt,
                  () => AppRouter.toVehicles(context),
                ),
                _buildActionCard(
                  context,
                  'Profit Analysis',
                  Icons.trending_up,
                  () => AppRouter.toReports(context),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatCard(
    BuildContext context,
    String title,
    String value,
    IconData icon,
    Color color,
  ) {
    return Card(
      child: Padding(
        padding: EdgeInsets.all(16.w),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              size: 32.w,
              color: color,
            ),
            Gap(8.h),
            Text(
              value,
              style: TextStyle(
                fontSize: 24.sp,
                fontWeight: FontWeight.bold,
                color: color,
              ),
            ),
            Text(
              title,
              style: TextStyle(
                fontSize: 12.sp,
                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildActionCard(
    BuildContext context,
    String title,
    IconData icon,
    VoidCallback onTap,
  ) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12.r),
        child: Padding(
          padding: EdgeInsets.all(16.w),
          child: Row(
            children: [
              Icon(
                icon,
                size: 24.w,
                color: Theme.of(context).colorScheme.primary,
              ),
              Gap(12.w),
              Expanded(
                child: Text(
                  title,
                  style: TextStyle(
                    fontSize: 14.sp,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
              Icon(
                Icons.arrow_forward_ios,
                size: 16.w,
                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.5),
              ),
            ],
          ),
        ),
      ),
    );
  }
}