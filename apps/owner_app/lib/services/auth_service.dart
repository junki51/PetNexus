import 'package:http/http.dart' as http;
import 'dart:convert';
import '../../api_config.dart';

class AuthService {
  Future<Map<String, dynamic>> login({
    required String email,
    required String password,
  }) async {
    final response = await http.post(
      Uri.parse(ApiConfig.login),
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'email': email,
        'password': password,
      }),
    );

    final data = jsonDecode(response.body);

    if (response.statusCode == 200) {
      return {
        "success": true,
        "accessToken": data["accessToken"],
      };
    }

    return {
      "success": false,
      "message": data["message"] ?? "เข้าสู่ระบบไม่สำเร็จ",
    };
  }
}