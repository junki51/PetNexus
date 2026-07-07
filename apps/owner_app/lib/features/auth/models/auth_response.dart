class AuthResponse {
  final String accessToken;
  final String message;

  const AuthResponse({
    required this.accessToken,
    required this.message,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    // Backend returns: {success, message, data: {accessToken, ...}}
    final data = json['data'] is Map ? json['data'] as Map<String, dynamic> : json;
    return AuthResponse(
      accessToken: data['accessToken'] as String? ??
          data['access_token'] as String? ?? '',
      message: json['message'] as String? ?? '',
    );
  }
}