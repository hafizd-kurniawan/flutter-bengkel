class CustomerModel {
  final int customerId;
  final String fullName;
  final String? email;
  final String phone;
  final String? address;
  final String? city;
  final String? postalCode;
  final String? identityNumber;
  final String customerType;
  final bool isActive;
  final int loyaltyPoints;
  final double totalSpending;
  final DateTime? lastVisit;
  final DateTime createdAt;

  const CustomerModel({
    required this.customerId,
    required this.fullName,
    this.email,
    required this.phone,
    this.address,
    this.city,
    this.postalCode,
    this.identityNumber,
    required this.customerType,
    required this.isActive,
    required this.loyaltyPoints,
    required this.totalSpending,
    this.lastVisit,
    required this.createdAt,
  });

  factory CustomerModel.fromJson(Map<String, dynamic> json) {
    return CustomerModel(
      customerId: json['customer_id'] as int,
      fullName: json['full_name'] as String,
      email: json['email'] as String?,
      phone: json['phone'] as String,
      address: json['address'] as String?,
      city: json['city'] as String?,
      postalCode: json['postal_code'] as String?,
      identityNumber: json['identity_number'] as String?,
      customerType: json['customer_type'] as String,
      isActive: json['is_active'] as bool? ?? true,
      loyaltyPoints: json['loyalty_points'] as int? ?? 0,
      totalSpending: (json['total_spending'] as num?)?.toDouble() ?? 0.0,
      lastVisit: json['last_visit'] != null 
          ? DateTime.parse(json['last_visit'] as String)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'customer_id': customerId,
      'full_name': fullName,
      'email': email,
      'phone': phone,
      'address': address,
      'city': city,
      'postal_code': postalCode,
      'identity_number': identityNumber,
      'customer_type': customerType,
      'is_active': isActive,
      'loyalty_points': loyaltyPoints,
      'total_spending': totalSpending,
      'last_visit': lastVisit?.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
    };
  }

  CustomerModel copyWith({
    int? customerId,
    String? fullName,
    String? email,
    String? phone,
    String? address,
    String? city,
    String? postalCode,
    String? identityNumber,
    String? customerType,
    bool? isActive,
    int? loyaltyPoints,
    double? totalSpending,
    DateTime? lastVisit,
    DateTime? createdAt,
  }) {
    return CustomerModel(
      customerId: customerId ?? this.customerId,
      fullName: fullName ?? this.fullName,
      email: email ?? this.email,
      phone: phone ?? this.phone,
      address: address ?? this.address,
      city: city ?? this.city,
      postalCode: postalCode ?? this.postalCode,
      identityNumber: identityNumber ?? this.identityNumber,
      customerType: customerType ?? this.customerType,
      isActive: isActive ?? this.isActive,
      loyaltyPoints: loyaltyPoints ?? this.loyaltyPoints,
      totalSpending: totalSpending ?? this.totalSpending,
      lastVisit: lastVisit ?? this.lastVisit,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'CustomerModel(customerId: $customerId, fullName: $fullName, phone: $phone)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    
    return other is CustomerModel &&
      other.customerId == customerId;
  }

  @override
  int get hashCode => customerId.hashCode;

  // Helper getters
  String get displayName => fullName;
  String get initials => fullName.split(' ').take(2).map((n) => n[0]).join().toUpperCase();
  String get formattedSpending => 'Rp ${totalSpending.toStringAsFixed(0)}';
  String get customerTypeDisplay => customerType == 'individual' ? 'Perorangan' : 'Perusahaan';
  bool get isNewCustomer => DateTime.now().difference(createdAt).inDays <= 30;
  bool get isVip => totalSpending >= 10000000; // VIP if spent > 10M
  int get daysSinceLastVisit => lastVisit != null 
      ? DateTime.now().difference(lastVisit!).inDays 
      : 999;
}