class VehicleSaleModel {
  final int saleId;
  final int inventoryId;
  final int customerId;
  final String customerName;
  final int outletId;
  final DateTime saleDate;
  final double sellingPrice;
  final double profitAmount;
  final String paymentMethod;
  final String financingType;
  final double? downPayment;
  final int? tenorMonths;
  final double? monthlyInstallment;
  final String? notes;
  final String status;
  final DateTime createdAt;

  const VehicleSaleModel({
    required this.saleId,
    required this.inventoryId,
    required this.customerId,
    required this.customerName,
    required this.outletId,
    required this.saleDate,
    required this.sellingPrice,
    required this.profitAmount,
    required this.paymentMethod,
    required this.financingType,
    this.downPayment,
    this.tenorMonths,
    this.monthlyInstallment,
    this.notes,
    required this.status,
    required this.createdAt,
  });

  factory VehicleSaleModel.fromJson(Map<String, dynamic> json) {
    return VehicleSaleModel(
      saleId: json['sale_id'] as int,
      inventoryId: json['inventory_id'] as int,
      customerId: json['customer_id'] as int,
      customerName: json['customer_name'] as String,
      outletId: json['outlet_id'] as int,
      saleDate: DateTime.parse(json['sale_date'] as String),
      sellingPrice: (json['selling_price'] as num).toDouble(),
      profitAmount: (json['profit_amount'] as num).toDouble(),
      paymentMethod: json['payment_method'] as String,
      financingType: json['financing_type'] as String,
      downPayment: json['down_payment'] != null 
          ? (json['down_payment'] as num).toDouble() 
          : null,
      tenorMonths: json['tenor_months'] as int?,
      monthlyInstallment: json['monthly_installment'] != null 
          ? (json['monthly_installment'] as num).toDouble()
          : null,
      notes: json['notes'] as String?,
      status: json['status'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'sale_id': saleId,
      'inventory_id': inventoryId,
      'customer_id': customerId,
      'customer_name': customerName,
      'outlet_id': outletId,
      'sale_date': saleDate.toIso8601String().split('T')[0],
      'selling_price': sellingPrice,
      'profit_amount': profitAmount,
      'payment_method': paymentMethod,
      'financing_type': financingType,
      'down_payment': downPayment,
      'tenor_months': tenorMonths,
      'monthly_installment': monthlyInstallment,
      'notes': notes,
      'status': status,
      'created_at': createdAt.toIso8601String(),
    };
  }

  VehicleSaleModel copyWith({
    int? saleId,
    int? inventoryId,
    int? customerId,
    String? customerName,
    int? outletId,
    DateTime? saleDate,
    double? sellingPrice,
    double? profitAmount,
    String? paymentMethod,
    String? financingType,
    double? downPayment,
    int? tenorMonths,
    double? monthlyInstallment,
    String? notes,
    String? status,
    DateTime? createdAt,
  }) {
    return VehicleSaleModel(
      saleId: saleId ?? this.saleId,
      inventoryId: inventoryId ?? this.inventoryId,
      customerId: customerId ?? this.customerId,
      customerName: customerName ?? this.customerName,
      outletId: outletId ?? this.outletId,
      saleDate: saleDate ?? this.saleDate,
      sellingPrice: sellingPrice ?? this.sellingPrice,
      profitAmount: profitAmount ?? this.profitAmount,
      paymentMethod: paymentMethod ?? this.paymentMethod,
      financingType: financingType ?? this.financingType,
      downPayment: downPayment ?? this.downPayment,
      tenorMonths: tenorMonths ?? this.tenorMonths,
      monthlyInstallment: monthlyInstallment ?? this.monthlyInstallment,
      notes: notes ?? this.notes,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'VehicleSaleModel(saleId: $saleId, customerName: $customerName, sellingPrice: $sellingPrice, profitAmount: $profitAmount)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    
    return other is VehicleSaleModel &&
      other.saleId == saleId;
  }

  @override
  int get hashCode => saleId.hashCode;

  // Helper getters
  String get formattedPrice => 'Rp ${sellingPrice.toStringAsFixed(0)}';
  String get formattedProfit => 'Rp ${profitAmount.toStringAsFixed(0)}';
  String get formattedDate => '${saleDate.day}/${saleDate.month}/${saleDate.year}';
  String get profitPercentage => '${((profitAmount / (sellingPrice - profitAmount)) * 100).toStringAsFixed(1)}%';
  bool get isCompleted => status == 'completed';
  bool get isPending => status == 'pending';
  bool get isCash => paymentMethod == 'cash';
  bool get isFinanced => financingType != 'cash';
}