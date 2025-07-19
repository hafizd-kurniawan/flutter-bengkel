import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/customer_model.dart';

class CustomerService {
  // Mock customer data
  static List<CustomerModel> _mockCustomers = [
    CustomerModel(
      customerId: 1,
      fullName: 'Budi Santoso',
      email: 'budi.santoso@email.com',
      phone: '08123456789',
      address: 'Jl. Merdeka No. 10, Jakarta',
      city: 'Jakarta',
      postalCode: '12345',
      identityNumber: '1234567890123456',
      customerType: 'individual',
      isActive: true,
      loyaltyPoints: 1250,
      totalSpending: 15000000,
      lastVisit: DateTime.now().subtract(const Duration(days: 3)),
      createdAt: DateTime.now().subtract(const Duration(days: 45)),
    ),
    CustomerModel(
      customerId: 2,
      fullName: 'Siti Rahayu',
      email: 'siti.rahayu@company.com',
      phone: '08234567890',
      address: 'Jl. Sudirman No. 25, Jakarta',
      city: 'Jakarta',
      postalCode: '12346',
      identityNumber: '2345678901234567',
      customerType: 'corporate',
      isActive: true,
      loyaltyPoints: 2850,
      totalSpending: 35000000,
      lastVisit: DateTime.now().subtract(const Duration(days: 1)),
      createdAt: DateTime.now().subtract(const Duration(days: 120)),
    ),
    CustomerModel(
      customerId: 3,
      fullName: 'Ahmad Firdaus',
      phone: '08345678901',
      address: 'Jl. Gatot Subroto No. 15, Tangerang',
      city: 'Tangerang',
      postalCode: '15111',
      customerType: 'individual',
      isActive: true,
      loyaltyPoints: 850,
      totalSpending: 8500000,
      lastVisit: DateTime.now().subtract(const Duration(days: 7)),
      createdAt: DateTime.now().subtract(const Duration(days: 20)),
    ),
    CustomerModel(
      customerId: 4,
      fullName: 'PT Maju Bersama',
      email: 'admin@majubersama.co.id',
      phone: '02154321098',
      address: 'Jl. HR Rasuna Said No. 50, Jakarta',
      city: 'Jakarta',
      postalCode: '12940',
      identityNumber: '9876543210987654',
      customerType: 'corporate',
      isActive: true,
      loyaltyPoints: 5200,
      totalSpending: 75000000,
      lastVisit: DateTime.now().subtract(const Duration(days: 2)),
      createdAt: DateTime.now().subtract(const Duration(days: 200)),
    ),
    CustomerModel(
      customerId: 5,
      fullName: 'Maria Gonzalez',
      email: 'maria.g@email.com',
      phone: '08456789012',
      address: 'Jl. Kemang Raya No. 88, Jakarta',
      city: 'Jakarta',
      postalCode: '12560',
      customerType: 'individual',
      isActive: true,
      loyaltyPoints: 450,
      totalSpending: 4500000,
      lastVisit: DateTime.now().subtract(const Duration(days: 14)),
      createdAt: DateTime.now().subtract(const Duration(days: 10)),
    ),
    CustomerModel(
      customerId: 6,
      fullName: 'John Smith',
      phone: '08567890123',
      address: 'Jl. Thamrin No. 100, Jakarta',
      city: 'Jakarta',
      customerType: 'individual',
      isActive: false,
      loyaltyPoints: 0,
      totalSpending: 1200000,
      lastVisit: DateTime.now().subtract(const Duration(days: 180)),
      createdAt: DateTime.now().subtract(const Duration(days: 365)),
    ),
  ];

