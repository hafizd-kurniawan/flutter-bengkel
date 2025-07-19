class UserModel {
  final int id;
  final String username;
  final String email;
  final String fullName;
  final String? phone;
  final String roleName;
  final String? outletName;
  final bool isActive;
  final DateTime? lastLoginAt;
  final DateTime createdAt;

  const UserModel({
    required this.id,
    required this.username,
    required this.email,
    required this.fullName,
    this.phone,
    required this.roleName,
    this.outletName,
    this.isActive = true,
    this.lastLoginAt,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as int,
      username: json['username'] as String,
      email: json['email'] as String,
      fullName: json['full_name'] as String,
      phone: json['phone'] as String?,
      roleName: json['role']?['name'] as String? ?? 'User',
      outletName: json['outlet']?['name'] as String?,
      isActive: json['is_active'] as bool? ?? true,
      lastLoginAt: json['last_login_at'] != null 
          ? DateTime.parse(json['last_login_at'] as String)
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'username': username,
      'email': email,
      'full_name': fullName,
      'phone': phone,
      'role': {'name': roleName},
      'outlet': outletName != null ? {'name': outletName} : null,
      'is_active': isActive,
      'last_login_at': lastLoginAt?.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
    };
  }

  UserModel copyWith({
    int? id,
    String? username,
    String? email,
    String? fullName,
    String? phone,
    String? roleName,
    String? outletName,
    bool? isActive,
    DateTime? lastLoginAt,
    DateTime? createdAt,
  }) {
    return UserModel(
      id: id ?? this.id,
      username: username ?? this.username,
      email: email ?? this.email,
      fullName: fullName ?? this.fullName,
      phone: phone ?? this.phone,
      roleName: roleName ?? this.roleName,
      outletName: outletName ?? this.outletName,
      isActive: isActive ?? this.isActive,
      lastLoginAt: lastLoginAt ?? this.lastLoginAt,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'UserModel(id: $id, username: $username, email: $email, fullName: $fullName)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    
    return other is UserModel &&
      other.id == id &&
      other.username == username &&
      other.email == email &&
      other.fullName == fullName &&
      other.phone == phone &&
      other.roleName == roleName &&
      other.outletName == outletName &&
      other.isActive == isActive &&
      other.lastLoginAt == lastLoginAt &&
      other.createdAt == createdAt;
  }

  @override
  int get hashCode {
    return Object.hash(
      id,
      username,
      email,
      fullName,
      phone,
      roleName,
      outletName,
      isActive,
      lastLoginAt,
      createdAt,
    );
  }
}