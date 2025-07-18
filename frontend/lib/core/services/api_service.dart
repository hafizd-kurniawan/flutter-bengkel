import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/api_models.dart';

class ApiService {
  final Dio _dio;
  final String baseUrl;

  ApiService({
    required this.baseUrl,
    Dio? dio,
  }) : _dio = dio ?? Dio() {
    _setupDio();
  }

  void _setupDio() {
    _dio.options.baseUrl = baseUrl;
    _dio.options.connectTimeout = const Duration(seconds: 30);
    _dio.options.receiveTimeout = const Duration(seconds: 30);
    _dio.options.headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    // Add request/response interceptors for logging
    _dio.interceptors.add(
      LogInterceptor(
        requestBody: true,
        responseBody: true,
        logPrint: (log) => print('API: $log'),
      ),
    );
  }

  // Health check
  Future<Map<String, dynamic>> healthCheck() async {
    try {
      final response = await _dio.get('/health');
      return response.data as Map<String, dynamic>;
    } catch (e) {
      throw _handleError(e);
    }
  }

  // Outlet operations with soft delete support
  Future<List<OutletModel>> getOutlets({bool includeDeleted = false}) async {
    try {
      final response = await _dio.get(
        '/api/v1/outlets',
        queryParameters: includeDeleted ? {'include_deleted': 'true'} : null,
      );

      final apiResponse = ApiResponse<List<dynamic>>.fromJson(
        response.data,
        (data) => data as List<dynamic>,
      );

      if (!apiResponse.success) {
        throw ApiException(
          message: apiResponse.message,
          error: apiResponse.error,
        );
      }

      return (apiResponse.data ?? [])
          .map((json) => OutletModel.fromJson(json as Map<String, dynamic>))
          .toList();
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<OutletModel> getOutletById(
    String id, {
    bool includeDeleted = false,
  }) async {
    try {
      final response = await _dio.get(
        '/api/v1/outlets/$id',
        queryParameters: includeDeleted ? {'include_deleted': 'true'} : null,
      );

      final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
        response.data,
        (data) => data as Map<String, dynamic>,
      );

      if (!apiResponse.success) {
        throw ApiException(
          message: apiResponse.message,
          error: apiResponse.error,
        );
      }

      return OutletModel.fromJson(apiResponse.data!);
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<OutletModel> createOutlet(CreateOutletRequest request) async {
    try {
      final response = await _dio.post(
        '/api/v1/outlets',
        data: request.toJson(),
      );

      final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
        response.data,
        (data) => data as Map<String, dynamic>,
      );

      if (!apiResponse.success) {
        throw ApiException(
          message: apiResponse.message,
          error: apiResponse.error,
        );
      }

      return OutletModel.fromJson(apiResponse.data!);
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> softDeleteOutlet(String id) async {
    try {
      final response = await _dio.delete('/api/v1/outlets/$id');

      final apiResponse = ApiResponse<void>.fromJson(response.data, null);

      if (!apiResponse.success) {
        throw ApiException(
          message: apiResponse.message,
          error: apiResponse.error,
        );
      }
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> restoreOutlet(String id) async {
    try {
      final response = await _dio.post('/api/v1/outlets/$id/restore');

      final apiResponse = ApiResponse<void>.fromJson(response.data, null);

      if (!apiResponse.success) {
        throw ApiException(
          message: apiResponse.message,
          error: apiResponse.error,
        );
      }
    } catch (e) {
      throw _handleError(e);
    }
  }

  Exception _handleError(dynamic error) {
    if (error is DioException) {
      switch (error.type) {
        case DioExceptionType.connectionTimeout:
        case DioExceptionType.sendTimeout:
        case DioExceptionType.receiveTimeout:
          return ApiException(
            message: 'Connection timeout. Please try again.',
            statusCode: error.response?.statusCode,
          );
        case DioExceptionType.badResponse:
          final statusCode = error.response?.statusCode;
          final data = error.response?.data;
          
          if (data is Map<String, dynamic>) {
            return ApiException(
              message: data['message'] as String? ?? 'Server error occurred',
              error: data['error'] as String?,
              statusCode: statusCode,
            );
          }
          
          return ApiException(
            message: 'Server error occurred',
            statusCode: statusCode,
          );
        case DioExceptionType.cancel:
          return ApiException(message: 'Request was cancelled');
        case DioExceptionType.unknown:
        default:
          return ApiException(
            message: 'Network error. Please check your connection.',
          );
      }
    }
    
    if (error is ApiException) {
      return error;
    }
    
    return ApiException(
      message: error.toString(),
    );
  }
}

class ApiException implements Exception {
  final String message;
  final String? error;
  final int? statusCode;

  const ApiException({
    required this.message,
    this.error,
    this.statusCode,
  });

  @override
  String toString() {
    final buffer = StringBuffer('ApiException: $message');
    if (statusCode != null) {
      buffer.write(' (Status: $statusCode)');
    }
    if (error != null) {
      buffer.write(' - $error');
    }
    return buffer.toString();
  }
}

// Provider for API service
final apiServiceProvider = Provider<ApiService>((ref) {
  return ApiService(
    baseUrl: 'http://localhost:8080', // Development URL
  );
});

// Provider for outlets data with soft delete support
final outletsProvider = FutureProvider.family<List<OutletModel>, bool>(
  (ref, includeDeleted) async {
    final apiService = ref.watch(apiServiceProvider);
    return apiService.getOutlets(includeDeleted: includeDeleted);
  },
);

final outletByIdProvider = FutureProvider.family<OutletModel, String>(
  (ref, id) async {
    final apiService = ref.watch(apiServiceProvider);
    return apiService.getOutletById(id);
  },
);

// State provider to track whether to include deleted outlets
final includeDeletedProvider = StateProvider<bool>((ref) => false);