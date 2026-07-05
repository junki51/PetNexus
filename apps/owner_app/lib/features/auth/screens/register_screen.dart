import 'package:flutter/material.dart';
import 'package:owner_app/features/auth/controllers/auth_controller.dart';
import 'package:owner_app/shared/widgets/app_button.dart';
import 'package:provider/provider.dart';

import '../../../app/app_routes.dart';
import '../../../core/constants.dart';
import '../../../layout/responsive_layout.dart';
import '../widgets/auth_screen_layout.dart';
import '../widgets/custom_input_field.dart';

class RegisterScreen extends StatefulWidget {
  const RegisterScreen({super.key});

  @override
  State<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends State<RegisterScreen> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _confirmPasswordController =
      TextEditingController();

  static const double _fieldSpacing = 22;
  static const double _buttonSpacing = 28;
  static const double _termsTopPadding = 2;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AuthController>();
    final isLoading = controller.state == AuthState.loading;

    return AuthScreenLayout(
      title: 'สร้างบัญชีใหม่',
      onBack: () => Navigator.pop(context),
      children: [
        CustomInputField(
          controller: _emailController,
          hintText: 'กรอกอีเมล*',
          prefixIcon: Icons.email_outlined,
          keyboardType: TextInputType.emailAddress,
          textInputAction: TextInputAction.next,
        ),
        AppSpacing.h(context, _fieldSpacing),
        CustomInputField(
          controller: _passwordController,
          hintText: 'อย่างน้อย 8 ตัวอักษร*',
          prefixIcon: Icons.lock_outline,
          isPassword: true,
          obscureText: !controller.isPasswordVisible,
          onToggleVisibility: controller.togglePasswordVisibility,
          textInputAction: TextInputAction.next,
        ),
        AppSpacing.h(context, _fieldSpacing),
        CustomInputField(
          controller: _confirmPasswordController,
          hintText: 'ยืนยันรหัสผ่าน*',
          isPassword: true,
          obscureText: !controller.isConfirmPasswordVisible,
          prefixIcon: Icons.lock_clock_outlined,
          onToggleVisibility: controller.toggleConfirmPasswordVisibility,
          textInputAction: TextInputAction.done,
        ),
        AppSpacing.h(context, _fieldSpacing),
        Align(
          alignment: Alignment.centerRight,
          child: TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(
              'มีบัญชีอยู่แล้ว?',
              style: AppTextStyles.caption(context).copyWith(
                color: AppColors.primary,
              ),
            ),
          ),
        ),
        AppSpacing.h(context, _fieldSpacing),
        Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SizedBox(
              width: context.nw(28),
              height: context.nh(28),
              child: Checkbox(
                value: controller.acceptedTerms,
                onChanged: controller.toggleAcceptedTerms,
                activeColor: AppColors.primary,
                materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                visualDensity: VisualDensity.compact,
              ),
            ),
            Expanded(
              child: Padding(
                padding: EdgeInsets.only(top: context.nh(_termsTopPadding)),
                child: Text.rich(
                  TextSpan(
                    text: 'ฉันยอมรับ',
                    style: AppTextStyles.caption(context),
                    children: [
                      TextSpan(
                        text: 'เงื่อนไขการใช้งาน',
                        style: AppTextStyles.caption(context).copyWith(
                          color: AppColors.primary,
                        ),
                      ),
                      const TextSpan(text: ' และ'),
                      TextSpan(
                        text: 'นโยบายความเป็นส่วนตัว',
                        style: AppTextStyles.caption(context).copyWith(
                          color: AppColors.primary,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
        AppSpacing.h(context, _buttonSpacing),
        AppButton.primary(
          text: 'สร้างบัญชีใหม่',
          icon: Icons.pets,
          loading: isLoading,
          onPressed: () => _register(context, controller),
        ),
      ],
    );
  }

  Future<void> _register(
    BuildContext context,
    AuthController controller,
  ) async {
    final email = _emailController.text.trim();
    final password = _passwordController.text;
    final confirmPassword = _confirmPasswordController.text;

    if (email.isEmpty || password.isEmpty || confirmPassword.isEmpty) {
      _showSnackBar(context, 'กรุณากรอกข้อมูลให้ครบ');
      return;
    }

    if (!controller.acceptedTerms) {
      _showSnackBar(context, 'กรุณายอมรับเงื่อนไขการใช้งาน');
      return;
    }

    final navigator = Navigator.of(context);
    final messenger = ScaffoldMessenger.of(context);

    final success = await controller.register(
      email: email,
      password: password,
      confirmPassword: confirmPassword,
    );

    if (!mounted) return;

    if (success) {
      navigator.pushReplacementNamed(AppRoutes.completeProfile);
      return;
    }

    messenger.showSnackBar(
      SnackBar(content: Text(controller.errorMessage ?? 'Register Failed')),
    );
  }

  void _showSnackBar(BuildContext context, String message) {
    ScaffoldMessenger.of(
      context,
    ).showSnackBar(SnackBar(content: Text(message)));
  }
}
