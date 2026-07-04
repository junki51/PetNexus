import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../layout/responsive_layout.dart';

class AppScaffold extends StatelessWidget {
  final Widget child;

  final PreferredSizeWidget? appBar;

  final Widget? floatingActionButton;

  final Color? backgroundColor;

  final bool scrollable;

  const AppScaffold({
    super.key,
    required this.child,
    this.appBar,
    this.floatingActionButton,
    this.backgroundColor,
    this.scrollable = true,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor:
          backgroundColor ??
          AppColors.background,

      appBar: appBar,

      floatingActionButton:
          floatingActionButton,

      body: SafeArea(
        child: scrollable
            ? SingleChildScrollView(
                padding: EdgeInsets.all(
                  context.nw(20),
                ),
                child: child,
              )
            : child,
      ),
    );
  }
}