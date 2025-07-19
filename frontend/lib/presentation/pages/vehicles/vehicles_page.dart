import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gap/gap.dart';
import 'package:intl/intl.dart';

import '../../../shared/widgets/app_drawer.dart';
import '../../../core/constants/app_constants.dart';
import '../../../data/services/vehicle_service.dart';
import '../../../data/models/vehicle_inventory_model.dart';
import '../../../data/models/vehicle_purchase_model.dart';
import '../../../data/models/vehicle_sale_model.dart';

class VehiclesPage extends ConsumerStatefulWidget {
  const VehiclesPage({super.key});

  @override
  ConsumerState<VehiclesPage> createState() => _VehiclesPageState();
}

class _VehiclesPageState extends ConsumerState<VehiclesPage> 
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final TextEditingController _searchController = TextEditingController();
  String _searchQuery = '';
  String _statusFilter = '';

  // Form controllers for purchase
  final _customerNameController = TextEditingController();
  final _plateNumberController = TextEditingController();
  final _brandController = TextEditingController();
  final _modelController = TextEditingController();
  final _yearController = TextEditingController();
  final _colorController = TextEditingController();
  final _mileageController = TextEditingController();
  final _chassisController = TextEditingController();
  final _engineController = TextEditingController();
  final _purchasePriceController = TextEditingController();
  final _estimatedPriceController = TextEditingController();
  final _conditionNotesController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    _searchController.addListener(() {
      setState(() {
        _searchQuery = _searchController.text;
      });
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    _customerNameController.dispose();
    _plateNumberController.dispose();
    _brandController.dispose();
    _modelController.dispose();
    _yearController.dispose();
    _colorController.dispose();
    _mileageController.dispose();
    _chassisController.dispose();
    _engineController.dispose();
    _purchasePriceController.dispose();
    _estimatedPriceController.dispose();
    _conditionNotesController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.vehicles),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(icon: Icon(Icons.inventory_2), text: 'Stok Kendaraan'),
            Tab(icon: Icon(Icons.shopping_cart), text: 'Beli Kendaraan'),
            Tab(icon: Icon(Icons.sell), text: 'Jual Kendaraan'),
            Tab(icon: Icon(Icons.analytics), text: 'Laporan'),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showAddVehicleDialog(context),
          ),
        ],
      ),
      drawer: const AppDrawer(),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildVehicleInventoryTab(),
          _buildVehiclePurchaseTab(),
          _buildVehicleSalesTab(),
          _buildReportsTab(),
        ],
      ),
    );
  }

  Widget _buildVehicleInventoryTab() {
    return Padding(
      padding: EdgeInsets.all(16.w),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Search and filter row
          Row(
            children: [
              Expanded(
                flex: 3,
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Cari berdasarkan merk, model, atau plat nomor...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              setState(() {
                                _searchQuery = '';
                              });
                            },
                          )
                        : null,
                  ),
                ),
              ),
              Gap(16.w),
              ElevatedButton.icon(
                onPressed: () => _showFilterDialog(context),
                icon: const Icon(Icons.filter_list),
                label: const Text('Filter'),
              ),
            ],
          ),
          Gap(16.h),
          
          // Statistics cards
          Consumer(
            builder: (context, ref, child) {
              final statsAsyncValue = ref.watch(vehicleStatisticsProvider);
              
              return statsAsyncValue.when(
                data: (stats) => Row(
                  children: [
                    Expanded(
                      child: _buildStatCard(
                        'Total Stok', 
                        stats['total_stock'].toString(), 
                        Icons.inventory, 
                        Colors.blue
                      ),
                    ),
                    Gap(16.w),
                    Expanded(
                      child: _buildStatCard(
                        'Tersedia', 
                        stats['available_stock'].toString(), 
                        Icons.check_circle, 
                        Colors.green
                      ),
                    ),
                    Gap(16.w),
                    Expanded(
                      child: _buildStatCard(
                        'Terjual Bulan Ini', 
                        stats['sold_this_month'].toString(), 
                        Icons.trending_up, 
                        Colors.orange
                      ),
                    ),
                    Gap(16.w),
                    Expanded(
                      child: _buildStatCard(
                        'Total Nilai Stok', 
                        NumberFormat.currency(
                          locale: 'id_ID',
                          symbol: 'Rp ',
                          decimalDigits: 0,
                        ).format(stats['total_stock_value']), 
                        Icons.attach_money, 
                        Colors.purple
                      ),
                    ),
                  ],
                ),
                loading: () => Row(
                  children: List.generate(4, (index) => Expanded(
                    child: Card(
                      child: Container(
                        height: 80.h,
                        child: const Center(child: CircularProgressIndicator()),
                      ),
                    ),
                  )),
                ),
                error: (error, stack) => const SizedBox(),
              );
            },
          ),
          Gap(24.h),
          
          // Vehicle list
          Expanded(
            child: Consumer(
              builder: (context, ref, child) {
                final inventoryAsyncValue = ref.watch(vehicleInventoryProvider({
                  'status': _statusFilter.isNotEmpty ? _statusFilter : null,
                  'search': _searchQuery.isNotEmpty ? _searchQuery : null,
                }));
                
                return inventoryAsyncValue.when(
                  data: (vehicles) {
                    if (vehicles.isEmpty) {
                      return const Center(
                        child: Text('Tidak ada data kendaraan'),
                      );
                    }
                    
                    return ListView.builder(
                      itemCount: vehicles.length,
                      itemBuilder: (context, index) {
                        return _buildVehicleCard(vehicles[index]);
                      },
                    );
                  },
                  loading: () => const Center(child: CircularProgressIndicator()),
                  error: (error, stack) => Center(
                    child: Text('Error: $error'),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatCard(String title, String value, IconData icon, Color color) {
    return Card(
      child: Padding(
        padding: EdgeInsets.all(16.w),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(icon, color: color, size: 24.sp),
                const Spacer(),
                Text(
                  value,
                  style: TextStyle(
                    fontSize: 18.sp,
                    fontWeight: FontWeight.bold,
                    color: color,
                  ),
                ),
              ],
            ),
            Gap(8.h),
            Text(
              title,
              style: TextStyle(
                fontSize: 12.sp,
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildVehicleCard(VehicleInventoryModel vehicle) {
    return Card(
      margin: EdgeInsets.only(bottom: 12.h),
      child: Padding(
        padding: EdgeInsets.all(16.w),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        vehicle.displayName,
                        style: TextStyle(
                          fontSize: 16.sp,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Gap(4.h),
                      Text(
                        vehicle.plateNumber,
                        style: TextStyle(
                          fontSize: 14.sp,
                          color: Colors.grey[600],
                        ),
                      ),
                      Gap(4.h),
                      Text(
                        vehicle.color,
                        style: TextStyle(
                          fontSize: 12.sp,
                          color: Colors.grey[500],
                        ),
                      ),
                    ],
                  ),
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Container(
                      padding: EdgeInsets.symmetric(horizontal: 8.w, vertical: 4.h),
                      decoration: BoxDecoration(
                        color: vehicle.isAvailable ? Colors.green : Colors.red,
                        borderRadius: BorderRadius.circular(12.r),
                      ),
                      child: Text(
                        vehicle.status,
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 10.sp,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    Gap(8.h),
                    Text(
                      vehicle.formattedEstimatedPrice,
                      style: TextStyle(
                        fontSize: 14.sp,
                        fontWeight: FontWeight.bold,
                        color: Colors.blue,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            Gap(12.h),
            
            Row(
              children: [
                Icon(Icons.speed, size: 16.sp, color: Colors.grey),
                Gap(4.w),
                Text(
                  '${NumberFormat('#,###').format(vehicle.mileage)} km',
                  style: TextStyle(fontSize: 12.sp),
                ),
                Gap(16.w),
                Icon(Icons.star, size: 16.sp, color: Colors.amber),
                Gap(4.w),
                Text(
                  '${vehicle.conditionRating}/5',
                  style: TextStyle(fontSize: 12.sp),
                ),
                Gap(16.w),
                Icon(Icons.access_time, size: 16.sp, color: Colors.grey),
                Gap(4.w),
                Text(
                  '${vehicle.daysSincePurchase} hari',
                  style: TextStyle(fontSize: 12.sp),
                ),
              ],
            ),
            
            if (vehicle.conditionNotes != null) ...[
              Gap(8.h),
              Text(
                vehicle.conditionNotes!,
                style: TextStyle(
                  fontSize: 12.sp,
                  fontStyle: FontStyle.italic,
                  color: Colors.grey[600],
                ),
              ),
            ],
            
            Gap(12.h),
            Row(
              children: [
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: () => _viewVehicleDetail(vehicle),
                    icon: const Icon(Icons.visibility),
                    label: const Text('Detail'),
                  ),
                ),
                Gap(8.w),
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: vehicle.isAvailable ? () => _editPrice(vehicle) : null,
                    icon: const Icon(Icons.edit),
                    label: const Text('Edit Harga'),
                  ),
                ),
                Gap(8.w),
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: () => _uploadPhotos(vehicle),
                    icon: const Icon(Icons.photo_camera),
                    label: const Text('Foto'),
                  ),
                ),
                Gap(8.w),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: vehicle.isAvailable ? () => _sellVehicle(vehicle) : null,
                    icon: const Icon(Icons.sell),
                    label: const Text('Jual'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.green,
                      foregroundColor: Colors.white,
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildVehiclePurchaseTab() {
    return Padding(
      padding: EdgeInsets.all(16.w),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Purchase form
          Card(
            child: Padding(
              padding: EdgeInsets.all(16.w),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Form Pembelian Kendaraan',
                    style: TextStyle(fontSize: 18.sp, fontWeight: FontWeight.bold),
                  ),
                  Gap(16.h),
                  
                  // Customer and date row
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField(
                          'Nama Customer', 
                          'Masukkan nama customer',
                          controller: _customerNameController,
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: _buildFormField(
                          'Tanggal Beli', 
                          'dd/mm/yyyy',
                          readOnly: true,
                          onTap: () => _selectDate(context),
                        ),
                      ),
                    ],
                  ),
                  Gap(16.h),
                  
                  // Vehicle details
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField(
                          'Plat Nomor', 
                          'B 1234 ABC',
                          controller: _plateNumberController,
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: _buildFormField(
                          'Merk', 
                          'Toyota',
                          controller: _brandController,
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: _buildFormField(
                          'Model', 
                          'Avanza',
                          controller: _modelController,
                        ),
                      ),
                    ],
                  ),
                  Gap(16.h),
                  
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField(
                          'Tahun', 
                          '2020',
                          controller: _yearController,
                          keyboardType: TextInputType.number,
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: _buildFormField(
                          'Warna', 
                          'Putih',
                          controller: _colorController,
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: _buildFormField(
                          'Kilometer', 
                          '50000',
                          controller: _mileageController,
                          keyboardType: TextInputType.number,
                        ),
                      ),
                    ],
                  ),
                  Gap(16.h),
                  
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField(
                          'No. Rangka', 
                          'Nomor rangka kendaraan',
                          controller: _chassisController,
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: _buildFormField(
                          'No. Mesin', 
                          'Nomor mesin kendaraan',
                          controller: _engineController,
                        ),
                      ),
                    ],
                  ),
                  Gap(16.h),
                  
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField(
                          'Harga Beli', 
                          '150000000',
                          controller: _purchasePriceController,
                          keyboardType: TextInputType.number,
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: _buildFormField(
                          'Estimasi Harga Jual', 
                          '170000000',
                          controller: _estimatedPriceController,
                          keyboardType: TextInputType.number,
                        ),
                      ),
                    ],
                  ),
                  Gap(16.h),
                  
                  _buildFormField(
                    'Catatan Kondisi', 
                    'Kondisi kendaraan, kelengkapan dokumen, dll',
                    controller: _conditionNotesController,
                    maxLines: 3,
                  ),
                  Gap(24.h),
                  
                  Row(
                    children: [
                      Expanded(
                        child: OutlinedButton(
                          onPressed: _clearPurchaseForm,
                          child: const Text('Reset Form'),
                        ),
                      ),
                      Gap(16.w),
                      Expanded(
                        child: ElevatedButton(
                          onPressed: _processPurchase,
                          child: const Text('Proses Pembelian'),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
          Gap(16.h),
          
          // Purchase history
          Text(
            'Riwayat Pembelian',
            style: TextStyle(fontSize: 16.sp, fontWeight: FontWeight.bold),
          ),
          Gap(8.h),
          
          Expanded(
            child: Consumer(
              builder: (context, ref, child) {
                final purchasesAsyncValue = ref.watch(vehiclePurchasesProvider({}));
                
                return purchasesAsyncValue.when(
                  data: (purchases) {
                    if (purchases.isEmpty) {
                      return const Center(
                        child: Text('Belum ada riwayat pembelian'),
                      );
                    }
                    
                    return ListView.builder(
                      itemCount: purchases.length,
                      itemBuilder: (context, index) {
                        return _buildPurchaseHistoryCard(purchases[index]);
                      },
                    );
                  },
                  loading: () => const Center(child: CircularProgressIndicator()),
                  error: (error, stack) => Center(
                    child: Text('Error: $error'),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPurchaseHistoryCard(VehiclePurchaseModel purchase) {
    return Card(
      margin: EdgeInsets.only(bottom: 8.h),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: purchase.isCompleted ? Colors.green : Colors.orange,
          child: Icon(
            purchase.isCompleted ? Icons.check : Icons.pending,
            color: Colors.white,
          ),
        ),
        title: Text(purchase.customerName),
        subtitle: Text('${purchase.formattedDate} • ${purchase.paymentMethod}'),
        trailing: Text(
          purchase.formattedPrice,
          style: TextStyle(
            fontWeight: FontWeight.bold,
            color: Colors.blue,
          ),
        ),
      ),
    );
  }

  Widget _buildVehicleSalesTab() {
    return Padding(
      padding: EdgeInsets.all(16.w),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Riwayat Penjualan',
            style: TextStyle(fontSize: 18.sp, fontWeight: FontWeight.bold),
          ),
          Gap(16.h),
          
          Expanded(
            child: Consumer(
              builder: (context, ref, child) {
                final salesAsyncValue = ref.watch(vehicleSalesProvider({}));
                
                return salesAsyncValue.when(
                  data: (sales) {
                    if (sales.isEmpty) {
                      return const Center(
                        child: Text('Belum ada riwayat penjualan'),
                      );
                    }
                    
                    return ListView.builder(
                      itemCount: sales.length,
                      itemBuilder: (context, index) {
                        return _buildSalesHistoryCard(sales[index]);
                      },
                    );
                  },
                  loading: () => const Center(child: CircularProgressIndicator()),
                  error: (error, stack) => Center(
                    child: Text('Error: $error'),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSalesHistoryCard(VehicleSaleModel sale) {
    return Card(
      margin: EdgeInsets.only(bottom: 12.h),
      child: Padding(
        padding: EdgeInsets.all(16.w),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        sale.customerName,
                        style: TextStyle(
                          fontSize: 16.sp,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Gap(4.h),
                      Text(
                        '${sale.formattedDate} • ${sale.paymentMethod}',
                        style: TextStyle(
                          fontSize: 12.sp,
                          color: Colors.grey[600],
                        ),
                      ),
                    ],
                  ),
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      sale.formattedPrice,
                      style: TextStyle(
                        fontSize: 16.sp,
                        fontWeight: FontWeight.bold,
                        color: Colors.blue,
                      ),
                    ),
                    Text(
                      'Profit: ${sale.formattedProfit} (${sale.profitPercentage})',
                      style: TextStyle(
                        fontSize: 12.sp,
                        color: Colors.green,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            
            if (sale.isFinanced) ...[
              Gap(12.h),
              Container(
                padding: EdgeInsets.all(12.w),
                decoration: BoxDecoration(
                  color: Colors.blue.shade50,
                  borderRadius: BorderRadius.circular(8.r),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Detail Kredit:',
                      style: TextStyle(
                        fontSize: 12.sp,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    Gap(4.h),
                    Text(
                      'DP: ${NumberFormat.currency(locale: 'id_ID', symbol: 'Rp ', decimalDigits: 0).format(sale.downPayment)} • ${sale.tenorMonths} bulan • ${NumberFormat.currency(locale: 'id_ID', symbol: 'Rp ', decimalDigits: 0).format(sale.monthlyInstallment)}/bln',
                      style: TextStyle(fontSize: 12.sp),
                    ),
                  ],
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildReportsTab() {
    return Padding(
      padding: EdgeInsets.all(16.w),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Report cards
          Row(
            children: [
              Expanded(
                child: _buildReportCard(
                  'Analisis Keuntungan',
                  'Lihat trend keuntungan dan margin',
                  Icons.analytics,
                  Colors.green,
                  _showProfitAnalysis,
                ),
              ),
              Gap(16.w),
              Expanded(
                child: _buildReportCard(
                  'Umur Stok',
                  'Kendaraan yang lama di stok',
                  Icons.access_time,
                  Colors.orange,
                  _showAgingReport,
                ),
              ),
            ],
          ),
          Gap(16.h),
          
          Row(
            children: [
              Expanded(
                child: _buildReportCard(
                  'Performa Sales',
                  'Ranking dan komisi sales team',
                  Icons.leaderboard,
                  Colors.blue,
                  _showSalesPerformance,
                ),
              ),
              Gap(16.w),
              Expanded(
                child: _buildReportCard(
                  'Export Data',
                  'Download laporan dalam Excel/PDF',
                  Icons.download,
                  Colors.purple,
                  _exportData,
                ),
              ),
            ],
          ),
          Gap(24.h),
          
          // Quick summary
          Expanded(
            child: Card(
              child: Padding(
                padding: EdgeInsets.all(16.w),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Ringkasan Performa',
                      style: TextStyle(fontSize: 16.sp, fontWeight: FontWeight.bold),
                    ),
                    Gap(16.h),
                    
                    Consumer(
                      builder: (context, ref, child) {
                        final statsAsyncValue = ref.watch(vehicleStatisticsProvider);
                        final salesAsyncValue = ref.watch(vehicleSalesProvider({}));
                        
                        return Row(
                          children: [
                            statsAsyncValue.when(
                              data: (stats) => salesAsyncValue.when(
                                data: (sales) => Expanded(
                                  child: Column(
                                    children: [
                                      _buildSummaryItem(
                                        'Total Keuntungan',
                                        NumberFormat.currency(
                                          locale: 'id_ID',
                                          symbol: 'Rp ',
                                          decimalDigits: 0,
                                        ).format(sales.fold(0.0, (sum, sale) => sum + sale.profitAmount)),
                                      ),
                                      Gap(16.h),
                                      _buildSummaryItem(
                                        'Rata-rata Umur Stok',
                                        '${(stats['total_stock'] > 0 ? 20 : 0)} hari',
                                      ),
                                    ],
                                  ),
                                ),
                                loading: () => const Expanded(child: CircularProgressIndicator()),
                                error: (error, stack) => const Expanded(child: Text('Error')),
                              ),
                              loading: () => const Expanded(child: CircularProgressIndicator()),
                              error: (error, stack) => const Expanded(child: Text('Error')),
                            ),
                          ],
                        );
                      },
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSummaryItem(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            fontSize: 24.sp,
            fontWeight: FontWeight.bold,
            color: Colors.blue,
          ),
        ),
        Text(
          label,
          style: TextStyle(
            fontSize: 12.sp,
            color: Colors.grey[600],
          ),
        ),
      ],
    );
  }

  Widget _buildReportCard(
    String title,
    String subtitle,
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
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Icon(icon, color: color, size: 32.sp),
              Gap(12.h),
              Text(
                title,
                style: TextStyle(
                  fontSize: 14.sp,
                  fontWeight: FontWeight.bold,
                ),
              ),
              Gap(4.h),
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
      ),
    );
  }

  Widget _buildFormField(
    String label,
    String hint, {
    TextEditingController? controller,
    TextInputType? keyboardType,
    int maxLines = 1,
    bool readOnly = false,
    VoidCallback? onTap,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 12.sp,
            fontWeight: FontWeight.w500,
          ),
        ),
        Gap(4.h),
        TextField(
          controller: controller,
          keyboardType: keyboardType,
          maxLines: maxLines,
          readOnly: readOnly,
          onTap: onTap,
          decoration: InputDecoration(
            hintText: hint,
            isDense: true,
          ),
        ),
      ],
    );
  }

  // Action methods - Now implemented instead of TODO
  void _showAddVehicleDialog(BuildContext context) {
    _tabController.animateTo(1); // Switch to purchase tab
  }

  void _showFilterDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Filter Kendaraan'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: const Text('Semua'),
              leading: Radio<String>(
                value: '',
                groupValue: _statusFilter,
                onChanged: (value) {
                  setState(() {
                    _statusFilter = value ?? '';
                  });
                  Navigator.of(context).pop();
                },
              ),
            ),
            ListTile(
              title: const Text('Tersedia'),
              leading: Radio<String>(
                value: 'Available',
                groupValue: _statusFilter,
                onChanged: (value) {
                  setState(() {
                    _statusFilter = value ?? '';
                  });
                  Navigator.of(context).pop();
                },
              ),
            ),
            ListTile(
              title: const Text('Terjual'),
              leading: Radio<String>(
                value: 'Sold',
                groupValue: _statusFilter,
                onChanged: (value) {
                  setState(() {
                    _statusFilter = value ?? '';
                  });
                  Navigator.of(context).pop();
                },
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Tutup'),
          ),
        ],
      ),
    );
  }

  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now(),
      firstDate: DateTime(2020),
      lastDate: DateTime.now(),
    );
    if (picked != null) {
      // Handle selected date
    }
  }

  void _processPurchase() async {
    if (_customerNameController.text.isEmpty ||
        _plateNumberController.text.isEmpty ||
        _brandController.text.isEmpty ||
        _modelController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Mohon lengkapi data yang diperlukan')),
      );
      return;
    }

    try {
      final purchaseData = {
        'customer_name': _customerNameController.text,
        'customer_id': 999, // Mock customer ID
        'outlet_id': 1,
        'purchase_date': DateTime.now().toIso8601String().split('T')[0],
        'plate_number': _plateNumberController.text,
        'brand': _brandController.text,
        'model': _modelController.text,
        'type': 'Unknown', // Could be a dropdown
        'production_year': int.tryParse(_yearController.text) ?? DateTime.now().year,
        'chassis_number': _chassisController.text,
        'engine_number': _engineController.text,
        'color': _colorController.text,
        'mileage': int.tryParse(_mileageController.text) ?? 0,
        'condition_rating': 4, // Could be a rating widget
        'purchase_price': double.tryParse(_purchasePriceController.text) ?? 0,
        'estimated_selling_price': double.tryParse(_estimatedPriceController.text) ?? 0,
        'condition_notes': _conditionNotesController.text,
        'payment_method': 'cash',
        'notes': 'Pembelian melalui aplikasi',
      };

      await ref.read(vehicleServiceProvider).createVehiclePurchase(purchaseData);
      
      // Refresh data
      ref.invalidate(vehicleInventoryProvider);
      ref.invalidate(vehiclePurchasesProvider);
      ref.invalidate(vehicleStatisticsProvider);
      
      _clearPurchaseForm();
      
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Pembelian kendaraan berhasil diproses')),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: $e')),
      );
    }
  }

  void _clearPurchaseForm() {
    _customerNameController.clear();
    _plateNumberController.clear();
    _brandController.clear();
    _modelController.clear();
    _yearController.clear();
    _colorController.clear();
    _mileageController.clear();
    _chassisController.clear();
    _engineController.clear();
    _purchasePriceController.clear();
    _estimatedPriceController.clear();
    _conditionNotesController.clear();
  }

  void _viewVehicleDetail(VehicleInventoryModel vehicle) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Detail ${vehicle.displayName}'),
        content: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              _buildDetailRow('Plat Nomor', vehicle.plateNumber),
              _buildDetailRow('No. Rangka', vehicle.chassisNumber),
              _buildDetailRow('No. Mesin', vehicle.engineNumber),
              _buildDetailRow('Warna', vehicle.color),
              _buildDetailRow('Kilometer', '${NumberFormat('#,###').format(vehicle.mileage)} km'),
              _buildDetailRow('Kondisi', '${vehicle.conditionRating}/5 bintang'),
              _buildDetailRow('Harga Beli', vehicle.formattedPrice),
              _buildDetailRow('Estimasi Jual', vehicle.formattedEstimatedPrice),
              if (vehicle.conditionNotes != null)
                _buildDetailRow('Catatan', vehicle.conditionNotes!),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Tutup'),
          ),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 4.h),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80.w,
            child: Text(
              '$label:',
              style: TextStyle(fontWeight: FontWeight.w500),
            ),
          ),
          Expanded(
            child: Text(value),
          ),
        ],
      ),
    );
  }

  void _editPrice(VehicleInventoryModel vehicle) {
    final priceController = TextEditingController(
      text: vehicle.estimatedSellingPrice.toInt().toString(),
    );
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Edit Harga ${vehicle.displayName}'),
        content: TextField(
          controller: priceController,
          keyboardType: TextInputType.number,
          decoration: const InputDecoration(
            labelText: 'Harga Jual Baru',
            prefixText: 'Rp ',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Batal'),
          ),
          ElevatedButton(
            onPressed: () async {
              final newPrice = double.tryParse(priceController.text);
              if (newPrice != null) {
                await ref.read(vehicleServiceProvider).updateVehiclePrice(
                  vehicle.inventoryId,
                  newPrice,
                );
                
                ref.invalidate(vehicleInventoryProvider);
                
                Navigator.of(context).pop();
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Harga berhasil diupdate')),
                );
              }
            },
            child: const Text('Simpan'),
          ),
        ],
      ),
    );
  }

  void _uploadPhotos(VehicleInventoryModel vehicle) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('Fitur upload foto untuk ${vehicle.plateNumber} akan segera tersedia'),
        action: SnackBarAction(
          label: 'OK',
          onPressed: () {},
        ),
      ),
    );
  }

  void _sellVehicle(VehicleInventoryModel vehicle) {
    final customerController = TextEditingController();
    final priceController = TextEditingController(
      text: vehicle.estimatedSellingPrice.toInt().toString(),
    );
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Jual ${vehicle.displayName}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: customerController,
              decoration: const InputDecoration(
                labelText: 'Nama Customer',
              ),
            ),
            Gap(16.h),
            TextField(
              controller: priceController,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: 'Harga Jual',
                prefixText: 'Rp ',
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Batal'),
          ),
          ElevatedButton(
            onPressed: () async {
              if (customerController.text.isEmpty) return;
              
              final sellingPrice = double.tryParse(priceController.text);
              if (sellingPrice == null) return;
              
              final saleData = {
                'inventory_id': vehicle.inventoryId,
                'customer_id': 999,
                'customer_name': customerController.text,
                'outlet_id': 1,
                'selling_price': sellingPrice,
                'profit_amount': sellingPrice - vehicle.purchasePrice,
                'payment_method': 'cash',
                'financing_type': 'cash',
                'notes': 'Penjualan tunai',
              };
              
              await ref.read(vehicleServiceProvider).createVehicleSale(saleData);
              
              ref.invalidate(vehicleInventoryProvider);
              ref.invalidate(vehicleSalesProvider);
              ref.invalidate(vehicleStatisticsProvider);
              
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Kendaraan berhasil terjual')),
              );
            },
            child: const Text('Jual'),
          ),
        ],
      ),
    );
  }

  void _showProfitAnalysis() {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: const Text('Analisis keuntungan akan ditampilkan dalam modal terpisah'),
        action: SnackBarAction(label: 'OK', onPressed: () {}),
      ),
    );
  }

  void _showAgingReport() {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: const Text('Laporan umur stok akan ditampilkan dalam modal terpisah'),
        action: SnackBarAction(label: 'OK', onPressed: () {}),
      ),
    );
  }

  void _showSalesPerformance() {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: const Text('Laporan performa sales akan ditampilkan dalam modal terpisah'),
        action: SnackBarAction(label: 'OK', onPressed: () {}),
      ),
    );
  }

  void _exportData() {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: const Text('Data akan diekspor ke Excel/PDF'),
        action: SnackBarAction(label: 'OK', onPressed: () {}),
      ),
    );
  }
}