  // Get all customers
  Future<List<CustomerModel>> getCustomers({
    String? search,
    String? customerType,
    bool? isActive,
    String? sortBy,
  }) async {
    await Future.delayed(const Duration(milliseconds: 500)); // Simulate API delay
    
    var result = _mockCustomers;
    
    // Apply filters
    if (search != null && search.isNotEmpty) {
      result = result.where((customer) => 
        customer.fullName.toLowerCase().contains(search.toLowerCase()) ||
        customer.phone.contains(search) ||
        (customer.email?.toLowerCase().contains(search.toLowerCase()) ?? false) ||
        (customer.address?.toLowerCase().contains(search.toLowerCase()) ?? false)
      ).toList();
    }
    
    if (customerType != null && customerType.isNotEmpty) {
      result = result.where((customer) => customer.customerType == customerType).toList();
    }
    
    if (isActive != null) {
      result = result.where((customer) => customer.isActive == isActive).toList();
    }
    
    // Apply sorting
    if (sortBy != null) {
      switch (sortBy) {
        case 'name':
          result.sort((a, b) => a.fullName.compareTo(b.fullName));
          break;
        case 'spending':
          result.sort((a, b) => b.totalSpending.compareTo(a.totalSpending));
          break;
        case 'points':
          result.sort((a, b) => b.loyaltyPoints.compareTo(a.loyaltyPoints));
          break;
        case 'recent':
          result.sort((a, b) => {
            if (a.lastVisit == null && b.lastVisit == null) return 0;
            if (a.lastVisit == null) return 1;
            if (b.lastVisit == null) return -1;
            return b.lastVisit!.compareTo(a.lastVisit!);
          });
          break;
        default:
          result.sort((a, b) => b.createdAt.compareTo(a.createdAt));
      }
    }
    
    return result;
  }

  // Get customer by ID
  Future<CustomerModel?> getCustomerById(int customerId) async {
    await Future.delayed(const Duration(milliseconds: 300));
    
    try {
      return _mockCustomers.firstWhere((customer) => customer.customerId == customerId);
    } catch (e) {
      return null;
    }
  }

  // Create new customer
  Future<CustomerModel> createCustomer(Map<String, dynamic> customerData) async {
    await Future.delayed(const Duration(milliseconds: 800));
    
    final newCustomer = CustomerModel(
      customerId: _mockCustomers.length + 1,
      fullName: customerData['full_name'],
      email: customerData['email'],
      phone: customerData['phone'],
      address: customerData['address'],
      city: customerData['city'],
      postalCode: customerData['postal_code'],
      identityNumber: customerData['identity_number'],
      customerType: customerData['customer_type'] ?? 'individual',
      isActive: true,
      loyaltyPoints: 0,
      totalSpending: 0.0,
      createdAt: DateTime.now(),
    );
    
    _mockCustomers.add(newCustomer);
    return newCustomer;
  }

  // Update customer
  Future<CustomerModel?> updateCustomer(int customerId, Map<String, dynamic> updateData) async {
    await Future.delayed(const Duration(milliseconds: 600));
    
    final index = _mockCustomers.indexWhere((customer) => customer.customerId == customerId);
    if (index != -1) {
      final customer = _mockCustomers[index];
      _mockCustomers[index] = customer.copyWith(
        fullName: updateData['full_name'] ?? customer.fullName,
        email: updateData['email'] ?? customer.email,
        phone: updateData['phone'] ?? customer.phone,
        address: updateData['address'] ?? customer.address,
        city: updateData['city'] ?? customer.city,
        postalCode: updateData['postal_code'] ?? customer.postalCode,
        identityNumber: updateData['identity_number'] ?? customer.identityNumber,
        customerType: updateData['customer_type'] ?? customer.customerType,
        isActive: updateData['is_active'] ?? customer.isActive,
      );
      return _mockCustomers[index];
    }
    return null;
  }

  // Delete customer (soft delete)
  Future<bool> deleteCustomer(int customerId) async {
    await Future.delayed(const Duration(milliseconds: 500));
    
    final index = _mockCustomers.indexWhere((customer) => customer.customerId == customerId);
    if (index != -1) {
      _mockCustomers[index] = _mockCustomers[index].copyWith(isActive: false);
      return true;
    }
    return false;
  }

  // Get customer statistics
  Future<Map<String, dynamic>> getCustomerStatistics() async {
    await Future.delayed(const Duration(milliseconds: 300));
    
    final activeCustomers = _mockCustomers.where((c) => c.isActive).length;
    final totalCustomers = _mockCustomers.length;
    final newCustomers = _mockCustomers.where((c) => c.isNewCustomer && c.isActive).length;
    final totalLoyaltyPoints = _mockCustomers
        .where((c) => c.isActive)
        .fold(0, (sum, c) => sum + c.loyaltyPoints);
    final totalRevenue = _mockCustomers
        .where((c) => c.isActive)
        .fold(0.0, (sum, c) => sum + c.totalSpending);
    final vipCustomers = _mockCustomers.where((c) => c.isVip && c.isActive).length;
    
    return {
      'total_customers': totalCustomers,
      'active_customers': activeCustomers,
      'new_customers': newCustomers,
      'total_loyalty_points': totalLoyaltyPoints,
      'total_revenue': totalRevenue,
      'vip_customers': vipCustomers,
      'customer_retention_rate': totalCustomers > 0 ? (activeCustomers / totalCustomers * 100).toDouble() : 0.0,
      'avg_spending': activeCustomers > 0 ? (totalRevenue / activeCustomers) : 0.0,
    };
  }

