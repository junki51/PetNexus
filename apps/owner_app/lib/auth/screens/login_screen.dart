import 'package:flutter/material.dart';
import 'package:owner_app/auth/widgets/custom_input_field.dart';
import '../controllers/login_controller.dart'; // ตรวจสอบรันพาธไฟล์ให้ถูกต้องนะครับ
import '../../layout/responsive_layout.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final LoginController _controller = LoginController();

  // Controllers สำหรับดึงค่าจากช่องกรอก
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  // กำหนดสีหลักตามภาพดีไซน์ image_ff48da.png
  final Color primaryTeal = const Color(0xFF319F9B); 
  final Color bgTopColor = const Color(0xFFCBE2E2);
  final Color bgBottomColor = const Color(0xFFFAF9F6);

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: bgTopColor,
      body: SafeArea(
        bottom: false,
        child: Column(
          children: [
            // --- Header: Back Button & Title ---
            Padding(
              padding: EdgeInsets.symmetric(
                horizontal: context.nw(16), 
                vertical: context.nh(8),
              ),
              child: Row(
                children: [
                  CircleAvatar(
                    backgroundColor: primaryTeal,
                    child: IconButton(
                      icon: const Icon(Icons.chevron_left, color: Colors.white),
                      onPressed: () => Navigator.maybePop(context),
                    ),
                  ),
                  SizedBox(width: context.nw(16)),
                  Text(
                    'เข้าสู่ระบบ',
                    style: TextStyle(
                      fontSize: context.nf(20),
                      fontWeight: FontWeight.bold,
                      color: const Color(0xFF2D3748),
                    ),
                  ),
                ],
              ),
            ),

            // --- Main Content Area ---
            Expanded(
              child: Container(
                width: double.infinity,
                margin: EdgeInsets.all(context.nh(20)),
                decoration: BoxDecoration(
                  color: bgBottomColor,
                  borderRadius: BorderRadius.only(
                    topLeft: Radius.elliptical(context.nw(250), context.nh(50)),
                    topRight: Radius.elliptical(context.nw(250), context.nh(50)),
                  ),
                ),
                child: SingleChildScrollView(
                  padding: EdgeInsets.all(context.nw(32)),
                  child: Column(
                    children: [
                      SizedBox(height: context.nh(20)),
                      
                      // โลโก้ PetNexus
                      RichText(
                        text: TextSpan(
                          style: TextStyle(
                            fontSize: context.nf(42),
                            fontWeight: FontWeight.bold,
                          ),
                          children: [
                            const TextSpan(
                              text: 'Pet',
                              style: TextStyle(color: Color(0xFF2D3748)),
                            ),
                            TextSpan(
                              text: 'Nexus',
                              style: TextStyle(color: primaryTeal),
                            ),
                          ],
                        ),
                      ),
                      SizedBox(height: context.nh(40)),

                      // Email Field
                      CustomInputField(
                        controller: _emailController,
                        hintText: 'กรอกอีเมล',
                        prefixIcon: Icons.email_outlined,
                        keyboardType: TextInputType.emailAddress,
                      ),
                      SizedBox(height: context.nh(20)),

                      // Password Field (ผูกเข้ากับคุณสมบัติจริงใน LoginController)
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) => CustomInputField(
                          controller: _passwordController,
                          hintText: 'รหัสผ่าน',
                          prefixIcon: Icons.lock_outline,
                          isPassword: true,
                          obscureText: !_controller.isPasswordVisible,
                          onToggleVisibility: _controller.togglePasswordVisibility,
                        ),
                      ),

                      // Link: ลืมรหัสผ่าน? (จัดชิดขวาตามภาพต้นฉบับ)
                      Align(
                        alignment: Alignment.centerRight,
                        child: TextButton(
                          onPressed: () {
                            // การจัดการหน้าลืมรหัสผ่าน
                          },
                          child: Text(
                            'ลืมรหัสผ่าน?',
                            style: TextStyle(
                              color: primaryTeal,
                              fontSize: context.nf(13),
                            ),
                          ),
                        ),
                      ),
                      SizedBox(height: context.nh(10)),

                      // Login Button (ผูกเข้ากับ AuthState และเมธอดของปุ่ม)
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) {
                          // ตรวจสอบสถานะการโหลดจาก Enum AuthState ของคุณเป๊ะๆ
                          final isLoading = _controller.state == AuthState.loading;
                          
                          return SizedBox(
                            width: double.infinity,
                            height: context.nh(56),
                            child: ElevatedButton(
                              onPressed: isLoading
                                ? null
                                : () {
                                    final emailInput = _emailController.text.trim();
                                    final passwordInput = _passwordController.text;

                                    if (emailInput.isEmpty || passwordInput.isEmpty) {
                                      ScaffoldMessenger.of(context).showSnackBar(
                                        const SnackBar(content: Text('กรุณากรอกอีเมลและรหัสผ่าน')),
                                      );
                                      return;
                                    }

                                    _controller.loginWithEmail(
                                      email: emailInput,
                                      password: passwordInput,
                                    );
                                  },
                                    
                              style: ElevatedButton.styleFrom(
                                backgroundColor: primaryTeal,
                                foregroundColor: Colors.white,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(context.nw(28)),
                                ),
                                elevation: 4,
                                shadowColor: primaryTeal.withValues(alpha: 0.4),
                              ),
                              child: isLoading
                                  ? const CircularProgressIndicator(color: Colors.white)
                                  : Row(
                                      mainAxisAlignment: MainAxisAlignment.center,
                                      children: [
                                        Icon(Icons.pets, size: context.nw(24)),
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
                          );
                        },
                      ),
                      SizedBox(height: context.nh(40)),

                      // เส้นแบ่งช่องทางโซเชียล
                      Text(
                        'หรือเข้าสู่ระบบด้วย',
                        style: TextStyle(
                          color: Colors.grey,
                          fontSize: context.nf(14),
                        ),
                      ),
                      SizedBox(height: context.nh(24)),

                      // Social Logins ด้านล่างสุด (ผูกฟังก์ชัน loginWithSocial)
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          _buildSocialButton(
                            icon: Icons.g_mobiledata,
                            color: Colors.red,
                            onTap: () => _controller.loginWithSocial('Google'),
                          ),
                          SizedBox(width: context.nw(24)),
                          _buildSocialButton(
                            icon: Icons.apple,
                            color: Colors.black,
                            onTap: () => _controller.loginWithSocial('Apple'),
                          ),
                          SizedBox(width: context.nw(24)),
                          _buildSocialButton(
                            icon: Icons.facebook,
                            color: Colors.blue,
                            onTap: () => _controller.loginWithSocial('Facebook'),
                          ),
                        ],
                      ),
                      SizedBox(height: context.nh(20)),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  // Helper สำหรับปุ่มทางเลือกไอคอนโซเชียลทรงกลม
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
            blurRadius: context.nw(10),
            offset: Offset(0, context.nh(4)),
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(50),
          onTap: onTap,
          child: Padding(
            padding: EdgeInsets.all(context.nw(12.0)),
            child: Icon(icon, color: color, size: context.nf(36)),
          ),
        ),
      ),
    );
  }
}