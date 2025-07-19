import 'package:flutter/material.dart';
import '../../../shared/widgets/app_drawer.dart';
import '../../../core/constants/app_constants.dart';

class ReportsPage extends StatefulWidget {
  const ReportsPage({super.key});

  @override
  State<ReportsPage> createState() => _ReportsPageState();
}

class _ReportsPageState extends State<ReportsPage> {
  String _selectedDateRange = 'This Month';
  String _selectedOutlet = 'All Outlets';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.reports),
        actions: [
          IconButton(
            icon: const Icon(Icons.date_range),
            onPressed: () {
              _showDateRangePicker();
            },
          ),
          IconButton(
            icon: const Icon(Icons.file_download),
            onPressed: () {
              _exportAllReports();
            },
          ),
        ],
      ),
      drawer: const AppDrawer(),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Filters
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<String>(
                    decoration: const InputDecoration(
                      labelText: 'Period',
                      border: OutlineInputBorder(),
                    ),
                    value: _selectedDateRange,
                    items: ['Today', 'This Week', 'This Month', 'This Quarter', 'This Year', 'Custom']
                        .map((String value) {
                      return DropdownMenuItem<String>(
                        value: value,
                        child: Text(value),
                      );
                    }).toList(),
                    onChanged: (String? newValue) {
                      setState(() {
                        _selectedDateRange = newValue!;
                      });
                    },
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: DropdownButtonFormField<String>(
                    decoration: const InputDecoration(
                      labelText: 'Outlet',
                      border: OutlineInputBorder(),
                    ),
                    value: _selectedOutlet,
                    items: ['All Outlets', 'Outlet 1', 'Outlet 2', 'Outlet 3']
                        .map((String value) {
                      return DropdownMenuItem<String>(
                        value: value,
                        child: Text(value),
                      );
                    }).toList(),
                    onChanged: (String? newValue) {
                      setState(() {
                        _selectedOutlet = newValue!;
                      });
                    },
                  ),
                ),
              ],
            ),
            const SizedBox(height: 24),

            // Quick Stats
            const Text(
              'Quick Overview',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),

            Row(
              children: [
                _buildOverviewCard('Total Revenue', 'Rp 850M', Icons.trending_up, Colors.green, '+12%'),
                const SizedBox(width: 16),
                _buildOverviewCard('Vehicle Sales', '87', Icons.directions_car, Colors.blue, '+8%'),
                const SizedBox(width: 16),
                _buildOverviewCard('Service Jobs', '245', Icons.build, Colors.orange, '+15%'),
                const SizedBox(width: 16),
                _buildOverviewCard('Net Profit', 'Rp 125M', Icons.account_balance, Colors.purple, '+18%'),
              ],
            ),
            const SizedBox(height: 32),

            // Report Categories
            const Text(
              'Report Categories',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),

            // Financial Reports
            _buildReportSection(
              'Financial Reports',
              Icons.account_balance,
              Colors.green,
              [
                _buildReportItem('Profit & Loss Statement', 'Comprehensive P&L analysis', () => _showProfitLoss()),
                _buildReportItem('Revenue Analysis', 'Revenue breakdown by service type', () => _showRevenueAnalysis()),
                _buildReportItem('Cash Flow Report', 'Track cash in and out', () => _showCashFlow()),
                _buildReportItem('Tax Reports', 'VAT and tax calculations', () => _showTaxReports()),
              ],
            ),

            const SizedBox(height: 24),

            // Vehicle Trading Reports
            _buildReportSection(
              'Vehicle Trading Reports',
              Icons.directions_car,
              Colors.blue,
              [
                _buildReportItem('Vehicle Sales Performance', 'Sales analytics and trends', () => _showVehicleSalesPerformance()),
                _buildReportItem('Inventory Aging Report', 'Vehicle stock aging analysis', () => _showInventoryAging()),
                _buildReportItem('Commission Report', 'Sales team commission tracking', () => _showCommissionReport()),
                _buildReportItem('Vehicle Profit Analysis', 'Margin analysis per vehicle', () => _showVehicleProfitAnalysis()),
              ],
            ),

            const SizedBox(height: 24),

            // Service Reports
            _buildReportSection(
              'Service Reports',
              Icons.build,
              Colors.orange,
              [
                _buildReportItem('Service Performance', 'Job completion and efficiency', () => _showServicePerformance()),
                _buildReportItem('Customer Satisfaction', 'Service quality metrics', () => _showCustomerSatisfaction()),
                _buildReportItem('Mechanic Performance', 'Individual mechanic statistics', () => _showMechanicPerformance()),
                _buildReportItem('Service Revenue', 'Revenue by service type', () => _showServiceRevenue()),
              ],
            ),

            const SizedBox(height: 24),

            // Inventory Reports
            _buildReportSection(
              'Inventory Reports',
              Icons.inventory,
              Colors.purple,
              [
                _buildReportItem('Stock Level Report', 'Current inventory status', () => _showStockLevel()),
                _buildReportItem('Stock Movement', 'In/out inventory tracking', () => _showStockMovement()),
                _buildReportItem('Inventory Valuation', 'Total stock value analysis', () => _showInventoryValuation()),
                _buildReportItem('Supplier Performance', 'Supplier analysis and ratings', () => _showSupplierPerformance()),
              ],
            ),

            const SizedBox(height: 24),

            // Customer Reports
            _buildReportSection(
              'Customer Reports',
              Icons.people,
              Colors.teal,
              [
                _buildReportItem('Customer Analysis', 'Customer behavior and trends', () => _showCustomerAnalysis()),
                _buildReportItem('Loyalty Program', 'Loyalty points and rewards', () => _showLoyaltyProgram()),
                _buildReportItem('Customer Lifetime Value', 'CLV analysis and segmentation', () => _showCustomerLTV()),
                _buildReportItem('Customer Retention', 'Retention rates and churn', () => _showCustomerRetention()),
              ],
            ),

            const SizedBox(height: 32),

            // Quick Charts
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Revenue Trend - Last 6 Months',
                      style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 16),
                    Container(
                      height: 200,
                      decoration: BoxDecoration(
                        border: Border.all(color: Colors.grey.shade300),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: const Center(
                        child: Text(
                          'Interactive Chart Area\nImplemented with fl_chart package',
                          textAlign: TextAlign.center,
                          style: TextStyle(color: Colors.grey),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOverviewCard(String title, String value, IconData icon, Color color, String change) {
    return Expanded(
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Icon(icon, color: color, size: 24),
                  const Spacer(),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: color.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      change,
                      style: TextStyle(
                        color: color,
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                value,
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                  color: color,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                title,
                style: const TextStyle(fontSize: 12, color: Colors.grey),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildReportSection(String title, IconData icon, Color color, List<Widget> reports) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(icon, color: color, size: 24),
                const SizedBox(width: 8),
                Text(
                  title,
                  style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                ),
              ],
            ),
            const SizedBox(height: 16),
            ...reports,
          ],
        ),
      ),
    );
  }

  Widget _buildReportItem(String title, String description, VoidCallback onTap) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 12.0),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(fontWeight: FontWeight.w500),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    description,
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey.shade600,
                    ),
                  ),
                ],
              ),
            ),
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextButton.icon(
                  onPressed: onTap,
                  icon: const Icon(Icons.visibility, size: 16),
                  label: const Text('View'),
                ),
                TextButton.icon(
                  onPressed: () => _exportReport(title),
                  icon: const Icon(Icons.file_download, size: 16),
                  label: const Text('Export'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  // Action methods
  void _showDateRangePicker() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Select Date Range'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: const Text('From Date'),
              trailing: const Text('01/12/2024'),
              onTap: () {
                // TODO: Show date picker
              },
            ),
            ListTile(
              title: const Text('To Date'),
              trailing: const Text('31/12/2024'),
              onTap: () {
                // TODO: Show date picker
              },
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Apply'),
          ),
        ],
      ),
    );
  }

  void _exportAllReports() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Exporting all reports...')),
    );
  }

  void _exportReport(String reportName) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Exporting $reportName...')),
    );
  }

  // Financial Reports
  void _showProfitLoss() {
    _navigateToReport('Profit & Loss Statement', _buildProfitLossReport());
  }

  void _showRevenueAnalysis() {
    _navigateToReport('Revenue Analysis', _buildRevenueAnalysisReport());
  }

  void _showCashFlow() {
    _navigateToReport('Cash Flow Report', _buildCashFlowReport());
  }

  void _showTaxReports() {
    _navigateToReport('Tax Reports', _buildTaxReport());
  }

  // Vehicle Trading Reports
  void _showVehicleSalesPerformance() {
    _navigateToReport('Vehicle Sales Performance', _buildVehicleSalesReport());
  }

  void _showInventoryAging() {
    _navigateToReport('Inventory Aging Report', _buildInventoryAgingReport());
  }

  void _showCommissionReport() {
    _navigateToReport('Commission Report', _buildCommissionReport());
  }

  void _showVehicleProfitAnalysis() {
    _navigateToReport('Vehicle Profit Analysis', _buildVehicleProfitReport());
  }

  // Other report methods...
  void _showServicePerformance() => _navigateToReport('Service Performance', _buildServiceReport());
  void _showCustomerSatisfaction() => _navigateToReport('Customer Satisfaction', _buildSatisfactionReport());
  void _showMechanicPerformance() => _navigateToReport('Mechanic Performance', _buildMechanicReport());
  void _showServiceRevenue() => _navigateToReport('Service Revenue', _buildServiceRevenueReport());
  void _showStockLevel() => _navigateToReport('Stock Level Report', _buildStockReport());
  void _showStockMovement() => _navigateToReport('Stock Movement', _buildStockMovementReport());
  void _showInventoryValuation() => _navigateToReport('Inventory Valuation', _buildInventoryValueReport());
  void _showSupplierPerformance() => _navigateToReport('Supplier Performance', _buildSupplierReport());
  void _showCustomerAnalysis() => _navigateToReport('Customer Analysis', _buildCustomerReport());
  void _showLoyaltyProgram() => _navigateToReport('Loyalty Program', _buildLoyaltyReport());
  void _showCustomerLTV() => _navigateToReport('Customer Lifetime Value', _buildLTVReport());
  void _showCustomerRetention() => _navigateToReport('Customer Retention', _buildRetentionReport());

  void _navigateToReport(String title, Widget content) {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (context) => Scaffold(
          appBar: AppBar(
            title: Text(title),
            actions: [
              IconButton(
                icon: const Icon(Icons.file_download),
                onPressed: () => _exportReport(title),
              ),
            ],
          ),
          body: content,
        ),
      ),
    );
  }

  Widget _buildProfitLossReport() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        children: [
          // Sample P&L data
          _buildReportTable(
            ['Category', 'Amount', 'Percentage'],
            [
              ['Revenue', 'Rp 850M', '100%'],
              ['Cost of Goods Sold', 'Rp 420M', '49.4%'],
              ['Gross Profit', 'Rp 430M', '50.6%'],
              ['Operating Expenses', 'Rp 305M', '35.9%'],
              ['Net Profit', 'Rp 125M', '14.7%'],
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildRevenueAnalysisReport() {
    return const Padding(
      padding: EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text('Revenue Analysis Report Content'),
          // TODO: Add actual revenue analysis content
        ],
      ),
    );
  }

  Widget _buildCashFlowReport() {
    return const Padding(
      padding: EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text('Cash Flow Report Content'),
          // TODO: Add actual cash flow content
        ],
      ),
    );
  }

  Widget _buildTaxReport() {
    return const Padding(
      padding: EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text('Tax Report Content'),
          // TODO: Add actual tax report content
        ],
      ),
    );
  }

  Widget _buildVehicleSalesReport() {
    return const Padding(
      padding: EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text('Vehicle Sales Performance Report'),
          // TODO: Add actual vehicle sales content
        ],
      ),
    );
  }

  Widget _buildInventoryAgingReport() {
    return const Padding(
      padding: EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text('Inventory Aging Report'),
          // TODO: Add actual aging report content
        ],
      ),
    );
  }

  Widget _buildCommissionReport() {
    return const Padding(
      padding: EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text('Commission Report'),
          // TODO: Add actual commission report content
        ],
      ),
    );
  }

  Widget _buildVehicleProfitReport() {
    return const Padding(
      padding: EdgeInsets.all(16.0),
      child: Column(
        children: [
          Text('Vehicle Profit Analysis'),
          // TODO: Add actual profit analysis content
        ],
      ),
    );
  }

  // Generic report widgets
  Widget _buildServiceReport() => const Text('Service Performance Report');
  Widget _buildSatisfactionReport() => const Text('Customer Satisfaction Report');
  Widget _buildMechanicReport() => const Text('Mechanic Performance Report');
  Widget _buildServiceRevenueReport() => const Text('Service Revenue Report');
  Widget _buildStockReport() => const Text('Stock Level Report');
  Widget _buildStockMovementReport() => const Text('Stock Movement Report');
  Widget _buildInventoryValueReport() => const Text('Inventory Valuation Report');
  Widget _buildSupplierReport() => const Text('Supplier Performance Report');
  Widget _buildCustomerReport() => const Text('Customer Analysis Report');
  Widget _buildLoyaltyReport() => const Text('Loyalty Program Report');
  Widget _buildLTVReport() => const Text('Customer Lifetime Value Report');
  Widget _buildRetentionReport() => const Text('Customer Retention Report');

  Widget _buildReportTable(List<String> headers, List<List<String>> rows) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            // Header
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.grey.shade100,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: headers.map((header) => Expanded(
                  child: Text(
                    header,
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                )).toList(),
              ),
            ),
            // Rows
            ...rows.map((row) => Container(
              padding: const EdgeInsets.all(12),
              child: Row(
                children: row.map((cell) => Expanded(
                  child: Text(cell),
                )).toList(),
              ),
            )),
          ],
        ),
      ),
    );
  }
}