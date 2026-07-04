import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:owner_app/features/auth/controllers/auth_controller.dart';
import 'package:provider/provider.dart';

import '../../../app/app_routes.dart';
import '../../../core/constants.dart';
import '../../../layout/responsive_layout.dart';
import '../../../shared/widgets.dart';

class FirstScreen extends StatelessWidget {
  const FirstScreen({super.key});

  static const int _headerFlex = 3;
  static const int _contentFlex = 5;

  @override
  Widget build(BuildContext context) {
    final controller = context.watch<AuthController>();
    final isLoading = controller.state == AuthState.loading;

    return AppScaffold(
      scrollable: false,
      backgroundColor: AppColors.primaryLight,
      child: Column(
        children: [
          const Expanded(
            flex: _headerFlex,
            child: Center(child: AppLogo()),
          ),
          Expanded(
            flex: _contentFlex,
            child: _FirstScreenContent(isLoading: isLoading),
          ),
        ],
      ),
    );
  }
}

class _FirstScreenContent extends StatelessWidget {
  final bool isLoading;

  const _FirstScreenContent({required this.isLoading});

  static const double _curveWidth = 250;
  static const double _curveHeight = 60;
  static const double _horizontalPadding = 32;
  static const double _topPadding = 40;
  static const double _bottomPadding = 40;
  static const double _taglineSpacing = 40;
  static const double _buttonSpacing = 16;
  static const double _sectionSpacing = 40;
  static const double _socialSpacing = 55;
  static const double _buttonHeight = 56;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      padding: EdgeInsets.zero,
      color: AppColors.background,
      borderRadius: BorderRadius.only(
        topLeft: Radius.elliptical(
          context.nw(_curveWidth),
          context.nh(_curveHeight),
        ),
        topRight: Radius.elliptical(
          context.nw(_curveWidth),
          context.nh(_curveHeight),
        ),
      ),
      child: SingleChildScrollView(
        padding: EdgeInsets.fromLTRB(
          context.nw(_horizontalPadding),
          context.nh(_topPadding),
          context.nw(_horizontalPadding),
          context.nh(_bottomPadding),
        ),
        child: Column(
          children: [
            const _Tagline(),
            AppSpacing.h(context, _taglineSpacing),
            AppButton.primary(
              text: 'เข้าสู่ระบบ',
              icon: Icons.pets,
              loading: isLoading,
              height: context.nh(_buttonHeight),
              onPressed: () => Navigator.pushNamed(context, AppRoutes.login),
            ),
            AppSpacing.h(context, _buttonSpacing),
            AppButton.secondary(
              text: 'สร้างบัญชีใหม่',
              loading: isLoading,
              height: context.nh(_buttonHeight),
              onPressed: () => Navigator.pushNamed(context, AppRoutes.register),
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
          ],
        ),
      ),
    );
  }

  Future<void> _showComingSoon(BuildContext context) {
    return AppDialog.showMessage(
      context: context,
      title: 'Coming Soon',
      message: 'ฟีเจอร์นี้ยังไม่เปิดใช้งาน',
    );
  }
}

class _Tagline extends StatelessWidget {
  const _Tagline();

  static const double _fontSize = 20;
  static const double _lineHeight = 1.12;

  @override
  Widget build(BuildContext context) {
    final style = AppTextStyles.body(context).copyWith(
      fontSize: context.nf(_fontSize),
      fontWeight: FontWeight.w600,
      height: _lineHeight,
    );

    return Text.rich(
      TextSpan(
        text: 'Everything your pet needs,\n',
        style: style,
        children: [
          TextSpan(
            text: 'all in one place.',
            style: style.copyWith(color: AppColors.primary),
          ),
        ],
      ),
      textAlign: TextAlign.center,
    );
  }
}
