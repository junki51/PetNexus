class UserModel {
  final String id;
  final String email;
  final String? phone;
  final String role;
  final String createdAt;

  const UserModel({
    required this.id,
    required this.email,
    this.phone,
    required this.role,
    required this.createdAt,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    // Backend wraps in {success, data: {...}} for /me
    final data = json['data'] is Map ? json['data'] as Map<String, dynamic> : json;
    return UserModel(
      id: data['id'] as String? ?? '',
      email: data['email'] as String? ?? '',
      phone: data['phone'] as String?,
      role: data['role'] as String? ?? 'owner',
      createdAt: data['createdAt'] as String? ??
          data['created_at'] as String? ?? '',
    );
  }
}