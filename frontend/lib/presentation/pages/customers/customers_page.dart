import 'package:flutter/material.dart';
import '../../../shared/widgets/app_drawer.dart';
import '../../../core/constants/app_constants.dart';

class CustomersPage extends StatefulWidget {
  const CustomersPage({super.key});

  @override
  State<CustomersPage> createState() => _CustomersPageState();
}

class _CustomersPageState extends State<CustomersPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
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
        title: const Text(AppStrings.customers),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(icon: Icon(Icons.people), text: 'Semua Customer'),
            Tab(icon: Icon(Icons.person_add), text: 'Tambah Customer'),
            Tab(icon: Icon(Icons.analytics), text: 'Customer Analytics'),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              _showAddCustomerDialog(context);
            },
          ),
          IconButton(
            icon: const Icon(Icons.download),
            onPressed: () {
              _exportCustomerData();
            },
          ),
        ],
      ),
      drawer: const AppDrawer(),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildCustomerListTab(),
          _buildAddCustomerTab(),
          _buildAnalyticsTab(),
        ],
      ),
    );
  }

  Widget _buildCustomerListTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        children: [
          // Search and filter row
          Row(
            children: [
              Expanded(
                flex: 3,
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Cari customer berdasarkan nama, telepon, atau kendaraan...',
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
              ElevatedButton.icon(
                onPressed: () {
                  _showFilterDialog();
                },
                icon: const Icon(Icons.filter_list),
                label: const Text('Filter'),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // Statistics row
          Row(
            children: [
              _buildStatCard('Total Customer', '1,247', Icons.people, Colors.blue),
              const SizedBox(width: 16),
              _buildStatCard('Customer Aktif', '892', Icons.person, Colors.green),
              const SizedBox(width: 16),
              _buildStatCard('Customer Baru', '45', Icons.person_add, Colors.orange),
              const SizedBox(width: 16),
              _buildStatCard('Loyalty Points', '25.4K', Icons.star, Colors.purple),
            ],
          ),
          const SizedBox(height: 24),

          // Customer list
          Expanded(
            child: ListView.builder(
              itemCount: 20, // Mock data
              itemBuilder: (context, index) {
                return _buildCustomerCard(index);
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAddCustomerTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Tambah Customer Baru',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 24),
              
              // Personal Information
              const Text(
                'Informasi Personal',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
              ),
              const SizedBox(height: 12),
              
              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Nama Lengkap *', 'Masukkan nama lengkap'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Email', 'customer@email.com'),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              
              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Nomor Telepon *', '+62 812 3456 7890'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Tanggal Lahir', 'dd/mm/yyyy'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildDropdownField('Jenis Kelamin', ['Laki-laki', 'Perempuan']),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              
              _buildFormField('Alamat', 'Masukkan alamat lengkap', maxLines: 3),
              const SizedBox(height: 16),
              
              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Kota', 'Jakarta'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Provinsi', 'DKI Jakarta'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Kode Pos', '12345'),
                  ),
                ],
              ),
              const SizedBox(height: 24),
              
              // Vehicle Information
              const Text(
                'Informasi Kendaraan (Opsional)',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
              ),
              const SizedBox(height: 12),
              
              Row(
                children: [
                  Expanded(
                    child: _buildFormField('Plat Nomor', 'B 1234 ABC'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Merk Kendaraan', 'Toyota'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Model', 'Avanza'),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildFormField('Tahun', '2020'),
                  ),
                ],
              ),
              const SizedBox(height: 32),
              
              // Action buttons
              Row(
                children: [
                  ElevatedButton.icon(
                    onPressed: () {
                      _saveCustomer();
                    },
                    icon: const Icon(Icons.save),
                    label: const Text('Simpan Customer'),
                    style: ElevatedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                    ),
                  ),
                  const SizedBox(width: 16),
                  OutlinedButton.icon(
                    onPressed: () {
                      _clearCustomerForm();
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
    );
  }

  Widget _buildAnalyticsTab() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        children: [
          // Analytics cards
          Row(
            children: [
              _buildAnalyticsCard('Customer Baru Bulan Ini', '127', Icons.trending_up, Colors.green),
              const SizedBox(width: 16),
              _buildAnalyticsCard('Average Kunjungan', '3.2x', Icons.repeat, Colors.blue),
              const SizedBox(width: 16),
              _buildAnalyticsCard('Customer Rating', '4.8/5', Icons.star, Colors.orange),
              const SizedBox(width: 16),
              _buildAnalyticsCard('Retention Rate', '89%', Icons.favorite, Colors.purple),
            ],
          ),
          const SizedBox(height: 24),
          
          // Charts and reports
          Expanded(
            child: Row(
              children: [
                Expanded(
                  child: Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Customer Growth',
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
                                  'Chart Pertumbuhan Customer\n6 Bulan Terakhir',
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
                const SizedBox(width: 16),
                Expanded(
                  child: Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Top Customers',
                            style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                          ),
                          const SizedBox(height: 16),
                          Expanded(
                            child: ListView.builder(
                              itemCount: 10,
                              itemBuilder: (context, index) {
                                return ListTile(
                                  leading: CircleAvatar(
                                    child: Text('${index + 1}'),
                                  ),
                                  title: Text('Customer ${index + 1}'),
                                  subtitle: Text('Total: Rp ${(50 - index * 3)}M'),
                                  trailing: Text('${(25 - index * 2)} kunjungan'),
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

  Widget _buildAnalyticsCard(String title, String value, IconData icon, Color color) {
    return Expanded(
      child: Card(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            children: [
              Icon(icon, color: color, size: 32),
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
                textAlign: TextAlign.center,
                style: const TextStyle(fontSize: 12, color: Colors.grey),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildCustomerCard(int index) {
    final customers = [
      {'name': 'Budi Santoso', 'phone': '+62 812 3456 7890', 'vehicle': 'Toyota Avanza 2020 (B 1234 ABC)', 'visits': '12', 'lastVisit': '2 hari lalu'},
      {'name': 'Siti Aminah', 'phone': '+62 821 9876 5432', 'vehicle': 'Honda Civic 2019 (B 5678 DEF)', 'visits': '8', 'lastVisit': '1 minggu lalu'},
      {'name': 'Ahmad Rahman', 'phone': '+62 813 2468 1357', 'vehicle': 'Suzuki Ertiga 2021 (B 9012 GHI)', 'visits': '5', 'lastVisit': '3 hari lalu'},
    ];
    
    final customer = customers[index % customers.length];
    
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          children: [
            CircleAvatar(
              backgroundColor: Colors.blue.shade100,
              child: Text(
                customer['name']!.substring(0, 1),
                style: const TextStyle(fontWeight: FontWeight.bold),
              ),
            ),
            const SizedBox(width: 16),
            
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    customer['name']!,
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                  Text('📞 ${customer['phone']}'),
                  Text('🚗 ${customer['vehicle']}'),
                  const SizedBox(height: 4),
                  Row(
                    children: [
                      Chip(
                        label: Text('${customer['visits']} kunjungan'),
                        backgroundColor: Colors.green.shade100,
                        labelStyle: TextStyle(color: Colors.green.shade700, fontSize: 12),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        'Terakhir: ${customer['lastVisit']}',
                        style: const TextStyle(color: Colors.grey, fontSize: 12),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            
            PopupMenuButton(
              itemBuilder: (context) => [
                PopupMenuItem(
                  child: const Text('Lihat Detail'),
                  onTap: () => _viewCustomerDetail(customer),
                ),
                PopupMenuItem(
                  child: const Text('Edit Customer'),
                  onTap: () => _editCustomer(customer),
                ),
                PopupMenuItem(
                  child: const Text('Riwayat Service'),
                  onTap: () => _viewServiceHistory(customer),
                ),
                PopupMenuItem(
                  child: const Text('Tambah Kendaraan'),
                  onTap: () => _addVehicle(customer),
                ),
              ],
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

  // Action methods
  void _showAddCustomerDialog(BuildContext context) {
    _tabController.animateTo(1); // Switch to add customer tab
  }

  void _exportCustomerData() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Data customer berhasil diekspor')),
    );
  }

  void _showFilterDialog() {
    // TODO: Implement filter dialog
  }

  void _saveCustomer() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Customer berhasil disimpan')),
    );
  }

  void _clearCustomerForm() {
    // TODO: Clear all form fields
  }

  void _viewCustomerDetail(Map<String, String> customer) {
    // TODO: Navigate to customer detail page
  }

  void _editCustomer(Map<String, String> customer) {
    // TODO: Show edit customer dialog
  }

  void _viewServiceHistory(Map<String, String> customer) {
    // TODO: Navigate to service history page
  }

  void _addVehicle(Map<String, String> customer) {
    // TODO: Show add vehicle dialog
  }
}