class OwnerProfileModel {
  final String id;
  final String userId;
  final String firstName;
  final String lastName;
  final String displayName;
  final String? gender;
  final String? dateOfBirth;
  final String phoneNumber;
  final String? avatarUrl;
  final String? addressLine1;
  final String? addressLine2;
  final String? province;
  final String? district;
  final String? subdistrict;
  final String? postalCode;
  final String createdAt;
  final String updatedAt;

  const OwnerProfileModel({
    required this.id,
    required this.userId,
    required this.firstName,
    required this.lastName,
    required this.displayName,
    this.gender,
    this.dateOfBirth,
    required this.phoneNumber,
    this.avatarUrl,
    this.addressLine1,
    this.addressLine2,
    this.province,
    this.district,
    this.subdistrict,
    this.postalCode,
    required this.createdAt,
    required this.updatedAt,
  });

  factory OwnerProfileModel.fromJson(Map<String, dynamic> json) {
    return OwnerProfileModel(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      firstName: json['first_name'] as String? ?? '',
      lastName: json['last_name'] as String? ?? '',
      displayName: json['display_name'] as String? ?? '',
      gender: json['gender'] as String?,
      dateOfBirth: json['date_of_birth'] as String?,
      phoneNumber: json['phone_number'] as String? ?? '',
      avatarUrl: json['avatar_url'] as String?,
      addressLine1: json['address_line1'] as String?,
      addressLine2: json['address_line2'] as String?,
      province: json['province'] as String?,
      district: json['district'] as String?,
      subdistrict: json['subdistrict'] as String?,
      postalCode: json['postal_code'] as String?,
      createdAt: json['created_at'] as String? ?? '',
      updatedAt: json['updated_at'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'first_name': firstName,
      'last_name': lastName,
      'gender': gender,
      'date_of_birth': dateOfBirth,
      'phone_number': phoneNumber,
      'avatar_url': avatarUrl,
      'address_line1': addressLine1,
      'address_line2': addressLine2,
      'province': province,
      'district': district,
      'subdistrict': subdistrict,
      'postal_code': postalCode,
    };
  }
}
