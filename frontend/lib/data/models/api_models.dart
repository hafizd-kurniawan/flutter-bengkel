// API models for the new PostgreSQL backend with UUID and Soft Delete

class ApiResponse<T> {
  final bool success;
  final String message;
  final T? data;
  final String? error;

  const ApiResponse({
    required this.success,
    required this.message,
    this.data,
    this.error,
  });

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Object?)? fromJsonT,
  ) {
    return ApiResponse<T>(
      success: json['success'] as bool,
      message: json['message'] as String,
      data: json['data'] != null && fromJsonT != null
          ? fromJsonT(json['data'])
          : null,
      error: json['error'] as String?,
    );
  }
}

class OutletModel {
  final String id; // UUID
  final String name;
  final String address;
  final String phone;
  final String email;
  final bool isActive;
  final DateTime createdAt;
  final DateTime updatedAt;
  final DateTime? deletedAt;
  final String? deletedBy; // UUID

  const OutletModel({
    required this.id,
    required this.name,
    required this.address,
    required this.phone,
    required this.email,
    required this.isActive,
    required this.createdAt,
    required this.updatedAt,
    this.deletedAt,
    this.deletedBy,
  });

  factory OutletModel.fromJson(Map<String, dynamic> json) {
    return OutletModel(
      id: json['id'] as String,
      name: json['name'] as String,
      address: json['address'] as String,
      phone: json['phone'] as String,
      email: json['email'] as String,
      isActive: json['is_active'] as bool,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      deletedAt: json['deleted_at'] != null
          ? DateTime.parse(json['deleted_at'] as String)
          : null,
      deletedBy: json['deleted_by'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'address': address,
      'phone': phone,
      'email': email,
      'is_active': isActive,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      if (deletedAt != null) 'deleted_at': deletedAt!.toIso8601String(),
      if (deletedBy != null) 'deleted_by': deletedBy,
    };
  }

  OutletModel copyWith({
    String? id,
    String? name,
    String? address,
    String? phone,
    String? email,
    bool? isActive,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? deletedAt,
    String? deletedBy,
  }) {
    return OutletModel(
      id: id ?? this.id,
      name: name ?? this.name,
      address: address ?? this.address,
      phone: phone ?? this.phone,
      email: email ?? this.email,
      isActive: isActive ?? this.isActive,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      deletedAt: deletedAt ?? this.deletedAt,
      deletedBy: deletedBy ?? this.deletedBy,
    );
  }

  bool get isDeleted => deletedAt != null;

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is OutletModel && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;

  @override
  String toString() {
    return 'OutletModel(id: $id, name: $name, isDeleted: $isDeleted)';
  }
}

class CreateOutletRequest {
  final String name;
  final String address;
  final String phone;
  final String email;

  const CreateOutletRequest({
    required this.name,
    required this.address,
    required this.phone,
    required this.email,
  });

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'address': address,
      'phone': phone,
      'email': email,
    };
  }
}