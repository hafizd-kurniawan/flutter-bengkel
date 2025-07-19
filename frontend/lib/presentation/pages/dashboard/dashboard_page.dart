import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gap/gap.dart';
import 'package:intl/intl.dart';

import '../../../core/services/auth_service.dart';
import '../../../core/constants/app_constants.dart';
import '../../../app/routes/app_router.dart';
import '../../../shared/widgets/app_drawer.dart';
import '../../../data/services/vehicle_service.dart';

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
              _showNotifications(context, ref);
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
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        Text(
                          DateFormat('EEEE, MMM dd', 'id_ID').format(DateTime.now()),
                          style: TextStyle(
                            fontSize: 14.sp,
                            color: Colors.grey[600],
                          ),
                        ),
                        Text(
                          DateFormat('HH:mm').format(DateTime.now()),
                          style: TextStyle(
                            fontSize: 12.sp,
                            color: Colors.grey[500],
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            Gap(24.h),
            
            // Real-time stats
            Text(
              'Real-Time Statistics',
              style: TextStyle(
                fontSize: 20.sp,
                fontWeight: FontWeight.bold,
              ),
            ),
            Gap(16.h),
            
            // Live data from vehicle service
            Consumer(
              builder: (context, ref, child) {
                final vehicleStatsAsync = ref.watch(vehicleStatisticsProvider);
                final vehicleInventoryAsync = ref.watch(vehicleInventoryProvider({'status': 'Available'}));
                final vehicleSalesAsync = ref.watch(vehicleSalesProvider({}));
                
                return vehicleStatsAsync.when(
                  data: (stats) => GridView.count(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    crossAxisCount: 4,
                    childAspectRatio: 1.1,
                    crossAxisSpacing: 16.w,
                    mainAxisSpacing: 16.h,
                    children: [
                      _buildStatCard(
                        context,
                        'Vehicle Stock',
                        stats['total_stock'].toString(),
                        Icons.directions_car,
                        Colors.blue,
                        () => AppRouter.toVehicles(context),
                      ),
                      _buildStatCard(
                        context,
                        'Available',
                        stats['available_stock'].toString(),
                        Icons.check_circle,
                        Colors.green,
                        () => AppRouter.toVehicles(context),
                      ),
                      _buildStatCard(
                        context,
                        'Sold This Month',
                        stats['sold_this_month'].toString(),
                        Icons.trending_up,
                        Colors.orange,
                        () => AppRouter.toVehicles(context),
                      ),
                      _buildStatCard(
                        context,
                        'Stock Value',
                        NumberFormat.compact(locale: 'id_ID').format(stats['total_stock_value']),
                        Icons.attach_money,
                        Colors.purple,
                        () => AppRouter.toReports(context),
                      ),
                      
                      // Additional mock stats for workshop operations
                      _buildStatCard(
                        context,
                        'Service Queue',
                        '8',
                        Icons.build,
                        Colors.indigo,
                        () => AppRouter.toServiceJobs(context),
                      ),
                      _buildStatCard(
                        context,
                        'Customers',
                        '156',
                        Icons.people,
                        Colors.teal,
                        () => AppRouter.toCustomers(context),
                      ),
                      _buildStatCard(
                        context,
                        'Today Revenue',
                        'Rp 45M',
                        Icons.payments,
                        Colors.red,
                        () => AppRouter.toFinancial(context),
                      ),
                      _buildStatCard(
                        context,
                        'Low Stock Alert',
                        '3',
                        Icons.warning,
                        Colors.amber,
                        () => AppRouter.toInventory(context),
                      ),
                    ],
                  ),
                  loading: () => _buildLoadingStats(),
                  error: (error, stack) => _buildErrorStats(context),
                );
              },
            ),
            Gap(24.h),
            
            // Recent Activities
            Text(
              'Recent Activities',
              style: TextStyle(
                fontSize: 20.sp,
                fontWeight: FontWeight.bold,
              ),
            ),
            Gap(16.h),
            
            Consumer(
              builder: (context, ref, child) {
                final salesAsync = ref.watch(vehicleSalesProvider({}));
                final purchasesAsync = ref.watch(vehiclePurchasesProvider({}));
                
                return Card(
                  child: Padding(
                    padding: EdgeInsets.all(16.w),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Latest Transactions',
                          style: TextStyle(
                            fontSize: 16.sp,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        Gap(12.h),
                        
                        salesAsync.when(
                          data: (sales) {
                            if (sales.isEmpty) {
                              return const Text('No recent sales');
                            }
                            return Column(
                              children: sales.take(3).map((sale) => 
                                _buildActivityItem(
                                  'Vehicle Sold',
                                  '${sale.customerName} - ${sale.formattedPrice}',
                                  Icons.sell,
                                  Colors.green,
                                  sale.saleDate,
                                )
                              ).toList(),
                            );
                          },
                          loading: () => const CircularProgressIndicator(),
                          error: (error, stack) => Text('Error: $error'),
                        ),
                        
                        purchasesAsync.when(
                          data: (purchases) {
                            if (purchases.isEmpty) {
                              return const SizedBox();
                            }
                            return Column(
                              children: purchases.take(2).map((purchase) => 
                                _buildActivityItem(
                                  'Vehicle Purchased',
                                  '${purchase.customerName} - ${purchase.formattedPrice}',
                                  Icons.shopping_cart,
                                  Colors.blue,
                                  purchase.purchaseDate,
                                )
                              ).toList(),
                            );
                          },
                          loading: () => const SizedBox(),
                          error: (error, stack) => const SizedBox(),
                        ),
                        
                        // Mock service activities
                        _buildActivityItem(
                          'Service Completed',
                          'B 1234 ABC - Oil Change Service',
                          Icons.build,
                          Colors.orange,
                          DateTime.now().subtract(const Duration(hours: 2)),
                        ),
                        _buildActivityItem(
                          'Customer Added',
                          'John Doe - New customer registration',
                          Icons.person_add,
                          Colors.purple,
                          DateTime.now().subtract(const Duration(hours: 4)),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
            Gap(24.h),
            
            // Quick Actions
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
              crossAxisCount: 2,
              childAspectRatio: 3,
              crossAxisSpacing: 16.w,
              mainAxisSpacing: 16.h,
              children: [
                _buildActionCard(
                  context,
                  'New Service Job',
                  Icons.build,
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
                  'Financial Summary',
                  Icons.account_balance,
                  () => AppRouter.toFinancial(context),
                ),
                _buildActionCard(
                  context,
                  'Settings',
                  Icons.settings,
                  () => AppRouter.toSettings(context),
                ),
              ],
            ),
            Gap(24.h),
            
            // Performance Summary
            Text(
              'Performance Summary',
              style: TextStyle(
                fontSize: 20.sp,
                fontWeight: FontWeight.bold,
              ),
            ),
            Gap(16.h),
            
            Consumer(
              builder: (context, ref, child) {
                final salesAsync = ref.watch(vehicleSalesProvider({}));
                
                return Card(
                  child: Padding(
                    padding: EdgeInsets.all(16.w),
                    child: salesAsync.when(
                      data: (sales) {
                        final totalProfit = sales.fold(0.0, (sum, sale) => sum + sale.profitAmount);
                        final avgDaysToSell = 18; // Mock calculation
                        final topPerformer = 'Ahmad (Sales)';
                        
                        return Column(
                          children: [
                            Row(
                              children: [
                                Expanded(
                                  child: _buildPerformanceMetric(
                                    'Total Profit',
                                    NumberFormat.currency(
                                      locale: 'id_ID',
                                      symbol: 'Rp ',
                                      decimalDigits: 0,
                                    ).format(totalProfit),
                                    Icons.trending_up,
                                    Colors.green,
                                  ),
                                ),
                                Expanded(
                                  child: _buildPerformanceMetric(
                                    'Avg. Days to Sell',
                                    '$avgDaysToSell days',
                                    Icons.schedule,
                                    Colors.blue,
                                  ),
                                ),
                              ],
                            ),
                            Gap(16.h),
                            Row(
                              children: [
                                Expanded(
                                  child: _buildPerformanceMetric(
                                    'Top Performer',
                                    topPerformer,
                                    Icons.star,
                                    Colors.amber,
                                  ),
                                ),
                                Expanded(
                                  child: _buildPerformanceMetric(
                                    'Success Rate',
                                    '94%',
                                    Icons.check_circle,
                                    Colors.green,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        );
                      },
                      loading: () => const CircularProgressIndicator(),
                      error: (error, stack) => Text('Error loading performance data'),
                    ),
                  ),
                );
              },
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
    VoidCallback onTap,
  ) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12.r),
        child: Padding(
          padding: EdgeInsets.all(16.w),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                icon,
                size: 28.w,
                color: color,
              ),
              Gap(8.h),
              Text(
                value,
                style: TextStyle(
                  fontSize: 20.sp,
                  fontWeight: FontWeight.bold,
                  color: color,
                ),
              ),
              Gap(4.h),
              Text(
                title,
                style: TextStyle(
                  fontSize: 11.sp,
                  color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
                ),
                textAlign: TextAlign.center,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildLoadingStats() {
    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 4,
      childAspectRatio: 1.1,
      crossAxisSpacing: 16.w,
      mainAxisSpacing: 16.h,
      children: List.generate(8, (index) => Card(
        child: Center(
          child: CircularProgressIndicator(
            strokeWidth: 2,
          ),
        ),
      )),
    );
  }

  Widget _buildErrorStats(BuildContext context) {
    return Card(
      child: Padding(
        padding: EdgeInsets.all(32.w),
        child: Column(
          children: [
            Icon(
              Icons.error_outline,
              size: 48.sp,
              color: Colors.red,
            ),
            Gap(16.h),
            Text(
              'Failed to load statistics',
              style: TextStyle(
                fontSize: 16.sp,
                fontWeight: FontWeight.w500,
              ),
            ),
            Gap(8.h),
            ElevatedButton(
              onPressed: () {
                // Refresh data
              },
              child: const Text('Retry'),
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

  Widget _buildActivityItem(
    String title,
    String subtitle,
    IconData icon,
    Color color,
    DateTime time,
  ) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 8.h),
      child: Row(
        children: [
          CircleAvatar(
            radius: 20.r,
            backgroundColor: color.withOpacity(0.1),
            child: Icon(icon, color: color, size: 16.sp),
          ),
          Gap(12.w),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 14.sp,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Text(
                  subtitle,
                  style: TextStyle(
                    fontSize: 12.sp,
                    color: Colors.grey[600],
                  ),
                ),
              ],
            ),
          ),
          Text(
            _timeAgo(time),
            style: TextStyle(
              fontSize: 10.sp,
              color: Colors.grey[500],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPerformanceMetric(
    String label,
    String value,
    IconData icon,
    Color color,
  ) {
    return Column(
      children: [
        Icon(icon, color: color, size: 24.sp),
        Gap(8.h),
        Text(
          value,
          style: TextStyle(
            fontSize: 18.sp,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        Text(
          label,
          style: TextStyle(
            fontSize: 12.sp,
            color: Colors.grey[600],
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  String _timeAgo(DateTime time) {
    final now = DateTime.now();
    final difference = now.difference(time);
    
    if (difference.inDays > 0) {
      return '${difference.inDays}d ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours}h ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes}m ago';
    } else {
      return 'Just now';
    }
  }

  void _showNotifications(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Notifications'),
        content: SizedBox(
          width: 300,
          height: 300,
          child: ListView(
            children: [
              _buildNotificationItem(
                'Vehicle Sold',
                'Honda Civic has been sold to Ahmad Firdaus',
                Icons.sell,
                Colors.green,
                DateTime.now().subtract(const Duration(minutes: 30)),
              ),
              _buildNotificationItem(
                'Service Completed',
                'Oil change service completed for B 1234 ABC',
                Icons.build,
                Colors.blue,
                DateTime.now().subtract(const Duration(hours: 2)),
              ),
              _buildNotificationItem(
                'Low Stock Alert',
                'Engine oil stock is running low (5 units left)',
                Icons.warning,
                Colors.orange,
                DateTime.now().subtract(const Duration(hours: 4)),
              ),
              _buildNotificationItem(
                'New Customer',
                'John Doe has been added to customer database',
                Icons.person_add,
                Colors.purple,
                DateTime.now().subtract(const Duration(hours: 6)),
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Close'),
          ),
          ElevatedButton(
            onPressed: () {
              // Mark all as read
              Navigator.of(context).pop();
            },
            child: const Text('Mark All Read'),
          ),
        ],
      ),
    );
  }

  Widget _buildNotificationItem(
    String title,
    String subtitle,
    IconData icon,
    Color color,
    DateTime time,
  ) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: color.withOpacity(0.1),
        child: Icon(icon, color: color, size: 20),
      ),
      title: Text(
        title,
        style: const TextStyle(fontWeight: FontWeight.w500),
      ),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(subtitle),
          Gap(4),
          Text(
            _timeAgo(time),
            style: TextStyle(
              fontSize: 10.sp,
              color: Colors.grey[500],
            ),
          ),
        ],
      ),
      isThreeLine: true,
    );
  }
}