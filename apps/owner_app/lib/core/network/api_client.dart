import 'package:dio/dio.dart';
import '../storage/token_storage.dart';
import 'api_config.dart';

class ApiClient {
  ApiClient._();

  static final ApiClient instance = ApiClient._();

  late final Dio dio = Dio(
    BaseOptions(
      baseUrl: ApiConfig.baseUrl,
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 15),
      headers: {
        "Content-Type": "application/json",
      },
    ),
  )..interceptors.add(
      InterceptorsWrapper(
        onRequest: (
          options,
          handler,
        ) async {

          final token =
              await TokenStorage.instance.getToken();

          if (token != null) {
            options.headers["Authorization"] =
                "Bearer $token";
          }

          handler.next(options);
        },
      ),
    );
}