import 'package:flutter/material.dart';
import '../../../core/constants/app_colors.dart';
import '../../../layout/responsive_layout.dart';
import '../../pet/models/pet_model.dart';

class PetAvatarWidget extends StatefulWidget {
  final PetModel? pet;

  const PetAvatarWidget({super.key, this.pet});

  @override
  State<PetAvatarWidget> createState() => _PetAvatarWidgetState();
}

class _PetAvatarWidgetState extends State<PetAvatarWidget>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _bounceAnim;
  double _shakeValue = 0;
  bool _isShaking = false;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1200),
    )..repeat(reverse: true);

    _bounceAnim = Tween<double>(begin: 0, end: -8).animate(
      CurvedAnimation(parent: _controller, curve: Curves.easeInOut),
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _onTap() async {
    if (_isShaking) return;
    setState(() => _isShaking = true);
    // Play shake
    final shakeCtrl = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 500),
    );
    final shakeAnim = TweenSequence([
      TweenSequenceItem(tween: Tween<double>(begin: 0, end: 12), weight: 1),
      TweenSequenceItem(tween: Tween<double>(begin: 12, end: -12), weight: 2),
      TweenSequenceItem(tween: Tween<double>(begin: -12, end: 8), weight: 2),
      TweenSequenceItem(tween: Tween<double>(begin: 8, end: -8), weight: 2),
      TweenSequenceItem(tween: Tween<double>(begin: -8, end: 0), weight: 1),
    ]).animate(shakeCtrl);

    shakeCtrl.addListener(() {
      setState(() {
        _shakeValue = shakeAnim.value;
      });
    });
    await shakeCtrl.forward();
    shakeCtrl.dispose();
    setState(() {
      _shakeValue = 0;
      _isShaking = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final pet = widget.pet;
    final avatarSize = context.nw(160);

    return GestureDetector(
      onTap: _onTap,
      child: AnimatedBuilder(
        animation: _controller,
        builder: (_, child) {
          return Transform.translate(
            offset: Offset(_shakeValue, _bounceAnim.value),
            child: child,
          );
        },
        child: Container(
          width: avatarSize,
          height: avatarSize,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: Colors.white,
            boxShadow: [
              BoxShadow(
                color: AppColors.primary.withValues(alpha: 0.2),
                blurRadius: 20,
                spreadRadius: 4,
              ),
            ],
          ),
          child: ClipOval(
            child: _buildAvatarContent(context, pet, avatarSize),
          ),
        ),
      ),
    );
  }

  Widget _buildAvatarContent(
      BuildContext context, PetModel? pet, double size) {
    // Use real avatar if available
    if (pet?.avatarUrl != null && pet!.avatarUrl!.isNotEmpty) {
      return Image.network(
        pet.avatarUrl!,
        fit: BoxFit.cover,
        errorBuilder: (ctx, err, stack) =>
            _buildPlaceholder(context, pet.species, size),
      );
    }
    // Fallback placeholder based on species
    return _buildPlaceholder(context, pet?.species ?? 'dog', size);
  }

  Widget _buildPlaceholder(BuildContext context, String species, double size) {
    final isDog = species == 'dog';
    return Container(
      color: AppColors.primaryLight,
      child: Center(
        child: Text(
          isDog ? '🐶' : '🐱',
          style: TextStyle(fontSize: size * 0.5),
        ),
      ),
    );
  }
}
