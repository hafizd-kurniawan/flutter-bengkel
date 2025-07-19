class VehiclePurchaseModel {
  final int purchaseId;
  final int customerId;
  final String customerName;
  final int outletId;
  final DateTime purchaseDate;
  final double purchasePrice;
  final String paymentMethod;
  final String? notes;
  final String status;
  final DateTime createdAt;

  const VehiclePurchaseModel({
    required this.purchaseId,
    required this.customerId,
    required this.customerName,
    required this.outletId,
    required this.purchaseDate,
    required this.purchasePrice,
    required this.paymentMethod,
    this.notes,
    required this.status,
    required this.createdAt,
  });

  factory VehiclePurchaseModel.fromJson(Map<String, dynamic> json) {
    return VehiclePurchaseModel(
      purchaseId: json['purchase_id'] as int,
      customerId: json['customer_id'] as int,
      customerName: json['customer_name'] as String,
      outletId: json['outlet_id'] as int,
      purchaseDate: DateTime.parse(json['purchase_date'] as String),
      purchasePrice: (json['purchase_price'] as num).toDouble(),
      paymentMethod: json['payment_method'] as String,
      notes: json['notes'] as String?,
      status: json['status'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'purchase_id': purchaseId,
      'customer_id': customerId,
      'customer_name': customerName,
      'outlet_id': outletId,
      'purchase_date': purchaseDate.toIso8601String().split('T')[0],
      'purchase_price': purchasePrice,
      'payment_method': paymentMethod,
      'notes': notes,
      'status': status,
      'created_at': createdAt.toIso8601String(),
    };
  }

  VehiclePurchaseModel copyWith({
    int? purchaseId,
    int? customerId,
    String? customerName,
    int? outletId,
    DateTime? purchaseDate,
    double? purchasePrice,
    String? paymentMethod,
    String? notes,
    String? status,
    DateTime? createdAt,
  }) {
    return VehiclePurchaseModel(
      purchaseId: purchaseId ?? this.purchaseId,
      customerId: customerId ?? this.customerId,
      customerName: customerName ?? this.customerName,
      outletId: outletId ?? this.outletId,
      purchaseDate: purchaseDate ?? this.purchaseDate,
      purchasePrice: purchasePrice ?? this.purchasePrice,
      paymentMethod: paymentMethod ?? this.paymentMethod,
      notes: notes ?? this.notes,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'VehiclePurchaseModel(purchaseId: $purchaseId, customerName: $customerName, purchasePrice: $purchasePrice)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    
    return other is VehiclePurchaseModel &&
      other.purchaseId == purchaseId;
  }

  @override
  int get hashCode => purchaseId.hashCode;

  // Helper getters
  String get formattedPrice => 'Rp ${purchasePrice.toStringAsFixed(0)}';
  String get formattedDate => '${purchaseDate.day}/${purchaseDate.month}/${purchaseDate.year}';
  bool get isCompleted => status == 'completed';
  bool get isPending => status == 'pending';
}