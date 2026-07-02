class AuthResponse {
  final String accessToken;
  final String? refreshToken;
  final String message;

  const AuthResponse({
    required this.accessToken,
    this.refreshToken,
    required this.message,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      accessToken: json["access_token"],
      refreshToken: json["refresh_token"],
      message: json["message"] ?? "",
    );
  }
}