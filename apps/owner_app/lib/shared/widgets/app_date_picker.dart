import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import 'app_text_field.dart';

class AppDatePicker extends StatelessWidget {
  final TextEditingController controller;

  final String label;

  const AppDatePicker({
    super.key,
    required this.controller,
    required this.label,
  });

  Future<void> _pickDate(
    BuildContext context,
  ) async {
    final date = await showDatePicker(
      context: context,
      firstDate: DateTime(1900),
      lastDate: DateTime.now(),
      initialDate: DateTime.now(),
    );

    if (date == null) return;

    controller.text = DateFormat(
      "dd/MM/yyyy",
    ).format(date);
  }

  @override
  Widget build(BuildContext context) {
    return AppTextField(
      controller: controller,

      label: label,

      hintText: "Select date",

      readOnly: true,

      prefixIcon: Icons.calendar_month,

      onTap: () => _pickDate(context),
    );
  }
}