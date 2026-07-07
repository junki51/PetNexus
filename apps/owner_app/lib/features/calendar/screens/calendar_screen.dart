import 'package:flutter/material.dart';

import '../../../core/constants/app_colors.dart';
import '../../../core/constants/app_text_styles.dart';
import '../../../layout/responsive_layout.dart';
import '../../home/widgets/activity_list_tile.dart';
import '../../pet/models/pet_model.dart';
import '../controllers/calendar_controller.dart';
import 'package:provider/provider.dart';
import '../../pet/controllers/pet_controller.dart';

class CalendarScreen extends StatefulWidget {
  const CalendarScreen({super.key});

  @override
  State<CalendarScreen> createState() => _CalendarScreenState();
}

class _CalendarScreenState extends State<CalendarScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<CalendarController>().init();
      context.read<PetController>().fetchMyPets();
    });
  }

  @override
  Widget build(BuildContext context) {
    final ctrl = context.watch<CalendarController>();
    final petCtrl = context.watch<PetController>();
    final pets = petCtrl.myPets;
    final selectedPet = pets.isNotEmpty ? pets[ctrl.selectedPetIndex] : null;

    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: Column(
          children: [
            // Header with pet selector
            _buildHeader(context, ctrl, pets, selectedPet),
            SizedBox(height: context.nh(8)),
            // Month calendar
            _buildCalendar(context, ctrl),
            Divider(color: AppColors.border, height: context.nh(16),
                indent: context.nw(16), endIndent: context.nw(16)),
            // Activities section header
            Padding(
              padding: EdgeInsets.symmetric(horizontal: context.nw(16)),
              child: Row(
                children: [
                  Text('กิจกรรมในเดือนนี้',
                      style: AppTextStyles.body(context).copyWith(
                          fontWeight: FontWeight.bold)),
                ],
              ),
            ),
            SizedBox(height: context.nh(8)),
            // Activity list
            Expanded(
              child: ctrl.activities.isEmpty
                  ? Center(
                      child: Text('ไม่มีกิจกรรม',
                          style: AppTextStyles.caption(context)),
                    )
                  : ListView.builder(
                      padding: EdgeInsets.symmetric(
                          horizontal: context.nw(16), vertical: 0),
                      itemCount: ctrl.activities.length,
                      itemBuilder: (ctx, i) =>
                          ActivityListTile(activity: ctrl.activities[i]),
                    ),
            ),
            // Add button
            Padding(
              padding: EdgeInsets.symmetric(
                  horizontal: context.nw(16), vertical: context.nh(12)),
              child: SizedBox(
                width: double.infinity,
                child: ElevatedButton.icon(
                  onPressed: () => ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('เพิ่มกำหนดการ — Coming soon!')),
                  ),
                  icon: const Icon(Icons.add_rounded),
                  label: const Text('เพิ่มกำหนดการ'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.primary,
                    foregroundColor: Colors.white,
                    padding: EdgeInsets.symmetric(vertical: context.nh(14)),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(context.radius(14)),
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context, CalendarController ctrl,
      List<PetModel> pets, PetModel? selectedPet) {
    return Padding(
      padding: EdgeInsets.symmetric(
          horizontal: context.nw(20), vertical: context.nh(16)),
      child: Row(
        children: [
          // Pet dropdown
          Expanded(
            child: pets.isEmpty
                ? Text('ไม่มีสัตว์เลี้ยง',
                    style: AppTextStyles.body(context)
                        .copyWith(color: AppColors.textSecondary))
                : DropdownButton<int>(
                    value: ctrl.selectedPetIndex,
                    underline: const SizedBox.shrink(),
                    items: pets.asMap().entries.map((e) {
                      return DropdownMenuItem(
                        value: e.key,
                        child: Row(
                          children: [
                            Text(e.value.species == 'dog' ? '🐶' : '🐱',
                                style: TextStyle(fontSize: context.nf(18))),
                            SizedBox(width: context.nw(8)),
                            Text(e.value.name,
                                style: AppTextStyles.body(context).copyWith(
                                    fontWeight: FontWeight.w600)),
                            Icon(Icons.keyboard_arrow_down_rounded,
                                color: AppColors.textSecondary,
                                size: context.icon(18)),
                          ],
                        ),
                      );
                    }).toList(),
                    onChanged: (i) {
                      if (i != null) ctrl.setSelectedPetIndex(i);
                    },
                  ),
          ),
          _NotificationBell(),
        ],
      ),
    );
  }

  Widget _buildCalendar(BuildContext context, CalendarController ctrl) {
    final focused = ctrl.focusedMonth;
    const weekdays = ['อา', 'จ', 'อ', 'พ', 'พฤ', 'ศ', 'ส'];

    // Days in this month
    final firstDay = DateTime(focused.year, focused.month, 1);
    final daysInMonth = DateTime(focused.year, focused.month + 1, 0).day;
    final startWeekday = firstDay.weekday % 7; // 0=Sun

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: context.nw(16)),
      child: Column(
        children: [
          // Month/Year navigation
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              IconButton(
                onPressed: ctrl.prevMonth,
                icon: Icon(Icons.chevron_left_rounded,
                    color: AppColors.textPrimary, size: context.icon(22)),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
              SizedBox(width: context.nw(8)),
              Text(_monthThai(focused.month),
                  style: AppTextStyles.body(context)
                      .copyWith(fontWeight: FontWeight.bold)),
              SizedBox(width: context.nw(4)),
              IconButton(
                onPressed: ctrl.nextMonth,
                icon: Icon(Icons.chevron_right_rounded,
                    color: AppColors.textPrimary, size: context.icon(22)),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
              SizedBox(width: context.nw(16)),
              IconButton(
                onPressed: ctrl.prevYear,
                icon: Icon(Icons.chevron_left_rounded,
                    color: AppColors.textPrimary, size: context.icon(22)),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
              SizedBox(width: context.nw(4)),
              Text('${focused.year + 543}',
                  style: AppTextStyles.body(context)
                      .copyWith(fontWeight: FontWeight.bold)),
              SizedBox(width: context.nw(4)),
              IconButton(
                onPressed: ctrl.nextYear,
                icon: Icon(Icons.chevron_right_rounded,
                    color: AppColors.textPrimary, size: context.icon(22)),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
            ],
          ),
          // Weekday headers
          Row(
            children: weekdays.map((d) {
              return Expanded(
                child: Center(
                  child: Text(d,
                      style: TextStyle(
                          fontSize: context.nf(12),
                          color: AppColors.textSecondary,
                          fontWeight: FontWeight.w600)),
                ),
              );
            }).toList(),
          ),
          SizedBox(height: context.nh(4)),
          // Day grid
          GridView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 7,
              childAspectRatio: 1,
            ),
            itemCount: startWeekday + daysInMonth,
            itemBuilder: (ctx, index) {
              if (index < startWeekday) return const SizedBox.shrink();
              final day = index - startWeekday + 1;
              final date = DateTime(focused.year, focused.month, day);
              final isToday = _isSameDay(date, DateTime.now());
              final isSelected = _isSameDay(date, ctrl.selectedDate);
              final hasActivity = ctrl.hasActivityOnDate(date);

              return GestureDetector(
                onTap: () => ctrl.setSelectedDate(date),
                child: Container(
                  margin: EdgeInsets.all(context.nw(2)),
                  decoration: BoxDecoration(
                    color: isSelected
                        ? AppColors.primary
                        : isToday
                            ? AppColors.primaryLight
                            : null,
                    shape: BoxShape.circle,
                  ),
                  child: Stack(
                    alignment: Alignment.center,
                    children: [
                      Text(
                        '$day',
                        style: TextStyle(
                          fontSize: context.nf(13),
                          fontWeight: isToday || isSelected
                              ? FontWeight.bold
                              : FontWeight.normal,
                          color: isSelected ? Colors.white : AppColors.textPrimary,
                        ),
                      ),
                      if (hasActivity)
                        Positioned(
                          bottom: context.nh(4),
                          child: Container(
                            width: context.nw(4),
                            height: context.nw(4),
                            decoration: BoxDecoration(
                              color: isSelected
                                  ? Colors.white
                                  : AppColors.primary,
                              shape: BoxShape.circle,
                            ),
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
    );
  }

  bool _isSameDay(DateTime a, DateTime b) =>
      a.year == b.year && a.month == b.month && a.day == b.day;

  String _monthThai(int m) {
    const months = [
      'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน',
      'พฤษภาคม', 'มิถุนายน', 'กรกฎาคม', 'สิงหาคม',
      'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
    ];
    return months[m - 1];
  }
}

class _NotificationBell extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Icon(Icons.notifications_outlined,
            color: AppColors.textPrimary, size: context.icon(26)),
        Positioned(
          right: 0,
          top: 0,
          child: Container(
            width: context.nw(8),
            height: context.nw(8),
            decoration: const BoxDecoration(
                color: Colors.red, shape: BoxShape.circle),
          ),
        ),
      ],
    );
  }
}
