import 'package:flutter/material.dart';
import 'package:owner_app/auth/widgets/custom_input_field.dart';
import 'package:owner_app/layout/responsive_layout.dart';
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
                      onPressed: () => Navigator.pop(context),
                    ),
                  ),
                  SizedBox(width: context.nw(16)),
                  Text(
                    'สร้างบัญชีใหม่',
                    style: TextStyle(
                      fontSize: context.nf(20),
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
                margin: EdgeInsets.all(context.nh(20)),
                decoration: BoxDecoration(
                  color: bgBottomColor,
                  borderRadius: BorderRadius.only(
                    topLeft: Radius.elliptical(context.nw(250), context.nh(50)),
                    topRight: Radius.elliptical(
                      context.nw(250),
                      context.nh(50),
                    ),
                  ),
                ),
                child: SingleChildScrollView(
                  padding: EdgeInsets.all(context.nw(32)),
                  child: Column(
                    children: [
                      SizedBox(height: context.nh(20)),
                      // Logo
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
                        hintText: 'กรอกอีเมล*',
                        prefixIcon: Icons.email_outlined,
                      ),
                      SizedBox(height: context.nh(20)),

                      // Password Field (พร้อมปุ่มเปิด/ปิดตา)
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) => CustomInputField(
                          controller: _passwordController,
                          hintText: 'อย่างน้อย 8 ตัวอักษร*',
                          prefixIcon: Icons.lock_outline,
                          isPassword: true,
                          obscureText: !_controller.isPasswordVisible,
                          onToggleVisibility:
                              _controller.togglePasswordVisibility,
                        ),
                      ),
                      SizedBox(height: context.nh(20)),

                      // Confirm Password Field (แยก State ตาชัดเจน)
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) => CustomInputField(
                          controller: _confirmPasswordController,
                          hintText: 'ยืนยันรหัสผ่าน*',
                          prefixIcon: Icons.lock_clock_outlined,
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
                              fontSize: context.nf(13),
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
                            Expanded(
                              child: Text(
                                'ฉันยอมรับเงื่อนไขในการใช้งานและนโยบายความเป็นส่วนตัว',
                                style: TextStyle(
                                  fontSize: context.nf(12),
                                  color: Colors.blueGrey,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      SizedBox(height: context.nh(30)),

                      // Register Button
                      ListenableBuilder(
                        listenable: _controller,
                        builder: (context, _) {
                          final isLoading =
                              _controller.state == RegisterState.loading;
                          return SizedBox(
                            width: double.infinity,
                            height: context.nh(56),
                            child: ElevatedButton(
                              onPressed:
                                  (isLoading || !_controller.isAcceptedTerms)
                                  ? null
                                  : () {
                                      final emailInput = _emailController.text.trim();
                                      final passwordInput =
                                          _passwordController.text;
                                      final confirmPasswordInput =
                                          _confirmPasswordController.text;
                                      _controller.register(
                                        emailInput,
                                        passwordInput,
                                        confirmPasswordInput,
                                      );
                                    },
                              style: ElevatedButton.styleFrom(
                                backgroundColor: primaryTeal,
                                foregroundColor: Colors.white,
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(
                                    context.nw(28),
                                  ),
                                ),
                                elevation: 4,
                                shadowColor: primaryTeal.withValues(alpha: 0.4),
                              ),
                              child: isLoading
                                  ? const CircularProgressIndicator(
                                      color: Colors.white,
                                    )
                                  : Row(
                                      mainAxisAlignment:
                                          MainAxisAlignment.center,
                                      children: [
                                        Icon(Icons.pets, size: context.nw(24)),
                                        SizedBox(width: context.nw(12)),
                                        Text(
                                          'สร้างบัญชีใหม่',
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
}
