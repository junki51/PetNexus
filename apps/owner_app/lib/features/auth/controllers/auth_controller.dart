import 'package:flutter/material.dart';

import '../models/user_model.dart';
import '../repositories/auth_repository.dart';

enum AuthState {
  initial,
  loading,
  authenticated,
  unauthenticated,
  error,
}

class AuthController extends ChangeNotifier {
  final AuthRepository _repository = AuthRepository();

  //======================
  // State
  //======================

  AuthState _state = AuthState.initial;
  AuthState get state => _state;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  UserModel? _currentUser;
  UserModel? get currentUser => _currentUser;

  //======================
  // UI State
  //======================

  bool _isPasswordVisible = false;
  bool get isPasswordVisible => _isPasswordVisible;

  bool _isConfirmPasswordVisible = false;
  bool get isConfirmPasswordVisible => _isConfirmPasswordVisible;

  bool _acceptedTerms = false;
  bool get acceptedTerms => _acceptedTerms;

  //======================
  // UI Methods
  //======================

  void togglePasswordVisibility() {
    _isPasswordVisible = !_isPasswordVisible;
    notifyListeners();
  }

  void toggleConfirmPasswordVisibility() {
    _isConfirmPasswordVisible = !_isConfirmPasswordVisible;
    notifyListeners();
  }

  void toggleAcceptedTerms(bool? value) {
    _acceptedTerms = value ?? false;
    notifyListeners();
  }

  void resetState() {
    _errorMessage = null;
    _state = AuthState.initial;
    notifyListeners();
  }

  //======================
  // Login
  //======================

  Future<bool> login({
    required String email,
    required String password,
  }) async {
    _setLoading();

    try {
      await _repository.login(
        email: email,
        password: password,
      );

      await loadCurrentUser();

      _state = AuthState.authenticated;
      notifyListeners();

      return true;
    } catch (e) {
      _setError(e.toString());
      return false;
    }
  }

  //======================
  // Register
  //======================

  Future<bool> register({
    required String email,
    required String password,
    required String confirmPassword,
  }) async {
    if (password != confirmPassword) {
      _setError("Passwords do not match");
      return false;
    }

    _setLoading();

    try {
      await _repository.register(
        email: email,
        password: password,
        confirmPassword: confirmPassword,
      );

      await loadCurrentUser();

      _state = AuthState.authenticated;
      notifyListeners();

      return true;
    } catch (e) {
      _setError(e.toString());
      return false;
    }
  }

  //======================
  // Auto Login
  //======================

  Future<bool> checkAuthentication() async {
    _setLoading();

    final token = await _repository.getToken();

    if (token == null) {
      _state = AuthState.unauthenticated;
      notifyListeners();
      return false;
    }

    try {
      await loadCurrentUser();

      _state = AuthState.authenticated;
      notifyListeners();

      return true;
    } catch (_) {
      await logout();
      return false;
    }
  }

  //======================
  // Current User
  //======================

    Future<void> loadCurrentUser() async {
    _currentUser = await _repository.getCurrentUser();
    notifyListeners();
  }

  //======================
  // Logout
  //======================

  Future<void> logout() async {
    await _repository.logout();

    _currentUser = null;
    _errorMessage = null;
    _state = AuthState.unauthenticated;

    notifyListeners();
  }

  //======================
  // Helper
  //======================

  void _setLoading() {
    _errorMessage = null;
    _state = AuthState.loading;
    notifyListeners();
  }

  void _setError(String message) {
    _errorMessage = message;
    _state = AuthState.error;
    notifyListeners();
  }
}