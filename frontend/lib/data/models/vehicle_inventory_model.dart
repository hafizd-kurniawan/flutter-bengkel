class VehicleInventoryModel {
  final int inventoryId;
  final String plateNumber;
  final String brand;
  final String model;
  final String type;
  final int productionYear;
  final String chassisNumber;
  final String engineNumber;
  final String color;
  final int mileage;
  final int conditionRating;
  final double purchasePrice;
  final double estimatedSellingPrice;
  final double? actualSellingPrice;
  final String status;
  final List<String> vehiclePhotos;
  final String? conditionNotes;
  final DateTime? sellingDate;
  final double? profitMargin;
  final DateTime createdAt;

  const VehicleInventoryModel({
    required this.inventoryId,
    required this.plateNumber,
    required this.brand,
    required this.model,
    required this.type,
    required this.productionYear,
    required this.chassisNumber,
    required this.engineNumber,
    required this.color,
    required this.mileage,
    required this.conditionRating,
    required this.purchasePrice,
    required this.estimatedSellingPrice,
    this.actualSellingPrice,
    required this.status,
    required this.vehiclePhotos,
    this.conditionNotes,
    this.sellingDate,
    this.profitMargin,
    required this.createdAt,
  });

  factory VehicleInventoryModel.fromJson(Map<String, dynamic> json) {
    return VehicleInventoryModel(
      inventoryId: json['inventory_id'] as int,
      plateNumber: json['plate_number'] as String,
      brand: json['brand'] as String,
      model: json['model'] as String,
      type: json['type'] as String,
      productionYear: json['production_year'] as int,
      chassisNumber: json['chassis_number'] as String,
      engineNumber: json['engine_number'] as String,
      color: json['color'] as String,
      mileage: json['mileage'] as int,
      conditionRating: json['condition_rating'] as int,
      purchasePrice: (json['purchase_price'] as num).toDouble(),
      estimatedSellingPrice: (json['estimated_selling_price'] as num).toDouble(),
      actualSellingPrice: json['actual_selling_price'] != null 
          ? (json['actual_selling_price'] as num).toDouble() 
          : null,
      status: json['status'] as String,
      vehiclePhotos: List<String>.from(json['vehicle_photos'] ?? []),
      conditionNotes: json['condition_notes'] as String?,
      sellingDate: json['selling_date'] != null 
          ? DateTime.parse(json['selling_date'] as String)
          : null,
      profitMargin: json['profit_margin'] != null 
          ? (json['profit_margin'] as num).toDouble()
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'inventory_id': inventoryId,
      'plate_number': plateNumber,
      'brand': brand,
      'model': model,
      'type': type,
      'production_year': productionYear,
      'chassis_number': chassisNumber,
      'engine_number': engineNumber,
      'color': color,
      'mileage': mileage,
      'condition_rating': conditionRating,
      'purchase_price': purchasePrice,
      'estimated_selling_price': estimatedSellingPrice,
      'actual_selling_price': actualSellingPrice,
      'status': status,
      'vehicle_photos': vehiclePhotos,
      'condition_notes': conditionNotes,
      'selling_date': sellingDate?.toIso8601String(),
      'profit_margin': profitMargin,
      'created_at': createdAt.toIso8601String(),
    };
  }

  VehicleInventoryModel copyWith({
    int? inventoryId,
    String? plateNumber,
    String? brand,
    String? model,
    String? type,
    int? productionYear,
    String? chassisNumber,
    String? engineNumber,
    String? color,
    int? mileage,
    int? conditionRating,
    double? purchasePrice,
    double? estimatedSellingPrice,
    double? actualSellingPrice,
    String? status,
    List<String>? vehiclePhotos,
    String? conditionNotes,
    DateTime? sellingDate,
    double? profitMargin,
    DateTime? createdAt,
  }) {
    return VehicleInventoryModel(
      inventoryId: inventoryId ?? this.inventoryId,
      plateNumber: plateNumber ?? this.plateNumber,
      brand: brand ?? this.brand,
      model: model ?? this.model,
      type: type ?? this.type,
      productionYear: productionYear ?? this.productionYear,
      chassisNumber: chassisNumber ?? this.chassisNumber,
      engineNumber: engineNumber ?? this.engineNumber,
      color: color ?? this.color,
      mileage: mileage ?? this.mileage,
      conditionRating: conditionRating ?? this.conditionRating,
      purchasePrice: purchasePrice ?? this.purchasePrice,
      estimatedSellingPrice: estimatedSellingPrice ?? this.estimatedSellingPrice,
      actualSellingPrice: actualSellingPrice ?? this.actualSellingPrice,
      status: status ?? this.status,
      vehiclePhotos: vehiclePhotos ?? this.vehiclePhotos,
      conditionNotes: conditionNotes ?? this.conditionNotes,
      sellingDate: sellingDate ?? this.sellingDate,
      profitMargin: profitMargin ?? this.profitMargin,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'VehicleInventoryModel(inventoryId: $inventoryId, plateNumber: $plateNumber, brand: $brand, model: $model)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    
    return other is VehicleInventoryModel &&
      other.inventoryId == inventoryId &&
      other.plateNumber == plateNumber &&
      other.chassisNumber == chassisNumber;
  }

  @override
  int get hashCode {
    return Object.hash(inventoryId, plateNumber, chassisNumber);
  }

  // Helper getters
  String get displayName => '$brand $model ($productionYear)';
  String get formattedPrice => 'Rp ${purchasePrice.toStringAsFixed(0)}';
  String get formattedEstimatedPrice => 'Rp ${estimatedSellingPrice.toStringAsFixed(0)}';
  bool get isAvailable => status == 'Available';
  bool get isSold => status == 'Sold';
  int get daysSincePurchase => DateTime.now().difference(createdAt).inDays;
}