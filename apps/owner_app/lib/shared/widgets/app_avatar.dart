import 'dart:io';

import 'package:flutter/material.dart';

import '../../core/constants/app_colors.dart';
import '../../core/constants/app_text_styles.dart';
import '../../layout/responsive_layout.dart';

class AppAvatar extends StatelessWidget {
  final File? imageFile;

  final String? imageUrl;

  final VoidCallback onTap;

  final String buttonText;

  final double? radius;

  const AppAvatar({
    super.key,
    this.imageFile,
    this.imageUrl,
    required this.onTap,
    this.buttonText = "Upload Image",
    this.radius,
  });

  @override
  Widget build(BuildContext context) {
    final avatarRadius = radius ?? context.nw(55);

    ImageProvider? provider;

    if (imageFile != null) {
      provider = FileImage(imageFile!);
    } else if (imageUrl != null && imageUrl!.isNotEmpty) {
      provider = NetworkImage(imageUrl!);
    }

    return Column(
      children: [
        GestureDetector(
          onTap: onTap,
          child: Stack(
            alignment: Alignment.bottomRight,
            children: [
              CircleAvatar(
                radius: avatarRadius,
                backgroundColor: Colors.grey.shade200,
                backgroundImage: provider,
                child: provider == null
                    ? Icon(
                        Icons.person,
                        size: context.icon(55),
                        color: Colors.grey,
                      )
                    : null,
              ),

              Container(
                padding: EdgeInsets.all(
                  context.nw(6),
                ),
                decoration: BoxDecoration(
                  color: AppColors.primary,
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.camera_alt,
                  color: Colors.white,
                  size: context.icon(18),
                ),
              ),
            ],
          ),
        ),

        SizedBox(
          height: context.nh(12),
        ),

        TextButton.icon(
          onPressed: onTap,
          icon: Icon(
            Icons.upload,
            size: context.icon(18),
          ),
          label: Text(
            buttonText,
            style: AppTextStyles.body(context),
          ),
        ),
      ],
    );
  }
}