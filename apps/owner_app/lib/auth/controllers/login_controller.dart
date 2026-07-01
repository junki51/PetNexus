import 'package:flutter/material.dart';
import '../../services/auth_service.dart';

// กำหนดสถานะของหน้าจอ Login
enum AuthState { initial, loading, success, error }

class LoginController extends ChangeNotifier {
  final AuthService _authService = AuthService();

  String? _accessToken;
  String? get accessToken => _accessToken;

  AuthState _state = AuthState.initial;
  AuthState get state => _state;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  bool _isPasswordVisible = false;
  bool get isPasswordVisible => _isPasswordVisible;

  void togglePasswordVisibility() {
    _isPasswordVisible = !_isPasswordVisible;
    notifyListeners();
  }

  Future<bool> loginWithEmail({
    required String email,
    required String password,
  }) async {
    _setState(AuthState.loading);

    try {
      final result = await _authService.login(
        email: email,
        password: password,
      );

      if (result["success"]) {
        _accessToken = result["accessToken"];
        _setState(AuthState.success);
        return true;
      }
      _errorMessage = result["message"];
      _setState(AuthState.error);
      return false;

    } catch (e) {
      _errorMessage = e.toString();
      _setState(AuthState.error);
      return false;
    }
  }

  Future<void> loginWithSocial(String provider) async {
    _setState(AuthState.loading);
    try {
      // 
      await Future.delayed(const Duration(seconds: 2));
      _setState(AuthState.success);
    } catch (e) {
      _errorMessage = "Social Login Error: $e";
      _setState(AuthState.error);
    }
  }

  void _setState(AuthState newState) {
    _state = newState;
    notifyListeners(); // แจ้งเตือน UI ให้ Re-build
  }
  void logout() {
    _accessToken = null;
    _state = AuthState.initial;
    notifyListeners();
  }
}
