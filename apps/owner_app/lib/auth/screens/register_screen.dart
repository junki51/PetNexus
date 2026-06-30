import 'package:flutter/material.dart';
import '../controllers/register_controller.dart';

class RegisterScreen extends StatefulWidget {
  const RegisterScreen({super.key});

  @override
  State<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends State<RegisterScreen> {
  final RegisterController _controller = RegisterController();

  // Controllers สำหรับดึงค่าจากช่องกรอก
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _confirmPasswordController =
      TextEditingController();

  final Color primaryTeal = const Color(0xFF38A3A5);
  final Color bgTopColor = const Color(0xFFCBE2E2);
  final Color bgBottomColor = const Color(0xFFFAF9F6);

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
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
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: Row(
                children: [
                  CircleAvatar(
                    backgroundColor: primaryTeal,
                    child: IconButton(
                      icon: const Icon(Icons.chevron_left, color: Colors.white),
                      onPressed: () => Navigator.pop(context),
                    ),
                  ),
                  const SizedBox(width: 16),
                  const Text(
                    'สร้างบัญชีใหม่',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: Color(0xFF2D3748),
                    ),
                  ),
                ],
              ),
            ),

            // --- Main Content Area ---
            Expanded(
              child: Container(
                width: double.infinity,
                margin: const EdgeInsets.only(top: 20),
                decoration: BoxDecoration(
                  color: bgBottomColor,
                  borderRadius: const BorderRadius.only(
                    topLeft: Radius.elliptical(250, 50),
                    topRight: Radius.elliptical(250, 50),
                  ),
                ),
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(32.0),
                  child: Column(
                    children: [
                      const SizedBox(height: 20),
                      // Logo
                      RichText(
                        text: TextSpan(
                          style: const TextStyle(
                            fontSize: 42,
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
                      const SizedBox(height: 40),

                      // Email Field
                      _buildTextField(
                        controller: _emailController,
                        hint: 'กรอกอีเมล*',
                        icon: Icons.email_outlined,
                      ),
                      const SizedBox(height: 20),

                      // Password Field (พร้อมปุ่มเปิด/ปิดตา)
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) => _buildTextField(
                          controller: _passwordController,
                          hint: 'อย่างน้อย 8 ตัวอักษร*',
                          icon: Icons.lock_outline,
                          isPassword: true,
                          obscureText: !_controller.isPasswordVisible,
                          onToggleVisibility:
                              _controller.togglePasswordVisibility,
                        ),
                      ),
                      const SizedBox(height: 20),

                      // Confirm Password Field (แยก State ตาชัดเจน)
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) => _buildTextField(
                          controller: _confirmPasswordController,
                          hint: 'ยืนยันรหัสผ่าน*',
                          icon: Icons.lock_clock_outlined,
                          isPassword: true,
                          obscureText: !_controller.isConfirmPasswordVisible,
                          onToggleVisibility:
                              _controller.toggleConfirmPasswordVisibility,
                        ),
                      ),

                      // Link: มีบัญชีอยู่แล้ว?
                      Align(
                        alignment: Alignment.centerRight,
                        child: TextButton(
                          onPressed: () => Navigator.pop(context),
                          child: Text(
                            'มีบัญชีอยู่แล้ว?',
                            style: TextStyle(
                              color: Colors.blue.shade700,
                              fontSize: 13,
                            ),
                          ),
                        ),
                      ),

                      // Terms & Conditions Checkbox
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) => Row(
                          children: [
                            Checkbox(
                              value: _controller.isAcceptedTerms,
                              onChanged: _controller.toggleAcceptedTerms,
                              activeColor: primaryTeal,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(4),
                              ),
                            ),
                            const Expanded(
                              child: Text(
                                'ฉันยอมรับเงื่อนไขในการใช้งานและนโยบายความเป็นส่วนตัว',
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.blueGrey,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 30),

                      // Register Button
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) {
                          final isLoading =
                              _controller.state == RegisterState.loading;
                          return SizedBox(
                            width: double.infinity,
                            height: 56,
                            child: ElevatedButton(
                              onPressed:
                                  (isLoading || !_controller.isAcceptedTerms)
                                  ? null
                                  : () => _controller.register(
                                      _emailController.text,
                                      _passwordController.text,
                                      _confirmPasswordController.text,
                                    ),
                              style: ElevatedButton.styleFrom(
                                backgroundColor: primaryTeal,
                                foregroundColor: Colors.white,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(28),
                                ),
                                elevation: 4,
                                shadowColor: primaryTeal.withValues(alpha: 0.4),
                              ),
                              child: isLoading
                                  ? const CircularProgressIndicator(
                                      color: Colors.white,
                                    )
                                  : const Row(
                                      mainAxisAlignment:
                                          MainAxisAlignment.center,
                                      children: [
                                        Icon(Icons.pets, size: 24),
                                        SizedBox(width: 12),
                                        Text(
                                          'สร้างบัญชีใหม่',
                                          style: TextStyle(
                                            fontSize: 18,
                                            fontWeight: FontWeight.bold,
                                          ),
                                        ),
                                      ],
                                    ),
                            ),
                          );
                        },
                      ),
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

  // Helper สำหรับสร้าง TextField ที่มีเงาและดีไซน์ตามรูป
  Widget _buildTextField({
    required TextEditingController controller,
    required String hint,
    required IconData icon,
    bool isPassword = false,
    bool obscureText = false,
    VoidCallback? onToggleVisibility,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(30),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.08),
            blurRadius: 15,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: TextField(
        controller: controller,
        obscureText: obscureText,
        decoration: InputDecoration(
          hintText: hint,
          hintStyle: const TextStyle(color: Colors.black38),
          prefixIcon: Icon(icon, color: Colors.black54),
          suffixIcon: isPassword
              ? IconButton(
                  icon: Icon(
                    obscureText
                        ? Icons.visibility_off_outlined
                        : Icons.visibility_outlined,
                    color: Colors.black45,
                  ),
                  onPressed: onToggleVisibility,
                )
              : null,
          border: InputBorder.none,
          contentPadding: const EdgeInsets.symmetric(
            horizontal: 20,
            vertical: 16,
          ),
        ),
      ),
    );
  }
}
