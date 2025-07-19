import 'package:flutter/material.dart';
import '../../../shared/widgets/app_drawer.dart';
import '../../../core/constants/app_constants.dart';

class InventoryPage extends StatefulWidget {
  const InventoryPage({super.key});

  @override
  State<InventoryPage> createState() => _InventoryPageState();
}

class _InventoryPageState extends State<InventoryPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final TextEditingController _searchController = TextEditingController();
  String _selectedCategory = 'All';
  String _sortBy = 'Name';

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(AppStrings.inventory),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(icon: Icon(Icons.inventory), text: 'Stok Barang'),
            Tab(icon: Icon(Icons.warning), text: 'Stok Menipis'),
            Tab(icon: Icon(Icons.add_box), text: 'Tambah Barang'),
            Tab(icon: Icon(Icons.analytics), text: 'Laporan'),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              _tabController.animateTo(2);
            },
          ),
          IconButton(
            icon: const Icon(Icons.file_download),
            onPressed: () {
              _exportInventoryData();
            },
          ),
        ],
      ),
      drawer: const AppDrawer(),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildInventoryTab(),
          _buildLowStockTab(),
          _buildAddProductTab(),
          _buildReportsTab(),
        ],
      ),
    );
  }

  Widget _buildInventoryTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        children: [
          // Search and filters
          Row(
            children: [
              Expanded(
                flex: 3,
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Cari produk berdasarkan nama atau kode...',
                    prefixIcon: const Icon(Icons.search),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  onChanged: (value) {
                    // TODO: Implement search
                  },
                ),
              ),
              const SizedBox(width: 16),
              DropdownButton<String>(
                value: _selectedCategory,
                onChanged: (String? newValue) {
                  setState(() {
                    _selectedCategory = newValue!;
                  });
                },
                items: ['All', 'Sparepart', 'Oil', 'Tire', 'Battery', 'Filter']
                    .map<DropdownMenuItem<String>>((String value) {
                  return DropdownMenuItem<String>(
                    value: value,
                    child: Text(value),
                  );
                }).toList(),
              ),
              const SizedBox(width: 16),
              DropdownButton<String>(
                value: _sortBy,
                onChanged: (String? newValue) {
                  setState(() {
                    _sortBy = newValue!;
                  });
                },
                items: ['Name', 'Stock', 'Price', 'Category']
                    .map<DropdownMenuItem<String>>((String value) {
                  return DropdownMenuItem<String>(
                    value: value,
                    child: Text('Sort by $value'),
                  );
                }).toList(),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // Statistics cards
          Row(
            children: [
              _buildStatCard('Total Items', '2,847', Icons.inventory, Colors.blue),
              const SizedBox(width: 16),
              _buildStatCard('Total Value', 'Rp 450M', Icons.attach_money, Colors.green),
              const SizedBox(width: 16),
              _buildStatCard('Low Stock', '23', Icons.warning, Colors.orange),
              const SizedBox(width: 16),
              _buildStatCard('Out of Stock', '5', Icons.error, Colors.red),
            ],
          ),
          const SizedBox(height: 24),

          // Product list header
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.grey.shade100,
              borderRadius: BorderRadius.circular(8),
            ),
            child: const Row(
              children: [
                Expanded(flex: 2, child: Text('Produk', style: TextStyle(fontWeight: FontWeight.bold))),
                Expanded(child: Text('Kategori', style: TextStyle(fontWeight: FontWeight.bold))),
                Expanded(child: Text('Stok', style: TextStyle(fontWeight: FontWeight.bold))),
                Expanded(child: Text('Harga Beli', style: TextStyle(fontWeight: FontWeight.bold))),
                Expanded(child: Text('Harga Jual', style: TextStyle(fontWeight: FontWeight.bold))),
                SizedBox(width: 100, child: Text('Aksi', style: TextStyle(fontWeight: FontWeight.bold))),
              ],
            ),
          ),
          const SizedBox(height: 8),

          // Product list
          Expanded(
            child: ListView.builder(
              itemCount: 20,
              itemBuilder: (context, index) {
                return _buildProductRow(index);
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLowStockTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Card(
            color: Colors.orange.shade50,
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                children: [
                  Icon(Icons.warning, color: Colors.orange, size: 32),
                  const SizedBox(width: 16),
                  const Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Peringatan Stok Menipis',
                          style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                        ),
                        Text('23 produk memerlukan restock segera'),
                      ],
                    ),
                  ),
                  ElevatedButton.icon(
                    onPressed: () {
                      _generatePurchaseOrder();
                    },
                    icon: const Icon(Icons.shopping_cart),
                    label: const Text('Buat PO'),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          const Text(
            'Produk dengan Stok Menipis',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 16),

          Expanded(
            child: ListView.builder(
              itemCount: 10,
              itemBuilder: (context, index) {
                return _buildLowStockCard(index);
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAddProductTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Tambah Produk Baru',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 24),

              // Product information
              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Kode Produk *', 'SPR001'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Nama Produk *', 'Oli Mesin 10W-40'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildDropdownField('Kategori *', ['Sparepart', 'Oil', 'Tire', 'Battery', 'Filter']),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Supplier', 'PT Supplier Indo'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Satuan', 'Liter'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildDropdownField('Serial Number?', ['Ya', 'Tidak']),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              // Pricing and stock
              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Harga Beli *', 'Rp 45.000'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Harga Jual *', 'Rp 65.000'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Stok Awal', '100'),
                  ),
                ],
              ),
              const SizedBox(height: 16),

              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Min Stok Level', '10'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Max Stok Level', '500'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Container(), // Spacer
                  ),
                ],
              ),
              const SizedBox(height: 16),

              _buildFormField('Deskripsi Produk', 'Deskripsi detail produk...', maxLines: 3),
              const SizedBox(height: 32),

              // Action buttons
              Row(
                children: [
                  ElevatedButton.icon(
                    onPressed: () {
                      _saveProduct();
                    },
                    icon: const Icon(Icons.save),
                    label: const Text('Simpan Produk'),
                    style: ElevatedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                    ),
                  ),
                  const SizedBox(width: 16),
                  OutlinedButton.icon(
                    onPressed: () {
                      _clearProductForm();
                    },
                    icon: const Icon(Icons.clear),
                    label: const Text('Clear Form'),
                  ),
                  const SizedBox(width: 16),
                  ElevatedButton.icon(
                    onPressed: () {
                      _importProducts();
                    },
                    icon: const Icon(Icons.file_upload),
                    label: const Text('Import Excel'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.green,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildReportsTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        children: [
          // Report cards
          Row(
            children: [
              Expanded(
                child: _buildReportCard(
                  'Inventory Valuation',
                  'Nilai total stok per kategori',
                  Icons.assessment,
                  Colors.blue,
                  () => _showInventoryValuation(),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: _buildReportCard(
                  'Stock Movement',
                  'Pergerakan stok masuk dan keluar',
                  Icons.trending_up,
                  Colors.green,
                  () => _showStockMovement(),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          Row(
            children: [
              Expanded(
                child: _buildReportCard(
                  'ABC Analysis',
                  'Analisis produk berdasarkan nilai',
                  Icons.analytics,
                  Colors.orange,
                  () => _showABCAnalysis(),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: _buildReportCard(
                  'Stock Aging',
                  'Laporan umur stok di gudang',
                  Icons.access_time,
                  Colors.purple,
                  () => _showStockAging(),
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),

          // Quick chart
          Expanded(
            child: Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Top 10 Produk Terlaris',
                      style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 16),
                    Expanded(
                      child: ListView.builder(
                        itemCount: 10,
                        itemBuilder: (context, index) {
                          return ListTile(
                            leading: CircleAvatar(
                              backgroundColor: Colors.blue.shade100,
                              child: Text('${index + 1}'),
                            ),
                            title: Text('Produk ${index + 1}'),
                            subtitle: Text('Terjual: ${(100 - index * 8)} unit'),
                            trailing: Text('Rp ${(50 - index * 3)}M'),
                          );
                        },
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

  Widget _buildProductRow(int index) {
    final products = [
      {'name': 'Oli Mesin 10W-40', 'code': 'OIL001', 'category': 'Oil', 'stock': '45', 'buyPrice': '45.000', 'sellPrice': '65.000'},
      {'name': 'Ban Dalam Motor', 'code': 'TIR001', 'category': 'Tire', 'stock': '12', 'buyPrice': '25.000', 'sellPrice': '35.000'},
      {'name': 'Aki Motor 12V', 'code': 'BAT001', 'category': 'Battery', 'stock': '8', 'buyPrice': '180.000', 'sellPrice': '250.000'},
    ];

    final product = products[index % products.length];
    final stockCount = int.parse(product['stock']!);
    final stockColor = stockCount < 15 ? Colors.red : stockCount < 30 ? Colors.orange : Colors.green;

    return Container(
      margin: const EdgeInsets.only(bottom: 4),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.grey.shade200),
      ),
      child: Row(
        children: [
          Expanded(
            flex: 2,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  product['name']!,
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
                Text(
                  product['code']!,
                  style: TextStyle(color: Colors.grey.shade600, fontSize: 12),
                ),
              ],
            ),
          ),
          Expanded(
            child: Chip(
              label: Text(product['category']!),
              backgroundColor: Colors.blue.shade100,
              labelStyle: TextStyle(color: Colors.blue.shade700, fontSize: 12),
            ),
          ),
          Expanded(
            child: Text(
              product['stock']!,
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: stockColor,
              ),
            ),
          ),
          Expanded(child: Text('Rp ${product['buyPrice']}')),
          Expanded(child: Text('Rp ${product['sellPrice']}')),
          SizedBox(
            width: 100,
            child: PopupMenuButton(
              itemBuilder: (context) => [
                PopupMenuItem(
                  child: const Text('Edit'),
                  onTap: () => _editProduct(product),
                ),
                PopupMenuItem(
                  child: const Text('Stock In'),
                  onTap: () => _stockIn(product),
                ),
                PopupMenuItem(
                  child: const Text('Stock Out'),
                  onTap: () => _stockOut(product),
                ),
                PopupMenuItem(
                  child: const Text('History'),
                  onTap: () => _viewStockHistory(product),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLowStockCard(int index) {
    final products = [
      {'name': 'Ban Dalam Motor', 'stock': '5', 'min': '15', 'supplier': 'PT Tire Indonesia'},
      {'name': 'Filter Udara', 'stock': '3', 'min': '20', 'supplier': 'PT Filter Jaya'},
      {'name': 'Busi Platinum', 'stock': '8', 'min': '25', 'supplier': 'PT Spark Indo'},
    ];

    final product = products[index % products.length];

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Colors.red.shade100,
          child: Icon(Icons.warning, color: Colors.red),
        ),
        title: Text(product['name']!),
        subtitle: Text('Stok: ${product['stock']} (Min: ${product['min']}) • ${product['supplier']}'),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              '${int.parse(product['min']!) - int.parse(product['stock']!)} needed',
              style: const TextStyle(color: Colors.red, fontWeight: FontWeight.bold),
            ),
            const SizedBox(width: 8),
            ElevatedButton(
              onPressed: () => _reorderProduct(product),
              child: const Text('Reorder'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFormField(String label, String hint, {int maxLines = 1}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(fontWeight: FontWeight.w500)),
        const SizedBox(height: 4),
        TextField(
          maxLines: maxLines,
          decoration: InputDecoration(
            hintText: hint,
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(8)),
            contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          ),
        ),
      ],
    );
  }

  Widget _buildDropdownField(String label, List<String> items) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(fontWeight: FontWeight.w500)),
        const SizedBox(height: 4),
        DropdownButtonFormField<String>(
          decoration: InputDecoration(
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(8)),
            contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          ),
          items: items.map((String item) {
            return DropdownMenuItem<String>(
              value: item,
              child: Text(item),
            );
          }).toList(),
          onChanged: (String? newValue) {
            // Handle dropdown change
          },
        ),
      ],
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
  void _exportInventoryData() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Data inventory berhasil diekspor')),
    );
  }

  void _generatePurchaseOrder() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Purchase Order berhasil dibuat')),
    );
  }

  void _saveProduct() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Produk berhasil disimpan')),
    );
  }

  void _clearProductForm() {
    // TODO: Clear all form fields
  }

  void _importProducts() {
    // TODO: Show import dialog
  }

  void _editProduct(Map<String, String> product) {
    // TODO: Show edit product dialog
  }

  void _stockIn(Map<String, String> product) {
    // TODO: Show stock in dialog
  }

  void _stockOut(Map<String, String> product) {
    // TODO: Show stock out dialog
  }

  void _viewStockHistory(Map<String, String> product) {
    // TODO: Navigate to stock history page
  }

  void _reorderProduct(Map<String, String> product) {
    // TODO: Show reorder dialog
  }

  void _showInventoryValuation() {
    // TODO: Show inventory valuation report
  }

  void _showStockMovement() {
    // TODO: Show stock movement report
  }

  void _showABCAnalysis() {
    // TODO: Show ABC analysis report
  }

  void _showStockAging() {
    // TODO: Show stock aging report
  }
}