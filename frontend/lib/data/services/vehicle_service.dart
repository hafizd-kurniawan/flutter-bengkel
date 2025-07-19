import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/vehicle_inventory_model.dart';
import '../models/vehicle_purchase_model.dart';
import '../models/vehicle_sale_model.dart';

// Mock data providers for now - will be replaced with actual API calls
class VehicleService {
  // Mock vehicle inventory data
  static List<VehicleInventoryModel> _mockInventory = [
    VehicleInventoryModel(
      inventoryId: 1,
      plateNumber: 'B 1234 ABC',
      brand: 'Toyota',
      model: 'Avanza',
      type: 'MPV',
      productionYear: 2019,
      chassisNumber: 'MHFM1BA5XJK123456',
      engineNumber: '3SZ-VE123456',
      color: 'Silver Metallic',
      mileage: 45000,
      conditionRating: 4,
      purchasePrice: 135000000,
      estimatedSellingPrice: 155000000,
      status: 'Available',
      vehiclePhotos: [
        'assets/images/vehicles/avanza_front.jpg',
        'assets/images/vehicles/avanza_side.jpg',
      ],
      conditionNotes: 'Kondisi baik, service record lengkap',
      createdAt: DateTime.now().subtract(const Duration(days: 15)),
    ),
    VehicleInventoryModel(
      inventoryId: 2,
      plateNumber: 'B 5678 DEF',
      brand: 'Honda',
      model: 'Civic',
      type: 'Sedan',
      productionYear: 2020,
      chassisNumber: 'JHGFD1BA5XJK789012',
      engineNumber: 'L15B-789012',
      color: 'White Pearl',
      mileage: 32000,
      conditionRating: 5,
      purchasePrice: 285000000,
      estimatedSellingPrice: 320000000,
      status: 'Available',
      vehiclePhotos: [
        'assets/images/vehicles/civic_front.jpg',
        'assets/images/vehicles/civic_interior.jpg',
      ],
      conditionNotes: 'Seperti baru, terawat sekali',
      createdAt: DateTime.now().subtract(const Duration(days: 8)),
    ),
    VehicleInventoryModel(
      inventoryId: 3,
      plateNumber: 'B 9012 GHI',
      brand: 'Mitsubishi',
      model: 'Xpander',
      type: 'MPV',
      productionYear: 2021,
      chassisNumber: 'MMBJNKA5XJK345678',
      engineNumber: '4A91-345678',
      color: 'Diamond Black',
      mileage: 18000,
      conditionRating: 5,
      purchasePrice: 195000000,
      estimatedSellingPrice: 225000000,
      actualSellingPrice: 220000000,
      status: 'Sold',
      vehiclePhotos: [
        'assets/images/vehicles/xpander_front.jpg',
      ],
      conditionNotes: 'Masih dalam garansi',
      sellingDate: DateTime.now().subtract(const Duration(days: 3)),
      profitMargin: 25000000,
      createdAt: DateTime.now().subtract(const Duration(days: 25)),
    ),
  ];

  static List<VehiclePurchaseModel> _mockPurchases = [
    VehiclePurchaseModel(
      purchaseId: 1,
      customerId: 101,
      customerName: 'Budi Santoso',
      outletId: 1,
      purchaseDate: DateTime.now().subtract(const Duration(days: 15)),
      purchasePrice: 135000000,
      paymentMethod: 'cash',
      notes: 'Pembelian tunai dari customer walk-in',
      status: 'completed',
      createdAt: DateTime.now().subtract(const Duration(days: 15)),
    ),
    VehiclePurchaseModel(
      purchaseId: 2,
      customerId: 102,
      customerName: 'Siti Rahayu',
      outletId: 1,
      purchaseDate: DateTime.now().subtract(const Duration(days: 8)),
      purchasePrice: 285000000,
      paymentMethod: 'bank_transfer',
      notes: 'Trade-in dengan mobil baru',
      status: 'completed',
      createdAt: DateTime.now().subtract(const Duration(days: 8)),
    ),
  ];

