import 'package:flutter/material.dart';

class ActivityItem {
  final IconData icon;
  final String title;
  final String subtitle;
  final String time;
  final Color color;

  const ActivityItem({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.time,
    required this.color,
  });
}

class HomeController extends ChangeNotifier {
  DateTime _selectedDate = DateTime.now();
  DateTime get selectedDate => _selectedDate;

  int _selectedPetIndex = 0;
  int get selectedPetIndex => _selectedPetIndex;

  List<ActivityItem> _todayActivities = [];
  List<ActivityItem> get todayActivities => _todayActivities;

  void init() {
    _loadMockActivities();
  }

  void setSelectedDate(DateTime date) {
    _selectedDate = date;
    _loadMockActivities();
    notifyListeners();
  }

  void setSelectedPetIndex(int index) {
    _selectedPetIndex = index;
    notifyListeners();
  }

  void _loadMockActivities() {
    final today = DateTime.now();
    final isToday = _selectedDate.year == today.year &&
        _selectedDate.month == today.month &&
        _selectedDate.day == today.day;

    if (isToday) {
      _todayActivities = [
        ActivityItem(
          icon: Icons.medical_services_rounded,
          title: 'ไปหาหมอ',
          subtitle: 'คลินิกรักสัตว์ สุขุมวิท',
          time: '10:00',
          color: const Color(0xFF66BB6A),
        ),
        ActivityItem(
          icon: Icons.vaccines_rounded,
          title: 'รับยาถ่ายพยาธิ',
          subtitle: 'นัดครั้งที่ 2',
          time: '14:00',
          color: const Color(0xFFFFA726),
        ),
      ];
    } else {
      _todayActivities = [];
    }
    notifyListeners();
  }
}
