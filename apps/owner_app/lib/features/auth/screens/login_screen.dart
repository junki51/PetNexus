import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:owner_app/features/auth/controllers/auth_controller.dart';
import 'package:owner_app/shared/widgets/app_button.dart';
import 'package:owner_app/shared/widgets/app_social_button.dart';
import 'package:provider/provider.dart';

import '../../../app/app_routes.dart';
import '../../../core/constants.dart';
import '../widgets/auth_screen_layout.dart';
import '../widgets/custom_input_field.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  static const double _fieldSpacing = 16;
  static const double _forgotSpacing = 6;
  static const double _sectionSpacing = 28;
  static const double _socialSpacing = 55;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AuthController>();
    final isLoading = controller.state == AuthState.loading;

    return AuthScreenLayout(
      title: 'เข้าสู่ระบบ',
      children: [
        CustomInputField(
          controller: _emailController,
          hintText: 'กรอกอีเมล',
          prefixIcon: Icons.email_outlined,
          keyboardType: TextInputType.emailAddress,
          textInputAction: TextInputAction.next,
        ),
        AppSpacing.h(context, _fieldSpacing),
        CustomInputField(
          controller: _passwordController,
          hintText: 'รหัสผ่าน',
          prefixIcon: Icons.lock_outline,
          isPassword: true,
          obscureText: !controller.isPasswordVisible,
          onToggleVisibility: controller.togglePasswordVisibility,
          textInputAction: TextInputAction.done,
        ),
        Align(
          alignment: Alignment.centerRight,
          child: TextButton(
            onPressed: () {},
            child: Text(
              'ลืมรหัสผ่าน?',
              style: AppTextStyles.caption(
                context,
              ).copyWith(color: AppColors.primary),
            ),
          ),
        ),
        AppSpacing.h(context, _forgotSpacing),
        AppButton.primary(
          text: 'เข้าสู่ระบบ',
          icon: Icons.pets,
          loading: isLoading,
          onPressed: () => _login(context, controller),
        ),
        AppSpacing.h(context, _sectionSpacing),
        Text('หรือเข้าสู่ระบบด้วย', style: AppTextStyles.caption(context)),
        AppSpacing.h(context, _socialSpacing),
        Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                AppSocialButton(
                  icon: FontAwesomeIcons.google,
                  color: AppColors.google,
                  onTap: () => _showComingSoon(context),
                ),
                AppSpacing.w(context, _socialSpacing),
                AppSocialButton(
                  icon: Icons.apple,
                  color: AppColors.apple,
                  onTap: () => _showComingSoon(context),
                ),
                AppSpacing.w(context, _socialSpacing),
                AppSocialButton(
                  icon: Icons.facebook,
                  color: AppColors.facebook,
                  onTap: () => _showComingSoon(context),
                ),
              ],
            ),
        AppSpacing.h(context, _fieldSpacing),
      ],
    );
  }

  Future<void> _login(BuildContext context, AuthController controller) async {
    final email = _emailController.text.trim();
    final password = _passwordController.text;

    if (email.isEmpty || password.isEmpty) {
      _showSnackBar(context, 'กรุณากรอกอีเมลและรหัสผ่าน');
      return;
    }

    final navigator = Navigator.of(context);
    final messenger = ScaffoldMessenger.of(context);

    final success = await controller.login(email: email, password: password);

    if (!mounted) return;

    if (success) {
      navigator.pushReplacementNamed(AppRoutes.home);
      return;
    }

    messenger.showSnackBar(
      SnackBar(content: Text(controller.errorMessage ?? 'Login Failed')),
    );
  }

  void _showComingSoon(BuildContext context) {
    _showSnackBar(context, 'Coming Soon');
  }

  void _showSnackBar(BuildContext context, String message) {
    ScaffoldMessenger.of(
      context,
    ).showSnackBar(SnackBar(content: Text(message)));
  }
}