  static List<VehicleSaleModel> _mockSales = [
    VehicleSaleModel(
      saleId: 1,
      inventoryId: 3,
      customerId: 201,
      customerName: 'Ahmad Firdaus',
      outletId: 1,
      saleDate: DateTime.now().subtract(const Duration(days: 3)),
      sellingPrice: 220000000,
      profitAmount: 25000000,
      paymentMethod: 'financing',
      financingType: 'bank_financing',
      downPayment: 66000000,
      tenorMonths: 48,
      monthlyInstallment: 4200000,
      notes: 'Kredit melalui Bank Mandiri',
      status: 'completed',
      createdAt: DateTime.now().subtract(const Duration(days: 3)),
    ),
  ];

  // Inventory operations
  Future<List<VehicleInventoryModel>> getVehicleInventory({
    String? status,
    String? search,
  }) async {
    await Future.delayed(const Duration(milliseconds: 500)); // Simulate API delay
    
    var result = _mockInventory;
    
    if (status != null) {
      result = result.where((v) => v.status == status).toList();
    }
    
    if (search != null && search.isNotEmpty) {
      result = result.where((v) => 
        v.brand.toLowerCase().contains(search.toLowerCase()) ||
        v.model.toLowerCase().contains(search.toLowerCase()) ||
        v.plateNumber.toLowerCase().contains(search.toLowerCase())
      ).toList();
    }
    
    return result;
  }

  Future<VehicleInventoryModel?> getVehicleById(int id) async {
    await Future.delayed(const Duration(milliseconds: 300));
    
    try {
      return _mockInventory.firstWhere((v) => v.inventoryId == id);
    } catch (e) {
      return null;
    }
  }

  Future<bool> updateVehiclePrice(int inventoryId, double newPrice) async {
    await Future.delayed(const Duration(milliseconds: 500));
    
    final index = _mockInventory.indexWhere((v) => v.inventoryId == inventoryId);
    if (index != -1) {
      _mockInventory[index] = _mockInventory[index].copyWith(
        estimatedSellingPrice: newPrice,
      );
      return true;
    }
    return false;
  }

  Future<bool> markVehicleAsSold(int inventoryId, double sellingPrice) async {
    await Future.delayed(const Duration(milliseconds: 500));
    
    final index = _mockInventory.indexWhere((v) => v.inventoryId == inventoryId);
    if (index != -1) {
      final vehicle = _mockInventory[index];
      _mockInventory[index] = vehicle.copyWith(
        status: 'Sold',
        actualSellingPrice: sellingPrice,
        sellingDate: DateTime.now(),
        profitMargin: sellingPrice - vehicle.purchasePrice,
      );
      return true;
    }
    return false;
  }

  // Purchase operations
  Future<List<VehiclePurchaseModel>> getVehiclePurchases({
    String? status,
    DateTime? fromDate,
    DateTime? toDate,
  }) async {
    await Future.delayed(const Duration(milliseconds: 500));
    
    var result = _mockPurchases;
    
    if (status != null) {
      result = result.where((p) => p.status == status).toList();
    }
    
    if (fromDate != null) {
      result = result.where((p) => p.purchaseDate.isAfter(fromDate)).toList();
    }
    
    if (toDate != null) {
      result = result.where((p) => p.purchaseDate.isBefore(toDate)).toList();
    }
    
    return result;
  }

  Future<VehiclePurchaseModel> createVehiclePurchase(Map<String, dynamic> purchaseData) async {
    await Future.delayed(const Duration(milliseconds: 800));
    
    final newPurchase = VehiclePurchaseModel(
      purchaseId: _mockPurchases.length + 1,
      customerId: purchaseData['customer_id'],
      customerName: purchaseData['customer_name'],
      outletId: purchaseData['outlet_id'],
      purchaseDate: DateTime.parse(purchaseData['purchase_date']),
      purchasePrice: purchaseData['purchase_price'],
      paymentMethod: purchaseData['payment_method'],
      notes: purchaseData['notes'],
      status: 'completed',
      createdAt: DateTime.now(),
    );
    
    _mockPurchases.add(newPurchase);
    
    // Also add to inventory
    final newInventory = VehicleInventoryModel(
      inventoryId: _mockInventory.length + 1,
      plateNumber: purchaseData['plate_number'],
      brand: purchaseData['brand'],
      model: purchaseData['model'],
      type: purchaseData['type'],
      productionYear: purchaseData['production_year'],
      chassisNumber: purchaseData['chassis_number'],
      engineNumber: purchaseData['engine_number'],
      color: purchaseData['color'],
      mileage: purchaseData['mileage'],
      conditionRating: purchaseData['condition_rating'],
      purchasePrice: purchaseData['purchase_price'],
      estimatedSellingPrice: purchaseData['estimated_selling_price'],
      status: 'Available',
      vehiclePhotos: [],
      conditionNotes: purchaseData['condition_notes'],
      createdAt: DateTime.now(),
    );
    
    _mockInventory.add(newInventory);
    
    return newPurchase;
  }

