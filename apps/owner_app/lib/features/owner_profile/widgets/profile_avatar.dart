import 'dart:io';
import 'package:flutter/material.dart';
import '../../../shared/widgets/app_avatar.dart';

class ProfileAvatar extends StatelessWidget {
  final File? imageFile;
  final VoidCallback onTap;

  const ProfileAvatar({
    super.key,
    this.imageFile,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return AppAvatar(
      imageFile: imageFile,
      onTap: onTap,
      showUploadButton: false, // Hide the upload text button as in mockup design
      editIcon: Icons.file_upload_outlined, // Use upload icon from design instead of camera
    );
  }
}
