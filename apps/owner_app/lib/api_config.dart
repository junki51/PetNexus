class ApiConfig {
  // 10.0.2.2 คือ IP พิเศษสำหรับ Android Emulator ที่วิ่งไปหา Localhost ของเครื่องคอมฯ นั้นๆ
  static const String baseUrl = "http://10.0.2.2:8080"; 
  static const String apiAuth = "$baseUrl/api/auth";

  static const String register = "$apiAuth/register";
  static const String login = "$apiAuth/login";
  static const String getMe = "$baseUrl/api/auth/me"; 
}