import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

import '../../../core/constants.dart';
import '../../../layout/responsive_layout.dart';
import '../../../shared/widgets.dart';

class AuthScreenLayout extends StatelessWidget {
  final String title;

  final List<Widget> children;

  final VoidCallback? onBack;

  final EdgeInsets? cardMargin;

  final EdgeInsets? cardPadding;

  final Color? cardColor;

  final BorderRadiusGeometry? cardBorderRadius;

  const AuthScreenLayout({
    super.key,
    required this.title,
    required this.children,
    this.onBack,
    this.cardMargin,
    this.cardPadding,
    this.cardColor,
    this.cardBorderRadius,
  });

  static const double _headerHeight = 108;
  static const double _headerHorizontalPadding = 16;
  static const double _headerVerticalPadding = 8;
  static const double _curveWidth = 250;
  static const double _curveHeight = 52;
  static const double _horizontalPadding = 20;
  static const double _topPadding = 48;
  static const double _bottomPadding = 18;
  static const double _logoBottomSpacing = 34;

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      scrollable: false,
      backgroundColor: AppColors.primaryLight,
      child: LayoutBuilder(
        builder: (context, constraints) {
          final headerHeight = context
              .nh(_headerHeight)
              .clamp(92.0, 124.0)
              .toDouble();

          return Column(
            children: [
              SizedBox(
                height: headerHeight,
                child: _AuthHeroHeader(
                  title: title,
                  onBack: onBack,
                  horizontalPadding: _headerHorizontalPadding,
                  verticalPadding: _headerVerticalPadding,
                ),
              ),
              Expanded(
                child: Padding(
                  padding: cardMargin ?? EdgeInsets.zero,
                  child: AppCard(
                    padding: EdgeInsets.zero,
                    color: cardColor ?? AppColors.background,
                    borderRadius:
                        cardBorderRadius ??
                        BorderRadius.only(
                          topLeft: Radius.elliptical(
                            context.nw(_curveWidth),
                            context.nh(_curveHeight),
                          ),
                          topRight: Radius.elliptical(
                            context.nw(_curveWidth),
                            context.nh(_curveHeight),
                          ),
                        ),
                    child: Padding(
                      padding:
                          cardPadding ??
                          EdgeInsets.fromLTRB(
                            context.nw(_horizontalPadding),
                            context.nh(_topPadding),
                            context.nw(_horizontalPadding),
                            context.nh(_bottomPadding),
                          ),
                      child: Column(
                        children: [
                          const AppLogo(),
                          SizedBox(height: context.nh(_logoBottomSpacing)),
                          ...children,
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}

class _AuthHeroHeader extends StatelessWidget {
  final String title;

  final VoidCallback? onBack;

  final double horizontalPadding;

  final double verticalPadding;

  const _AuthHeroHeader({
    required this.title,
    required this.onBack,
    required this.horizontalPadding,
    required this.verticalPadding,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Positioned(
          left: context.nw(horizontalPadding),
          top: context.nh(verticalPadding +12),
          child: Row(
            children: [
              Semantics(
                label: title,
                button: true,
                child: CircleAvatar(
                  radius: context.nw(25),
                  backgroundColor: AppColors.primary,
                  child: IconButton(
                    padding: EdgeInsets.zero,
                    icon: Icon(
                      Icons.chevron_left,
                      color: AppColors.background,
                      size: context.icon(30),
                    ),
                    onPressed: onBack ?? () => Navigator.maybePop(context),
                  ),
                ),
              ),
              SizedBox(width: context.nw(16)),
              Text(
                title,
                style: AppTextStyles.caption(context).copyWith(
                  color: AppColors.textPrimary,
                  fontWeight: FontWeight.w600,
                  fontSize: context.nf(20),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class AuthSocialLoginRow extends StatelessWidget {
  final VoidCallback onGoogleTap;

  final VoidCallback onAppleTap;

  final VoidCallback onFacebookTap;

  const AuthSocialLoginRow({
    super.key,
    required this.onGoogleTap,
    required this.onAppleTap,
    required this.onFacebookTap,
  });

  static const double _spacing = 28;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        AppSocialButton(
          icon: CupertinoIcons.globe,
          color: AppColors.google,
          onTap: onGoogleTap,
        ),
        AppSpacing.w(context, _spacing),
        AppSocialButton(
          icon: Icons.apple,
          color: AppColors.apple,
          onTap: onAppleTap,
        ),
        AppSpacing.w(context, _spacing),
        AppSocialButton(
          icon: Icons.facebook,
          color: AppColors.facebook,
          onTap: onFacebookTap,
        ),
      ],
    );
  }
}

class AuthPrimaryButton extends StatelessWidget {
  final String text;

  final bool loading;

  final VoidCallback? onPressed;

  const AuthPrimaryButton({
    super.key,
    required this.text,
    this.loading = false,
    this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    final height = context.nh(48).clamp(46.0, 52.0).toDouble();
    final radius = BorderRadius.circular(height / 2);

    return SizedBox(
      width: double.infinity,
      height: height,
      child: ElevatedButton(
        onPressed: loading ? null : onPressed,
        style: ElevatedButton.styleFrom(
          elevation: 0,
          backgroundColor: AppColors.primary,
          foregroundColor: AppColors.background,
          disabledBackgroundColor: AppColors.primary.withValues(alpha: 0.55),
          disabledForegroundColor: AppColors.background,
          shape: RoundedRectangleBorder(borderRadius: radius),
          padding: EdgeInsets.symmetric(horizontal: context.nw(18)),
        ),
        child: loading
            ? SizedBox(
                width: context.nw(22),
                height: context.nw(22),
                child: const CircularProgressIndicator(
                  strokeWidth: 2,
                  color: AppColors.background,
                ),
              )
            : Row(
                children: [
                  Icon(Icons.pets, size: context.icon(50)),
                  SizedBox(width: context.nw(8)),
                  Expanded(
                    child: Text(
                      text,
                      textAlign: TextAlign.center,
                      style: AppTextStyles.button(context).copyWith(
                        color: AppColors.background,
                        fontSize: context.nf(35),
                      ),
                    ),
                  ),
                  SizedBox(width: context.icon(50) + context.nw(8)),
                ],
              ),
      ),
    );
  }
}
