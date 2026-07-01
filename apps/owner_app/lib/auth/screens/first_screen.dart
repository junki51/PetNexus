import 'package:flutter/material.dart';
import 'package:owner_app/auth/screens/login_screen.dart';
import 'package:owner_app/auth/screens/register_screen.dart';
import '../controllers/login_controller.dart';
import '../../layout/responsive_layout.dart';

// นำเข้า Extension (แก้ไข path ให้ตรงกับโปรเจกต์ของคุณ)
// import 'path_to_your_extension/responsive_context.dart'; 

class FirstScreen extends StatefulWidget {
  const FirstScreen({super.key});

  @override
  State<FirstScreen> createState() => _FirstScreenState();
}

class _FirstScreenState extends State<FirstScreen> {
  // สร้าง Instance ของ Controller
  final LoginController _controller = LoginController();

  // กำหนดสีตาม Design System
  final Color primaryTeal = const Color(0xFF38A3A5);
  final Color bgTopColor = const Color(0xFFCBE2E2);
  final Color bgBottomColor = const Color(0xFFFAF9F6);
  final Color textDark = const Color(0xFF2D3748);

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: bgTopColor,
      body: SafeArea(
        bottom: false, // ปล่อยขอบล่างให้เต็มจอ
        child: Column(
          children: [
            // --- Section 1: Logo (Top) ---
            Expanded(
              flex: 3,
              child: Center(
                child: RichText(
                  text: TextSpan(
                    style: TextStyle(
                      fontSize: context.nf(42), // ปรับขนาดฟอนต์ Logo
                      fontWeight: FontWeight.bold,
                      letterSpacing: -1.0,
                    ),
                    children: [
                      TextSpan(text: 'Pet', style: TextStyle(color: textDark)),
                      TextSpan(text: 'Nexus', style: TextStyle(color: primaryTeal)),
                    ],
                  ),
                ),
              ),
            ),

            // --- Section 2: Content & Actions (Bottom Curved) ---
            Expanded(
              flex: 5,
              child: Container(
                width: double.infinity,
                decoration: BoxDecoration(
                  color: bgBottomColor,
                  borderRadius: BorderRadius.only(
                    // ทำเส้นโค้งด้านบนแบบ Responsive
                    topLeft: Radius.elliptical(context.nw(250), context.nh(60)),
                    topRight: Radius.elliptical(context.nw(250), context.nh(60)),
                  ),
                ),
                child: Padding(
                  padding: EdgeInsets.symmetric(horizontal: context.nw(32.0)), // ปรับความห่างซ้าย-ขวา
                  child: SingleChildScrollView(
                    child: Column(
                      children: [
                        SizedBox(height: context.nh(40)), // ปรับระยะห่างแนวตั้ง

                        // Subtitle
                        Text(
                          'Everything your pet needs,',
                          style: TextStyle(
                            fontSize: context.nf(16), // ปรับฟอนต์ข้อความอธิบาย
                            fontWeight: FontWeight.w600,
                            color: textDark,
                          ),
                        ),
                        Text(
                          'all in one place.',
                          style: TextStyle(
                            fontSize: context.nf(16),
                            fontWeight: FontWeight.w600,
                            color: primaryTeal,
                          ),
                        ),
                        SizedBox(height: context.nh(40)),

                        // ใช้ ListenableBuilder เพื่ออัปเดต UI เฉพาะจุดที่ State เปลี่ยน
                        ListenableBuilder(
                          listenable: _controller,
                          builder: (context, _) {
                            final isLoading = _controller.state == AuthState.loading;

                            return Column(
                              children: [
                                // Login Button
                                SizedBox(
                                  width: double.infinity,
                                  height: context.nh(56), // ปรับความสูงปุ่ม
                                  child: ElevatedButton(
                                    onPressed: isLoading ? null : () {
                                      Navigator.push(context, MaterialPageRoute(builder: (context) => const LoginScreen()));
                                    },
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor: primaryTeal,
                                      foregroundColor: Colors.white,
                                      elevation: 0,
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(context.nh(28)), // ล้อตามความสูงปุ่มที่เปลี่ยนไป
                                      ),
                                    ),
                                    child: isLoading
                                        ? SizedBox(
                                            width: context.nw(24),
                                            height: context.nw(24),
                                            child: const CircularProgressIndicator(
                                              color: Colors.white,
                                              strokeWidth: 3,
                                            ),
                                          )
                                        : Row(
                                            mainAxisAlignment: MainAxisAlignment.center,
                                            children: [
                                              Icon(Icons.pets, size: context.nf(24)), // ปรับขนาดไอคอน
                                              SizedBox(width: context.nw(12)),
                                              Text(
                                                'เข้าสู่ระบบ',
                                                style: TextStyle(
                                                  fontSize: context.nf(18),
                                                  fontWeight: FontWeight.bold,
                                                ),
                                              ),
                                            ],
                                          ),
                                  ),
                                ),
                                SizedBox(height: context.nh(16)),

                                // Create Account Button
                                SizedBox(
                                  width: double.infinity,
                                  height: context.nh(56),
                                  child: OutlinedButton(
                                    onPressed: isLoading ? null : () {
                                      Navigator.push(context, MaterialPageRoute(builder: (context) => const RegisterScreen()));
                                    },
                                    style: OutlinedButton.styleFrom(
                                      backgroundColor: Colors.white,
                                      foregroundColor: textDark,
                                      side: const BorderSide(color: Colors.black12),
                                      elevation: 0,
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(context.nh(28)),
                                      ),
                                    ),
                                    child: Text(
                                      'สร้างบัญชีใหม่',
                                      style: TextStyle(
                                        fontSize: context.nf(18),
                                        fontWeight: FontWeight.w600,
                                      ),
                                    ),
                                  ),
                                ),
                              ],
                            );
                          },
                        ),
                        SizedBox(height: context.nh(40)),

                        // Divider
                        Text(
                          'หรือเข้าสู่ระบบด้วย',
                          style: TextStyle(
                            color: Colors.grey,
                            fontSize: context.nf(14),
                          ),
                        ),
                        SizedBox(height: context.nh(24)),

                        // Social Login Buttons
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            _buildSocialButton(
                              context: context, // ส่ง context เข้าไปใช้งาน
                              icon: Icons.g_mobiledata, 
                              color: Colors.red,
                              onTap: () => _controller.loginWithSocial('Google'),
                            ),
                            SizedBox(width: context.nw(24)),
                            _buildSocialButton(
                              context: context,
                              icon: Icons.apple,
                              color: Colors.black,
                              onTap: () => _controller.loginWithSocial('Apple'),
                            ),
                            SizedBox(width: context.nw(24)),
                            _buildSocialButton(
                              context: context,
                              icon: Icons.facebook,
                              color: Colors.blue,
                              onTap: () => _controller.loginWithSocial('Facebook'),
                            ),
                          ],
                        ),
                        SizedBox(height: context.nh(40)),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  // Widget ย่อยสำหรับปุ่ม Social เพิ่มพารามิเตอร์ BuildContext context เข้าไป
  Widget _buildSocialButton({
    required BuildContext context,
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: context.nw(10), // ปรับระยะเบลอตามสัดส่วนจอ
            offset: Offset(0, context.nh(4)), // ปรับตำแหน่งเงาตามความสูง
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(50),
          onTap: onTap,
          child: Padding(
            padding: EdgeInsets.all(context.nw(12.0)), // ปรับ Padding ด้านในปุ่มวงกลม
            child: Icon(icon, color: color, size: context.nf(36)), // ปรับขนาดไอคอนโซเชียล
          ),
        ),
      ),
    );
  }
}