import 'package:dio/dio.dart';
import '../../../../core/network/api_client.dart';
import '../models/owner_profile_model.dart';

class OwnerProfileRepository {
  Future<OwnerProfileModel> createProfile({
    required String firstName,
    required String lastName,
    required String phoneNumber,
    String? gender,
    String? dateOfBirth,
    String? avatarUrl,
    String? addressLine1,
    String? addressLine2,
    String? province,
    String? district,
    String? subdistrict,
    String? postalCode,
  }) async {
    try {
      final response = await ApiClient.instance.dio.post(
        "/owner/profile",
        data: {
          "first_name": firstName,
          "last_name": lastName,
          "phone_number": phoneNumber,
          "gender": gender,
          "date_of_birth": dateOfBirth,
          "avatar_url": avatarUrl,
          "address_line1": addressLine1,
          "address_line2": addressLine2,
          "province": province,
          "district": district,
          "subdistrict": subdistrict,
          "postal_code": postalCode,
        },
      );

      // Backend wraps response in a success utility: {"success": true, "message": "...", "data": {...}}
      final data = response.data['data'] as Map<String, dynamic>;
      return OwnerProfileModel.fromJson(data);
    } on DioException catch (e) {
      throw Exception(e.response?.data["message"] ?? "Failed to create profile");
    }
  }

  Future<OwnerProfileModel> getProfile() async {
    try {
      final response = await ApiClient.instance.dio.get(
        "/owner/profile",
      );
      final data = response.data['data'] as Map<String, dynamic>;
      return OwnerProfileModel.fromJson(data);
    } on DioException catch (e) {
      throw Exception(e.response?.data["message"] ?? "Failed to fetch profile");
    }
  }

  Future<OwnerProfileModel> updateProfile({
    String? firstName,
    String? lastName,
    String? phoneNumber,
    String? gender,
    String? dateOfBirth,
    String? avatarUrl,
    String? addressLine1,
    String? addressLine2,
    String? province,
    String? district,
    String? subdistrict,
    String? postalCode,
  }) async {
    try {
      final Map<String, dynamic> updateData = {};
      if (firstName != null) updateData["first_name"] = firstName;
      if (lastName != null) updateData["last_name"] = lastName;
      if (phoneNumber != null) updateData["phone_number"] = phoneNumber;
      if (gender != null) updateData["gender"] = gender;
      if (dateOfBirth != null) updateData["date_of_birth"] = dateOfBirth;
      if (avatarUrl != null) updateData["avatar_url"] = avatarUrl;
      if (addressLine1 != null) updateData["address_line1"] = addressLine1;
      if (addressLine2 != null) updateData["address_line2"] = addressLine2;
      if (province != null) updateData["province"] = province;
      if (district != null) updateData["district"] = district;
      if (subdistrict != null) updateData["subdistrict"] = subdistrict;
      if (postalCode != null) updateData["postal_code"] = postalCode;

      final response = await ApiClient.instance.dio.patch(
        "/owner/profile",
        data: updateData,
      );
      final data = response.data['data'] as Map<String, dynamic>;
      return OwnerProfileModel.fromJson(data);
    } on DioException catch (e) {
      throw Exception(e.response?.data["message"] ?? "Failed to update profile");
    }
  }
}