  // Sales operations  
  Future<List<VehicleSaleModel>> getVehicleSales({
    String? status,
    DateTime? fromDate,
    DateTime? toDate,
  }) async {
    await Future.delayed(const Duration(milliseconds: 500));
    
    var result = _mockSales;
    
    if (status != null) {
      result = result.where((s) => s.status == status).toList();
    }
    
    if (fromDate != null) {
      result = result.where((s) => s.saleDate.isAfter(fromDate)).toList();
    }
    
    if (toDate != null) {
      result = result.where((s) => s.saleDate.isBefore(toDate)).toList();
    }
    
    return result;
  }

  Future<VehicleSaleModel> createVehicleSale(Map<String, dynamic> saleData) async {
    await Future.delayed(const Duration(milliseconds: 800));
    
    final newSale = VehicleSaleModel(
      saleId: _mockSales.length + 1,
      inventoryId: saleData['inventory_id'],
      customerId: saleData['customer_id'],
      customerName: saleData['customer_name'],
      outletId: saleData['outlet_id'],
      saleDate: DateTime.now(),
      sellingPrice: saleData['selling_price'],
      profitAmount: saleData['profit_amount'],
      paymentMethod: saleData['payment_method'],
      financingType: saleData['financing_type'],
      downPayment: saleData['down_payment'],
      tenorMonths: saleData['tenor_months'],
      monthlyInstallment: saleData['monthly_installment'],
      notes: saleData['notes'],
      status: 'completed',
      createdAt: DateTime.now(),
    );
    
    _mockSales.add(newSale);
    
    // Mark vehicle as sold
    await markVehicleAsSold(saleData['inventory_id'], saleData['selling_price']);
    
    return newSale;
  }

  // Statistics
  Future<Map<String, dynamic>> getVehicleStatistics() async {
    await Future.delayed(const Duration(milliseconds: 300));
    
    final totalStock = _mockInventory.length;
    final availableStock = _mockInventory.where((v) => v.status == 'Available').length;
    final soldThisMonth = _mockSales.where((s) => 
      s.saleDate.month == DateTime.now().month &&
      s.saleDate.year == DateTime.now().year
    ).length;
    final totalStockValue = _mockInventory
        .where((v) => v.status == 'Available')
        .fold(0.0, (sum, v) => sum + v.purchasePrice);
    
    return {
      'total_stock': totalStock,
      'available_stock': availableStock,
      'sold_this_month': soldThisMonth,
      'total_stock_value': totalStockValue,
    };
  }
}

// Providers
final vehicleServiceProvider = Provider<VehicleService>((ref) {
  return VehicleService();
});

final vehicleInventoryProvider = FutureProvider.family<List<VehicleInventoryModel>, Map<String, dynamic>>((ref, params) async {
  final service = ref.watch(vehicleServiceProvider);
  return service.getVehicleInventory(
    status: params['status'],
    search: params['search'],
  );
});

final vehicleStatisticsProvider = FutureProvider<Map<String, dynamic>>((ref) async {
  final service = ref.watch(vehicleServiceProvider);
  return service.getVehicleStatistics();
});

final vehiclePurchasesProvider = FutureProvider.family<List<VehiclePurchaseModel>, Map<String, dynamic>>((ref, params) async {
  final service = ref.watch(vehicleServiceProvider);
  return service.getVehiclePurchases(
    status: params['status'],
    fromDate: params['from_date'],
    toDate: params['to_date'],
  );
});

final vehicleSalesProvider = FutureProvider.family<List<VehicleSaleModel>, Map<String, dynamic>>((ref, params) async {
  final service = ref.watch(vehicleServiceProvider);
  return service.getVehicleSales(
    status: params['status'],
    fromDate: params['from_date'],
    toDate: params['to_date'],
  );
});