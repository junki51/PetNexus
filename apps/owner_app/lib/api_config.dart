class ApiConfig {
  // 10.0.2.2 คือ IP พิเศษสำหรับ Android Emulator ที่วิ่งไปหา Localhost ของเครื่องคอมฯ นั้นๆ
  static const String baseUrl = "http://10.0.2.2:8080"; 
  
  static const String login = "$baseUrl/login";
  static const String register = "$baseUrl/register";
  static const String contact = "$baseUrl/contact"; // ตัวอย่างถ้ามี API อื่นเพิ่ม
}