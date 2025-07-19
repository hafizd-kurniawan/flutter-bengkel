import 'package:flutter/material.dart';
import '../../../shared/widgets/app_drawer.dart';
import '../../../core/constants/app_constants.dart';

class VehiclesPage extends StatefulWidget {
  const VehiclesPage({super.key});

  @override
  State<VehiclesPage> createState() => _VehiclesPageState();
}

class _VehiclesPageState extends State<VehiclesPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
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
            onPressed: () {
              _showAddVehicleDialog(context);
            },
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
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Search and filter row
          Row(
            children: [
              Expanded(
                flex: 3,
                child: TextField(
                  decoration: InputDecoration(
                    hintText: 'Cari berdasarkan merk, model, atau plat nomor...',
                    prefixIcon: const Icon(Icons.search),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 16),
              ElevatedButton.icon(
                onPressed: () {
                  _showFilterDialog(context);
                },
                icon: const Icon(Icons.filter_list),
                label: const Text('Filter'),
              ),
            ],
          ),
          const SizedBox(height: 16),
          
          // Statistics cards
          Row(
            children: [
              _buildStatCard('Total Stok', '125', Icons.inventory, Colors.blue),
              const SizedBox(width: 16),
              _buildStatCard('Tersedia', '98', Icons.check_circle, Colors.green),
              const SizedBox(width: 16),
              _buildStatCard('Terjual Bulan Ini', '27', Icons.trending_up, Colors.orange),
              const SizedBox(width: 16),
              _buildStatCard('Total Nilai Stok', 'Rp 2.5M', Icons.attach_money, Colors.purple),
            ],
          ),
          const SizedBox(height: 24),
          
          // Vehicle list
          Expanded(
            child: ListView.builder(
              itemCount: 10, // Mock data
              itemBuilder: (context, index) {
                return _buildVehicleCard();
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildVehiclePurchaseTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Purchase form
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Form Pembelian Kendaraan',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 16),
                  
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField('Customer', 'Pilih customer'),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildFormField('Tanggal Beli', 'dd/mm/yyyy'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField('Plat Nomor', 'B 1234 XYZ'),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildFormField('Merk', 'Toyota'),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildFormField('Model', 'Avanza'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField('Tahun', '2020'),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildFormField('Warna', 'Hitam'),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildFormField('Harga Beli', 'Rp 150.000.000'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  
                  Row(
                    children: [
                      Expanded(
                        child: _buildFormField('Kondisi (1-5)', '4'),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildFormField('Estimasi Harga Jual', 'Rp 165.000.000'),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Container(), // Spacer
                      ),
                    ],
                  ),
                  const SizedBox(height: 24),
                  
                  Row(
                    children: [
                      ElevatedButton.icon(
                        onPressed: () {
                          _processPurchase();
                        },
                        icon: const Icon(Icons.save),
                        label: const Text('Simpan Pembelian'),
                        style: ElevatedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                        ),
                      ),
                      const SizedBox(width: 16),
                      OutlinedButton.icon(
                        onPressed: () {
                          _clearForm();
                        },
                        icon: const Icon(Icons.clear),
                        label: const Text('Clear Form'),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          
          // Recent purchases
          const Text(
            'Pembelian Terbaru',
            style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 8),
          Expanded(
            child: ListView.builder(
              itemCount: 5,
              itemBuilder: (context, index) {
                return _buildPurchaseHistoryCard();
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildVehicleSalesTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Sales statistics
          Row(
            children: [
              _buildStatCard('Penjualan Hari Ini', '3', Icons.today, Colors.green),
              const SizedBox(width: 16),
              _buildStatCard('Penjualan Bulan Ini', '27', Icons.calendar_month, Colors.blue),
              const SizedBox(width: 16),
              _buildStatCard('Total Keuntungan', 'Rp 450M', Icons.trending_up, Colors.orange),
              const SizedBox(width: 16),
              _buildStatCard('Komisi Sales', 'Rp 12M', Icons.person, Colors.purple),
            ],
          ),
          const SizedBox(height: 24),
          
          // Quick sale section
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                children: [
                  Expanded(
                    flex: 2,
                    child: _buildFormField('Pilih Kendaraan', 'Cari kendaraan tersedia'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Customer', 'Pilih customer'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Harga Jual', 'Rp 165.000.000'),
                  ),
                  const SizedBox(width: 16),
                  ElevatedButton.icon(
                    onPressed: () {
                      _processSale();
                    },
                    icon: const Icon(Icons.sell),
                    label: const Text('Proses Jual'),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          
          // Sales history
          const Text(
            'Riwayat Penjualan',
            style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 8),
          Expanded(
            child: ListView.builder(
              itemCount: 10,
              itemBuilder: (context, index) {
                return _buildSalesHistoryCard();
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildReportsTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
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
                  () => _showProfitAnalysis(),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: _buildReportCard(
                  'Umur Stok',
                  'Kendaraan yang lama di stok',
                  Icons.access_time,
                  Colors.orange,
                  () => _showAgingReport(),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          
          Row(
            children: [
              Expanded(
                child: _buildReportCard(
                  'Performa Sales',
                  'Ranking dan komisi sales team',
                  Icons.leaderboard,
                  Colors.blue,
                  () => _showSalesPerformance(),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: _buildReportCard(
                  'Export Data',
                  'Download laporan dalam Excel/PDF',
                  Icons.download,
                  Colors.purple,
                  () => _exportData(),
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          
          // Quick charts
          Expanded(
            child: Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Grafik Penjualan 6 Bulan Terakhir',
                      style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 16),
                    Expanded(
                      child: Container(
                        decoration: BoxDecoration(
                          border: Border.all(color: Colors.grey.shade300),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: const Center(
                          child: Text(
                            'Chart akan ditampilkan di sini\nDengan fl_chart package',
                            textAlign: TextAlign.center,
                            style: TextStyle(color: Colors.grey),
                          ),
                        ),
                      ),
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

  Widget _buildStatCard(String title, String value, IconData icon, Color color) {
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
                  Text(
                    value,
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: color,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
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

  Widget _buildVehicleCard() {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          children: [
            // Vehicle image placeholder
            Container(
              width: 80,
              height: 60,
              decoration: BoxDecoration(
                color: Colors.grey.shade300,
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Icon(Icons.directions_car, size: 32),
            ),
            const SizedBox(width: 16),
            
            // Vehicle info
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Toyota Avanza 2020',
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                  Text('B 1234 XYZ • Hitam • 45.000 km'),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Icon(Icons.star, color: Colors.orange, size: 16),
                      Text(' 4.5 • Rp 165.000.000'),
                      const Spacer(),
                      Chip(
                        label: Text('Tersedia'),
                        backgroundColor: Colors.green.shade100,
                        labelStyle: TextStyle(color: Colors.green.shade700),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            
            // Actions
            PopupMenuButton(
              itemBuilder: (context) => [
                PopupMenuItem(
                  child: Text('Lihat Detail'),
                  onTap: () => _viewVehicleDetail(),
                ),
                PopupMenuItem(
                  child: Text('Edit Harga'),
                  onTap: () => _editPrice(),
                ),
                PopupMenuItem(
                  child: Text('Upload Foto'),
                  onTap: () => _uploadPhotos(),
                ),
                PopupMenuItem(
                  child: Text('Jual Kendaraan'),
                  onTap: () => _sellVehicle(),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFormField(String label, String hint) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(fontWeight: FontWeight.w500)),
        const SizedBox(height: 4),
        TextField(
          decoration: InputDecoration(
            hintText: hint,
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(8)),
            contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          ),
        ),
      ],
    );
  }

  Widget _buildPurchaseHistoryCard() {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Colors.blue.shade100,
          child: const Icon(Icons.shopping_cart, color: Colors.blue),
        ),
        title: const Text('Honda Civic 2019 - B 5678 ABC'),
        subtitle: const Text('15 Des 2024 • Rp 200.000.000 • dari Budi Santoso'),
        trailing: const Text('Berhasil', style: TextStyle(color: Colors.green)),
      ),
    );
  }

  Widget _buildSalesHistoryCard() {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Colors.green.shade100,
          child: const Icon(Icons.sell, color: Colors.green),
        ),
        title: const Text('Toyota Avanza 2020 - B 1234 XYZ'),
        subtitle: const Text('14 Des 2024 • Rp 165.000.000 • ke Siti Aminah'),
        trailing: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Keuntungan', style: TextStyle(fontSize: 10)),
            Text(
              'Rp 15M',
              style: TextStyle(color: Colors.green, fontWeight: FontWeight.bold),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildReportCard(String title, String subtitle, IconData icon, Color color, VoidCallback onTap) {
    return Card(
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            children: [
              Icon(icon, size: 48, color: color),
              const SizedBox(height: 8),
              Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 4),
              Text(subtitle, textAlign: TextAlign.center, style: const TextStyle(fontSize: 12)),
            ],
          ),
        ),
      ),
    );
  }

  // Action methods
  void _showAddVehicleDialog(BuildContext context) {
    // TODO: Implement add vehicle dialog
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Tambah Kendaraan'),
        content: const Text('Form tambah kendaraan akan ditampilkan di sini'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Tutup'),
          ),
        ],
      ),
    );
  }

  void _showFilterDialog(BuildContext context) {
    // TODO: Implement filter dialog
  }

  void _processPurchase() {
    // TODO: Implement purchase processing
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Pembelian kendaraan berhasil diproses')),
    );
  }

  void _clearForm() {
    // TODO: Clear form fields
  }

  void _processSale() {
    // TODO: Implement sale processing
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Penjualan kendaraan berhasil diproses')),
    );
  }

  void _viewVehicleDetail() {
    // TODO: Navigate to vehicle detail page
  }

  void _editPrice() {
    // TODO: Show edit price dialog
  }

  void _uploadPhotos() {
    // TODO: Show photo upload dialog
  }

  void _sellVehicle() {
    // TODO: Show sell vehicle dialog
  }

  void _showProfitAnalysis() {
    // TODO: Navigate to profit analysis page
  }

  void _showAgingReport() {
    // TODO: Navigate to aging report page
  }

  void _showSalesPerformance() {
    // TODO: Navigate to sales performance page
  }

  void _exportData() {
    // TODO: Implement data export
  }
}