  // Search customers by various criteria
  Future<List<CustomerModel>> searchCustomers(String query) async {
    await Future.delayed(const Duration(milliseconds: 400));
    
    return _mockCustomers.where((customer) =>
      customer.fullName.toLowerCase().contains(query.toLowerCase()) ||
      customer.phone.contains(query) ||
      (customer.email?.toLowerCase().contains(query.toLowerCase()) ?? false) ||
      (customer.address?.toLowerCase().contains(query.toLowerCase()) ?? false) ||
      customer.customerId.toString() == query
    ).toList();
  }

  // Get top customers by spending
  Future<List<CustomerModel>> getTopCustomers({int limit = 10}) async {
    await Future.delayed(const Duration(milliseconds: 300));
    
    final sortedCustomers = List<CustomerModel>.from(_mockCustomers)
      ..sort((a, b) => b.totalSpending.compareTo(a.totalSpending));
    
    return sortedCustomers.take(limit).toList();
  }

  // Get customers by loyalty tier
  Future<Map<String, List<CustomerModel>>> getCustomersByTier() async {
    await Future.delayed(const Duration(milliseconds: 400));
    
    final bronze = <CustomerModel>[];
    final silver = <CustomerModel>[];
    final gold = <CustomerModel>[];
    final platinum = <CustomerModel>[];
    
    for (final customer in _mockCustomers.where((c) => c.isActive)) {
      if (customer.loyaltyPoints >= 5000) {
        platinum.add(customer);
      } else if (customer.loyaltyPoints >= 2000) {
        gold.add(customer);
      } else if (customer.loyaltyPoints >= 500) {
        silver.add(customer);
      } else {
        bronze.add(customer);
      }
    }
    
    return {
      'bronze': bronze,
      'silver': silver,
      'gold': gold,
      'platinum': platinum,
    };
  }

  // Award loyalty points
  Future<bool> awardLoyaltyPoints(int customerId, int points) async {
    await Future.delayed(const Duration(milliseconds: 300));
    
    final index = _mockCustomers.indexWhere((customer) => customer.customerId == customerId);
    if (index != -1) {
      final customer = _mockCustomers[index];
      _mockCustomers[index] = customer.copyWith(
        loyaltyPoints: customer.loyaltyPoints + points,
      );
      return true;
    }
    return false;
  }

  // Update last visit
  Future<bool> updateLastVisit(int customerId) async {
    await Future.delayed(const Duration(milliseconds: 200));
    
    final index = _mockCustomers.indexWhere((customer) => customer.customerId == customerId);
    if (index != -1) {
      final customer = _mockCustomers[index];
      _mockCustomers[index] = customer.copyWith(lastVisit: DateTime.now());
      return true;
    }
    return false;
  }
}

// Providers
final customerServiceProvider = Provider<CustomerService>((ref) {
  return CustomerService();
});

final customersProvider = FutureProvider.family<List<CustomerModel>, Map<String, dynamic>>((ref, params) async {
  final service = ref.watch(customerServiceProvider);
  return service.getCustomers(
    search: params['search'],
    customerType: params['customer_type'],
    isActive: params['is_active'],
    sortBy: params['sort_by'],
  );
});

final customerStatisticsProvider = FutureProvider<Map<String, dynamic>>((ref) async {
  final service = ref.watch(customerServiceProvider);
  return service.getCustomerStatistics();
});

final topCustomersProvider = FutureProvider.family<List<CustomerModel>, int>((ref, limit) async {
  final service = ref.watch(customerServiceProvider);
  return service.getTopCustomers(limit: limit);
});

final customersByTierProvider = FutureProvider<Map<String, List<CustomerModel>>>((ref) async {
  final service = ref.watch(customerServiceProvider);
  return service.getCustomersByTier();
});