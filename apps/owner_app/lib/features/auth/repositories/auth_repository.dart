import 'package:dio/dio.dart';
import 'package:owner_app/features/auth/models/user_model.dart';
import '../../../../core/network/api_client.dart';
import '../../../../core/storage/token_storage.dart';
import '../models/auth_response.dart';

class AuthRepository {
  Future<AuthResponse> login({
    required String email,
    required String password,
  }) async {
    try {
      final response =
          await ApiClient.instance.dio.post(
        "/auth/login",
        data: {
          "email": email,
          "password": password,
        },
      );

      final auth =
          AuthResponse.fromJson(response.data);

      await TokenStorage.instance.saveToken(
          auth.accessToken);

      return auth;
    } on DioException catch (e) {
      throw Exception(e.response?.data["message"] ?? "Login failed");
    }
  }

  Future<AuthResponse> register({
    required String email,
    required String password,
    required String confirmPassword,
  }) async {
    try {
      final response = await ApiClient.instance.dio.post(
        "/auth/register",
        data: {
          "email": email,
          "password": password,
          "confirmPassword": confirmPassword,
        },
      );
      final auth =
    AuthResponse.fromJson(response.data);

    await TokenStorage.instance.saveToken(auth.accessToken);

    return auth;
    } on DioException catch (e) {
      throw Exception(e.response?.data["message"] ?? "Registration failed");
    }
  }

  Future<UserModel> getCurrentUser() async {
    try {
      final response = await ApiClient.instance.dio.get(
        "/auth/me",
      );

      return UserModel.fromJson(response.data);

    } on DioException catch (e) {
      throw Exception(
        e.response?.data["message"] ?? "Unable to load profile",
      );
    }
  }
  Future<String?> getToken() async {
    return await TokenStorage.instance.getToken();
  }

  Future<void> logout() async {
    await TokenStorage.instance.deleteToken();
  }
}
