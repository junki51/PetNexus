import 'package:flutter/material.dart';

import '../../core/constants/app_text_styles.dart';
import '../../layout/responsive_layout.dart';

class AppDropdown<T> extends StatelessWidget {
  final String label;

  final T? value;

  final List<T> items;

  final ValueChanged<T?> onChanged;

  const AppDropdown({
    super.key,
    required this.label,
    required this.items,
    required this.onChanged,
    this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: AppTextStyles.body(context),
        ),

        SizedBox(
          height: context.nh(8),
        ),

        DropdownButtonFormField<T>(
          initialValue: value,

          decoration: InputDecoration(
            filled: true,
            fillColor: Colors.white,

            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(context.radius(18)),
              borderSide: const BorderSide(color: Colors.black, width: 1.5), // เปลี่ยนเป็นสีดำ
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(context.radius(18)),
              borderSide: const BorderSide(color: Colors.black, width: 1.5), // เปลี่ยนเป็นสีดำ
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(context.radius(18)),
              borderSide: const BorderSide(color: Colors.black, width: 2.5), // เปลี่ยนเป็นสีดำหนาขึ้น
            ),
          ),

          items: items.map((item) {
            return DropdownMenuItem<T>(
              value: item,
              child: Text(
                item.toString(),
              ),
            );
          }).toList(),

          onChanged: onChanged,
        ),
      ],
    );
  }
}