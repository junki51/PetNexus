import 'package:flutter/material.dart';
import 'package:owner_app/auth/screens/register_screen.dart';
import '../controllers/login_controller.dart';

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
                    style: const TextStyle(
                      fontSize: 42,
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
                  borderRadius: const BorderRadius.only(
                    // ทำเส้นโค้งด้านบน
                    topLeft: Radius.elliptical(250, 60),
                    topRight: Radius.elliptical(250, 60),
                  ),
                ),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 32.0),
                  child: SingleChildScrollView(
                    child: Column(
                      children: [
                        const SizedBox(height: 40),
                        
                        // Subtitle
                        Text(
                          'Everything your pet needs,',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                            color: textDark,
                          ),
                        ),
                        Text(
                          'all in one place.',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                            color: primaryTeal,
                          ),
                        ),
                        const SizedBox(height: 40),

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
                                  height: 56,
                                  child: ElevatedButton(
                                    onPressed: isLoading ? null : _controller.loginWithEmail,
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor: primaryTeal,
                                      foregroundColor: Colors.white,
                                      elevation: 0,
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(28),
                                      ),
                                    ),
                                    child: isLoading
                                        ? const SizedBox(
                                            width: 24,
                                            height: 24,
                                            child: CircularProgressIndicator(
                                              color: Colors.white,
                                              strokeWidth: 3,
                                            ),
                                          )
                                        : const Row(
                                            mainAxisAlignment: MainAxisAlignment.center,
                                            children: [
                                              Icon(Icons.pets, size: 24),
                                              SizedBox(width: 12),
                                              Text(
                                                'เข้าสู่ระบบ',
                                                style: TextStyle(
                                                  fontSize: 18,
                                                  fontWeight: FontWeight.bold,
                                                ),
                                              ),
                                            ],
                                          ),
                                  ),
                                ),
                                const SizedBox(height: 16),

                                // Create Account Button
                                SizedBox(
                                  width: double.infinity,
                                  height: 56,
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
                                        borderRadius: BorderRadius.circular(28),
                                      ),
                                    ),
                                    child: const Text(
                                      'สร้างบัญชีใหม่',
                                      style: TextStyle(
                                        fontSize: 18,
                                        fontWeight: FontWeight.w600,
                                      ),
                                    ),
                                  ),
                                ),
                              ],
                            );
                          },
                        ),
                        const SizedBox(height: 40),

                        // Divider
                        const Text(
                          'หรือเข้าสู่ระบบด้วย',
                          style: TextStyle(
                            color: Colors.grey,
                            fontSize: 14,
                          ),
                        ),
                        const SizedBox(height: 24),

                        // Social Login Buttons
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            _buildSocialButton(
                              icon: Icons.g_mobiledata, // เปลี่ยนเป็นโลโก้จริงด้วย flutter_svg ได้
                              color: Colors.red,
                              onTap: () => _controller.loginWithSocial('Google'),
                            ),
                            const SizedBox(width: 24),
                            _buildSocialButton(
                              icon: Icons.apple,
                              color: Colors.black,
                              onTap: () => _controller.loginWithSocial('Apple'),
                            ),
                            const SizedBox(width: 24),
                            _buildSocialButton(
                              icon: Icons.facebook,
                              color: Colors.blue,
                              onTap: () => _controller.loginWithSocial('Facebook'),
                            ),
                          ],
                        ),
                        const SizedBox(height: 40),
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

  // Widget ย่อยสำหรับปุ่ม Social เพื่อลดความซ้ำซ้อนของโค้ด
  Widget _buildSocialButton({
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
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(50),
          onTap: onTap,
          child: Padding(
            padding: const EdgeInsets.all(12.0),
            child: Icon(icon, color: color, size: 36),
          ),
        ),
      ),
    );
  }
}