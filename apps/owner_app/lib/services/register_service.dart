import 'dart:convert';

import 'package:http/http.dart' as http;

import '../api_config.dart';

class RegisterService {
  Future<Map<String, dynamic>> register({
    required String email,
    required String password,
    required String confirmPassword,
  }) async {
    final response = await http.post(
      Uri.parse(ApiConfig.register),
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'email': email,
        'password': password,
        'password_confirmation': confirmPassword,
      }),
    );

    final data = jsonDecode(response.body);

    if (response.statusCode == 200 || response.statusCode == 201) {
      return {
        "success": true,
      };
    }

    return {
      "success": false,
      "message": data["message"] ?? "สมัครสมาชิกไม่สำเร็จ",
    };
  }
}